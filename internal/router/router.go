package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"yourz-itinerary/infrastructure/database"
	permissioncache "yourz-itinerary/internal/cache/permission"
	appConfigHandler "yourz-itinerary/internal/handlers/http/appconfig"
	auditHandler "yourz-itinerary/internal/handlers/http/audit"
	itineraryDayHandler "yourz-itinerary/internal/handlers/http/itineraryday"
	itineraryItemHandler "yourz-itinerary/internal/handlers/http/itineraryitem"
	locationHandler "yourz-itinerary/internal/handlers/http/location"
	menuHandler "yourz-itinerary/internal/handlers/http/menu"
	permissionHandler "yourz-itinerary/internal/handlers/http/permission"
	roleHandler "yourz-itinerary/internal/handlers/http/role"
	sessionHandler "yourz-itinerary/internal/handlers/http/session"
	tripHandler "yourz-itinerary/internal/handlers/http/trip"
	tripMemberHandler "yourz-itinerary/internal/handlers/http/tripmember"
	userHandler "yourz-itinerary/internal/handlers/http/user"
	interfacesession "yourz-itinerary/internal/interfaces/session"
	appConfigRepo "yourz-itinerary/internal/repositories/appconfig"
	auditRepo "yourz-itinerary/internal/repositories/audit"
	authRepo "yourz-itinerary/internal/repositories/auth"
	itineraryDayRepo "yourz-itinerary/internal/repositories/itineraryday"
	itineraryItemRepo "yourz-itinerary/internal/repositories/itineraryitem"
	locationRepo "yourz-itinerary/internal/repositories/location"
	menuRepo "yourz-itinerary/internal/repositories/menu"
	otpRepo "yourz-itinerary/internal/repositories/otp"
	permissionRepo "yourz-itinerary/internal/repositories/permission"
	resetRepo "yourz-itinerary/internal/repositories/reset"
	roleRepo "yourz-itinerary/internal/repositories/role"
	sessionRepo "yourz-itinerary/internal/repositories/session"
	tripRepo "yourz-itinerary/internal/repositories/trip"
	tripMemberRepo "yourz-itinerary/internal/repositories/tripmember"
	userRepo "yourz-itinerary/internal/repositories/user"
	appConfigSvc "yourz-itinerary/internal/services/appconfig"
	auditSvc "yourz-itinerary/internal/services/audit"
	itineraryDaySvc "yourz-itinerary/internal/services/itineraryday"
	itineraryItemSvc "yourz-itinerary/internal/services/itineraryitem"
	locationSvc "yourz-itinerary/internal/services/location"
	menuSvc "yourz-itinerary/internal/services/menu"
	otpSvc "yourz-itinerary/internal/services/otp"
	permissionSvc "yourz-itinerary/internal/services/permission"
	resetSvc "yourz-itinerary/internal/services/reset"
	roleSvc "yourz-itinerary/internal/services/role"
	sessionSvc "yourz-itinerary/internal/services/session"
	tripSvc "yourz-itinerary/internal/services/trip"
	tripMemberSvc "yourz-itinerary/internal/services/tripmember"
	userSvc "yourz-itinerary/internal/services/user"
	"yourz-itinerary/middlewares"
	"yourz-itinerary/pkg/config"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/mailer"
	"yourz-itinerary/pkg/security"
	"yourz-itinerary/utils"
)

type Routes struct {
	App *gin.Engine
	DB  *gorm.DB
}

func NewRoutes() *Routes {
	app := gin.Default()

	app.Use(middlewares.CORS())
	app.Use(gin.CustomRecovery(middlewares.ErrorHandler))
	app.Use(middlewares.SetContextId())
	app.Use(middlewares.RequestLogger())

	app.GET("/healthcheck", func(ctx *gin.Context) {
		logger.WriteLogWithContext(ctx, logger.LogLevelDebug, "ClientIP: "+ctx.ClientIP())
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK!!",
		})
	})

	return &Routes{
		App: app,
	}
}

func (r *Routes) UserRoutes() {
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	repo := userRepo.NewUserRepo(r.DB)
	rRepo := roleRepo.NewRoleRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	redisClient := database.GetRedisClient()
	permissionInvalidator := permissioncache.NewInvalidator(redisClient)
	uc := userSvc.NewUserService(repo, blacklistRepo, rRepo, pRepo, permissionInvalidator)
	var userSessionSvc interfacesession.ServiceSessionInterface
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	repoAppConfig := appConfigRepo.NewAppConfigRepo(r.DB)
	svcAppConfig := appConfigSvc.NewAppConfigService(repoAppConfig)

	// Setup login limiter if Redis is available
	var loginLimiter security.LoginLimiter
	var registerOTPService = otpSvc.NewOTPService(nil, nil, config.LoadOTPConfig())
	var passwordResetService = resetSvc.NewPasswordResetService(nil, nil, config.LoadPasswordResetConfig())
	if redisClient != nil {
		loginLimiter = security.NewRedisLoginLimiter(
			redisClient,
			utils.GetEnv("LOGIN_ATTEMPT_LIMIT", 5),
			time.Duration(utils.GetEnv("LOGIN_ATTEMPT_WINDOW_SECONDS", 60))*time.Second,
			time.Duration(utils.GetEnv("LOGIN_BLOCK_DURATION_SECONDS", 300))*time.Second,
		)

		sRepo := sessionRepo.NewSessionRepository(redisClient)
		userSessionSvc = sessionSvc.NewSessionService(sRepo)

		sender, err := mailer.NewBrevoSenderFromEnv()
		if err != nil {
			logger.WriteLog(logger.LogLevelWarn, "Email sender not configured: ", err)
		} else {
			registerOTPService = otpSvc.NewOTPService(otpRepo.NewOTPRepository(redisClient), sender, config.LoadOTPConfig())
			passwordResetService = resetSvc.NewPasswordResetService(resetRepo.NewPasswordResetRepository(redisClient), sender, config.LoadPasswordResetConfig())
		}
	}

	h := userHandler.NewUserHandler(uc, blacklistRepo, userSessionSvc, loginLimiter, svcAudit, svcAppConfig, registerOTPService, passwordResetService)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Setup register rate limiter
	registerLimit := utils.GetEnv("REGISTER_RATE_LIMIT", 5)
	registerWindowSeconds := utils.GetEnv("REGISTER_RATE_WINDOW_SECONDS", 60)
	if registerWindowSeconds <= 0 {
		registerWindowSeconds = 60
	}
	registerLimiter := middlewares.IPRateLimitMiddleware(
		redisClient,
		"user_register",
		registerLimit,
		time.Duration(registerWindowSeconds)*time.Second,
	)

	user := r.App.Group("/api/user")
	{
		user.GET("/register/status", h.GetRegisterStatus)
		user.POST("/register", registerLimiter, h.Register)
		user.POST("/register/otp/send", registerLimiter, h.SendRegisterOTP)
		user.POST("/login", h.Login)
		user.POST("/google/login", h.GoogleLogin)
		user.POST("/refresh-token", h.RefreshToken)
		user.POST("/forgot-password", h.ForgotPassword)
		user.POST("/reset-password", h.ResetPassword)

		userPriv := user.Group("").Use(mdw.AuthMiddleware())
		{
			userPriv.POST("/logout", h.Logout)
			userPriv.GET("", h.GetUserByAuth)
			userPriv.GET("/:id", mdw.PermissionMiddleware("users", "view"), h.GetUserById)
			userPriv.POST("/:id/impersonate", mdw.PermissionMiddleware("users", "impersonate"), h.ImpersonateUser)
			userPriv.POST("/stop-impersonation", h.StopImpersonation)
			userPriv.PUT("", h.Update)
			userPriv.PUT("/:id", mdw.PermissionMiddleware("users", "update"), h.UpdateUserById)
			userPriv.PUT("/change/password", h.ChangePassword)
			userPriv.DELETE("", h.Delete)
			userPriv.DELETE("/:id", mdw.PermissionMiddleware("users", "delete"), h.DeleteUserById)

			// Admin create user endpoint (with role selection)
			userPriv.POST("", mdw.PermissionMiddleware("users", "create"), h.AdminCreateUser)
		}
	}

	r.App.GET("/api/users", mdw.AuthMiddleware(), mdw.PermissionMiddleware("users", "list"), h.GetAllUsers)
}

func (r *Routes) RoleRoutes() {
	repoRole := roleRepo.NewRoleRepo(r.DB)
	repoPermission := permissionRepo.NewPermissionRepo(r.DB)
	repoMenu := menuRepo.NewMenuRepo(r.DB)
	permissionInvalidator := permissioncache.NewInvalidator(database.GetRedisClient())
	svc := roleSvc.NewRoleService(repoRole, repoPermission, repoMenu, permissionInvalidator)
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	h := roleHandler.NewRoleHandler(svc, svcAudit)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, repoPermission)

	// List endpoints
	r.App.GET("/api/roles", mdw.AuthMiddleware(), mdw.PermissionMiddleware("roles", "list"), h.GetAll)

	// CRUD endpoints
	role := r.App.Group("/api/role").Use(mdw.AuthMiddleware())
	{
		role.POST("", mdw.PermissionMiddleware("roles", "create"), h.Create)
		role.GET("/:id", mdw.PermissionMiddleware("roles", "view"), h.GetByID)
		role.PUT("/:id", mdw.PermissionMiddleware("roles", "update"), h.Update)
		role.DELETE("/:id", mdw.PermissionMiddleware("roles", "delete"), h.Delete)

		// Permission and menu assignment
		role.POST("/:id/permissions", mdw.PermissionMiddleware("roles", "assign_permissions"), h.AssignPermissions)
	}
}

func (r *Routes) PermissionRoutes() {
	repo := permissionRepo.NewPermissionRepo(r.DB)
	permissionInvalidator := permissioncache.NewInvalidator(database.GetRedisClient())
	svc := permissionSvc.NewPermissionService(repo, permissionInvalidator)
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	h := permissionHandler.NewPermissionHandler(svc, svcAudit)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, repo)

	// List endpoints
	r.App.GET("/api/permissions", mdw.AuthMiddleware(), mdw.PermissionMiddleware("permissions", "list"), h.GetAll)

	// Get current user's permissions
	r.App.GET("/api/permissions/me", mdw.AuthMiddleware(), h.GetUserPermissions)

	// CRUD endpoints
	permission := r.App.Group("/api/permission").Use(mdw.AuthMiddleware())
	{
		permission.POST("", mdw.PermissionMiddleware("permissions", "create"), h.Create)
		permission.GET("/:id", mdw.PermissionMiddleware("permissions", "view"), h.GetByID)
		permission.PUT("/:id", mdw.PermissionMiddleware("permissions", "update"), h.Update)
		permission.DELETE("/:id", mdw.PermissionMiddleware("permissions", "delete"), h.Delete)
	}
}

func (r *Routes) MenuRoutes() {
	repo := menuRepo.NewMenuRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	svc := menuSvc.NewMenuService(repo, pRepo)
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	h := menuHandler.NewMenuHandler(svc, svcAudit)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Public endpoints for authenticated users
	r.App.GET("/api/menus/active", mdw.AuthMiddleware(), h.GetActiveMenus)
	r.App.GET("/api/menus/me", mdw.AuthMiddleware(), h.GetUserMenus)

	// List endpoints
	r.App.GET("/api/menus", mdw.AuthMiddleware(), mdw.PermissionMiddleware("menus", "list"), h.GetAll)

	// CRUD endpoints
	menu := r.App.Group("/api/menu").Use(mdw.AuthMiddleware())
	{
		menu.GET("/:id", mdw.PermissionMiddleware("menus", "view"), h.GetByID)
		menu.PUT("/:id", mdw.PermissionMiddleware("menus", "update"), h.Update)
	}
}

func (r *Routes) AppConfigRoutes() {
	repo := appConfigRepo.NewAppConfigRepo(r.DB)
	svc := appConfigSvc.NewAppConfigService(repo)
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	h := appConfigHandler.NewAppConfigHandler(svc, svcAudit)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	r.App.GET("/api/configs", mdw.AuthMiddleware(), mdw.PermissionMiddleware("configs", "list"), h.GetAll)

	appConfig := r.App.Group("/api/config").Use(mdw.AuthMiddleware())
	{
		appConfig.GET("/:id", mdw.PermissionMiddleware("configs", "view"), h.GetByID)
		appConfig.PUT("/:id", mdw.PermissionMiddleware("configs", "update"), h.Update)
	}
}

func (r *Routes) AuditRoutes() {
	repo := auditRepo.NewAuditRepo(r.DB)
	svc := auditSvc.NewAuditService(repo)
	h := auditHandler.NewAuditHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	r.App.GET("/api/audits", mdw.AuthMiddleware(), mdw.PermissionMiddleware("audits", "list"), h.GetAll)

	audit := r.App.Group("/api/audit").Use(mdw.AuthMiddleware())
	{
		audit.GET("/:id", mdw.PermissionMiddleware("audits", "view"), h.GetByID)
	}
}

func (r *Routes) SessionRoutes() {
	redisClient := database.GetRedisClient()
	if redisClient == nil {
		logger.WriteLog(logger.LogLevelDebug, "Redis not available, session routes will not be registered")
		return
	}

	repo := sessionRepo.NewSessionRepository(redisClient)
	svc := sessionSvc.NewSessionService(repo)
	repoAudit := auditRepo.NewAuditRepo(r.DB)
	svcAudit := auditSvc.NewAuditService(repoAudit)
	h := sessionHandler.NewSessionHandler(svc, svcAudit)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	// Session management endpoints (authenticated users only)
	sessionGroup := r.App.Group("/api/user").Use(mdw.AuthMiddleware())
	{
		sessionGroup.GET("/sessions", h.GetActiveSessions)
		sessionGroup.DELETE("/session/:session_id", h.RevokeSession)
		sessionGroup.POST("/sessions/revoke-others", h.RevokeAllOtherSessions)
	}

	logger.WriteLog(logger.LogLevelInfo, "Session management routes registered")
}

func (r *Routes) LocationRoutes() {
	repo := locationRepo.NewLocationRepo(r.DB)
	svc := locationSvc.NewLocationService(repo, database.GetRedisClient())
	h := locationHandler.NewLocationHandler(svc)
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	mdw := middlewares.NewMiddleware(blacklistRepo, pRepo)

	location := r.App.Group("/api/location")
	{
		location.GET("/province", h.GetProvince)
		location.GET("/city", h.GetCity)
		location.GET("/district", h.GetDistrict)
		location.GET("/village", h.GetVillage)
	}

	locationPriv := r.App.Group("/api/location").Use(mdw.AuthMiddleware())
	{
		locationPriv.POST("/sync", mdw.PermissionMiddleware("locations", "sync"), h.Sync)
		locationPriv.GET("/sync/:id", mdw.PermissionMiddleware("locations", "sync"), h.GetSyncJob)
	}
}

func (r *Routes) newItineraryAuthMiddleware() *middlewares.Middleware {
	blacklistRepo := authRepo.NewBlacklistRepo(r.DB)
	pRepo := permissionRepo.NewPermissionRepo(r.DB)
	return middlewares.NewMiddleware(blacklistRepo, pRepo)
}

func (r *Routes) TripRoutes() {
	repo := tripRepo.NewTripRepo(r.DB)
	memberRepo := tripMemberRepo.NewTripMemberRepo(r.DB)
	dayRepo := itineraryDayRepo.NewItineraryDayRepo(r.DB)
	itemRepo := itineraryItemRepo.NewItineraryItemRepo(r.DB)
	uRepo := userRepo.NewUserRepo(r.DB)
	svc := tripSvc.NewTripService(repo, memberRepo, dayRepo, itemRepo, uRepo)
	h := tripHandler.NewTripHandler(svc)
	mdw := r.newItineraryAuthMiddleware()

	trips := r.App.Group("/api/trips").Use(mdw.AuthMiddleware())
	{
		trips.POST("", h.CreateTrip)
		trips.GET("", h.ListTrips)
		trips.GET("/:id", h.GetTripDetail)
		trips.PUT("/:id", h.UpdateTrip)
		trips.DELETE("/:id", h.DeleteTrip)
	}
}

func (r *Routes) TripMemberRoutes() {
	tRepo := tripRepo.NewTripRepo(r.DB)
	mRepo := tripMemberRepo.NewTripMemberRepo(r.DB)
	uRepo := userRepo.NewUserRepo(r.DB)
	svc := tripMemberSvc.NewTripMemberService(tRepo, mRepo, uRepo)
	h := tripMemberHandler.NewTripMemberHandler(svc)
	mdw := r.newItineraryAuthMiddleware()

	trips := r.App.Group("/api/trips").Use(mdw.AuthMiddleware())
	{
		trips.POST("/:id/members", h.AddMember)
		trips.PUT("/:id/members/:member_id", h.UpdateMemberRole)
		trips.DELETE("/:id/members/:member_id", h.RemoveMember)
		trips.DELETE("/:id/leave", h.LeaveTrip)
	}
}

func (r *Routes) ItineraryDayRoutes() {
	mRepo := tripMemberRepo.NewTripMemberRepo(r.DB)
	dRepo := itineraryDayRepo.NewItineraryDayRepo(r.DB)
	tRepo := tripRepo.NewTripRepo(r.DB)
	svc := itineraryDaySvc.NewItineraryDayService(mRepo, dRepo, tRepo)
	h := itineraryDayHandler.NewItineraryDayHandler(svc)
	mdw := r.newItineraryAuthMiddleware()

	trips := r.App.Group("/api/trips").Use(mdw.AuthMiddleware())
	{
		trips.POST("/:id/days", h.CreateDay)
	}

	dayRoutes := r.App.Group("/api/itinerary-days").Use(mdw.AuthMiddleware())
	{
		dayRoutes.PUT("/:id", h.UpdateDay)
		dayRoutes.DELETE("/:id", h.DeleteDay)
	}
}

func (r *Routes) ItineraryItemRoutes() {
	mRepo := tripMemberRepo.NewTripMemberRepo(r.DB)
	dRepo := itineraryDayRepo.NewItineraryDayRepo(r.DB)
	iRepo := itineraryItemRepo.NewItineraryItemRepo(r.DB)
	svc := itineraryItemSvc.NewItineraryItemService(mRepo, dRepo, iRepo)
	h := itineraryItemHandler.NewItineraryItemHandler(svc)
	mdw := r.newItineraryAuthMiddleware()

	dayRoutes := r.App.Group("/api/itinerary-days").Use(mdw.AuthMiddleware())
	{
		dayRoutes.POST("/:id/items", h.CreateItem)
		dayRoutes.PUT("/:id/items/reorder", h.ReorderItems)
	}

	itemRoutes := r.App.Group("/api/itinerary-items").Use(mdw.AuthMiddleware())
	{
		itemRoutes.PUT("/:id", h.UpdateItem)
		itemRoutes.DELETE("/:id", h.DeleteItem)
	}
}
