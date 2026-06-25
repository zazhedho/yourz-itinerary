package handlerrole

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	domainaudit "yourz-itinerary/internal/domain/audit"
	"yourz-itinerary/internal/dto"
	interfaceaudit "yourz-itinerary/internal/interfaces/audit"
	interfacerole "yourz-itinerary/internal/interfaces/role"
	"yourz-itinerary/pkg/filter"
	"yourz-itinerary/pkg/logger"
	"yourz-itinerary/pkg/messages"
	"yourz-itinerary/pkg/response"
	"yourz-itinerary/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	Service      interfacerole.ServiceRoleInterface
	AuditService interfaceaudit.ServiceAuditInterface
}

func NewRoleHandler(s interfacerole.ServiceRoleInterface, auditService interfaceaudit.ServiceAuditInterface) *RoleHandler {
	return &RoleHandler{
		Service:      s,
		AuditService: auditService,
	}
}

func (h *RoleHandler) Create(ctx *gin.Context) {
	var req dto.RoleCreate
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][Create]"
	reqCtx := ctx.Request.Context()

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	data, err := h.Service.Create(reqCtx, req)
	if err != nil {
		h.writeAudit(ctx, domainaudit.AuditEvent{
			Action:       domainaudit.ActionCreate,
			Resource:     "role",
			Status:       domainaudit.StatusFailed,
			Message:      "Failed to create role",
			ErrorMessage: err.Error(),
			AfterData:    req,
		})
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.Create; Error: %+v", logPrefix, err))
		statusCode, res := roleMutationErrorResponse(logId, err)
		ctx.JSON(statusCode, res)
		return
	}
	h.writeAudit(ctx, domainaudit.AuditEvent{
		Action:     domainaudit.ActionCreate,
		Resource:   "role",
		ResourceID: data.Id,
		Status:     domainaudit.StatusSuccess,
		Message:    "Created role",
		AfterData:  data,
	})

	res := response.Response(http.StatusCreated, "Role created successfully", logId, data)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusCreated, res)
}

func (h *RoleHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][GetByID]"
	reqCtx := ctx.Request.Context()

	data, err := h.Service.GetByIDWithDetails(reqCtx, id)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.GetByIDWithDetails; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusNotFound, "Role not found", logId, nil)
		res.Error = response.Errors{Code: http.StatusNotFound, Message: "role not found"}
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := response.Response(http.StatusOK, "Get role successfully", logId, data)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *RoleHandler) GetAll(ctx *gin.Context) {
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][GetAll]"
	reqCtx := ctx.Request.Context()

	params, err := filter.GetBaseParams(ctx, "name", "asc", 10)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; GetBaseParams; Error: %+v", logPrefix, err))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = response.Errors{Code: http.StatusBadRequest, Message: "invalid query parameters"}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, total, err := h.Service.GetAll(reqCtx, params)
	if err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.GetAll; Error: %+v", logPrefix, err))
		res := response.InternalServerError(logId)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.PaginationResponse(http.StatusOK, int(total), params.Page, params.Limit, logId, data)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *RoleHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.RoleUpdate
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][Update]"
	reqCtx := ctx.Request.Context()

	if err := ctx.BindJSON(&req); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, logId, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(req), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", logPrefix, utils.JsonEncode(req)))

	before, _ := h.Service.GetByID(reqCtx, id)
	data, err := h.Service.Update(reqCtx, id, req)
	if err != nil {
		h.writeAudit(ctx, domainaudit.AuditEvent{
			Action:       domainaudit.ActionUpdate,
			Resource:     "role",
			ResourceID:   id,
			Status:       domainaudit.StatusFailed,
			Message:      "Failed to update role",
			ErrorMessage: err.Error(),
			BeforeData:   before,
			AfterData:    req,
		})
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.Update; Error: %+v", logPrefix, err))
		statusCode, res := roleMutationErrorResponse(logId, err)
		ctx.JSON(statusCode, res)
		return
	}
	h.writeAudit(ctx, domainaudit.AuditEvent{
		Action:     domainaudit.ActionUpdate,
		Resource:   "role",
		ResourceID: data.Id,
		Status:     domainaudit.StatusSuccess,
		Message:    "Updated role",
		BeforeData: before,
		AfterData:  data,
	})

	res := response.Response(http.StatusOK, "Role updated successfully", logId, data)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Response: %+v;", logPrefix, utils.JsonEncode(data)))
	ctx.JSON(http.StatusOK, res)
}

func (h *RoleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][Delete]"
	reqCtx := ctx.Request.Context()
	before, _ := h.Service.GetByID(reqCtx, id)

	if err := h.Service.Delete(reqCtx, id); err != nil {
		h.writeAudit(ctx, domainaudit.AuditEvent{
			Action:       domainaudit.ActionDelete,
			Resource:     "role",
			ResourceID:   id,
			Status:       domainaudit.StatusFailed,
			Message:      "Failed to delete role",
			ErrorMessage: err.Error(),
			BeforeData:   before,
		})
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; Service.Delete; Error: %+v", logPrefix, err))
		statusCode, res := roleMutationErrorResponse(logId, err)
		ctx.JSON(statusCode, res)
		return
	}
	h.writeAudit(ctx, domainaudit.AuditEvent{
		Action:     domainaudit.ActionDelete,
		Resource:   "role",
		ResourceID: id,
		Status:     domainaudit.StatusSuccess,
		Message:    "Deleted role",
		BeforeData: before,
	})

	res := response.Response(http.StatusOK, "Role deleted successfully", logId, nil)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Response: Role deleted", logPrefix))
	ctx.JSON(http.StatusOK, res)
}

func (h *RoleHandler) AssignPermissions(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.AssignPermissions
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][AssignPermissions]"
	h.assignRoleRelation(ctx, roleRelationAssignment{
		id:             id,
		request:        &req,
		requestValue:   func() interface{} { return req },
		logID:          logId,
		logPrefix:      logPrefix,
		beforeKey:      "permission_ids",
		resource:       "role_permissions",
		failedMessage:  "Failed to assign permissions to role",
		successMessage: "Assigned permissions to role",
		responseText:   "Permissions assigned successfully",
		debugText:      "Permissions assigned",
		serviceName:    "Service.AssignPermissions",
		before:         h.Service.GetRolePermissions,
		assign: func(reqCtx context.Context, id string) error {
			return h.Service.AssignPermissions(reqCtx, id, req)
		},
	})
}

func (h *RoleHandler) AssignMenus(ctx *gin.Context) {
	id := ctx.Param("id")
	var req dto.AssignMenus
	logId := utils.GenerateLogId(ctx)
	logPrefix := "[RoleHandler][AssignMenus]"
	h.assignRoleRelation(ctx, roleRelationAssignment{
		id:             id,
		request:        &req,
		requestValue:   func() interface{} { return req },
		logID:          logId,
		logPrefix:      logPrefix,
		beforeKey:      "menu_ids",
		resource:       "role_menus",
		failedMessage:  "Failed to assign menus to role",
		successMessage: "Assigned menus to role",
		responseText:   "Menus assigned successfully",
		debugText:      "Menus assigned",
		serviceName:    "Service.AssignMenus",
		before:         h.Service.GetRoleMenus,
		assign: func(reqCtx context.Context, id string) error {
			return h.Service.AssignMenus(reqCtx, id, req)
		},
	})
}

type roleRelationAssignment struct {
	id             string
	request        interface{}
	requestValue   func() interface{}
	logID          uuid.UUID
	logPrefix      string
	beforeKey      string
	resource       string
	failedMessage  string
	successMessage string
	responseText   string
	debugText      string
	serviceName    string
	before         func(context.Context, string) ([]string, error)
	assign         func(context.Context, string) error
}

func (h *RoleHandler) assignRoleRelation(ctx *gin.Context, assignment roleRelationAssignment) {
	reqCtx := ctx.Request.Context()

	if err := ctx.BindJSON(assignment.request); err != nil {
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; BindJSON ERROR: %s;", assignment.logPrefix, err.Error()))
		res := response.Response(http.StatusBadRequest, messages.InvalidRequest, assignment.logID, nil)
		res.Error = utils.ValidateError(err, reflect.TypeOf(assignment.requestValue()), "json")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	requestValue := assignment.requestValue()
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; Request: %+v;", assignment.logPrefix, utils.JsonEncode(requestValue)))

	beforeIDs, _ := assignment.before(reqCtx, assignment.id)
	if err := assignment.assign(reqCtx, assignment.id); err != nil {
		h.writeAudit(ctx, domainaudit.AuditEvent{
			Action:       domainaudit.ActionAssign,
			Resource:     assignment.resource,
			ResourceID:   assignment.id,
			Status:       domainaudit.StatusFailed,
			Message:      assignment.failedMessage,
			ErrorMessage: err.Error(),
			BeforeData: map[string]interface{}{
				assignment.beforeKey: beforeIDs,
			},
			AfterData: requestValue,
		})
		logger.WriteLogWithContext(ctx, logger.LogLevelError, fmt.Sprintf("%s; %s; Error: %+v", assignment.logPrefix, assignment.serviceName, err))
		statusCode, res := roleMutationErrorResponse(assignment.logID, err)
		ctx.JSON(statusCode, res)
		return
	}
	h.writeAudit(ctx, domainaudit.AuditEvent{
		Action:     domainaudit.ActionAssign,
		Resource:   assignment.resource,
		ResourceID: assignment.id,
		Status:     domainaudit.StatusSuccess,
		Message:    assignment.successMessage,
		BeforeData: map[string]interface{}{
			assignment.beforeKey: beforeIDs,
		},
		AfterData: requestValue,
	})

	res := response.Response(http.StatusOK, assignment.responseText, assignment.logID, nil)
	logger.WriteLogWithContext(ctx, logger.LogLevelDebug, fmt.Sprintf("%s; %s", assignment.logPrefix, assignment.debugText))
	ctx.JSON(http.StatusOK, res)
}
