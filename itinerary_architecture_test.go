package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestItineraryFeatureInterfacesExist(t *testing.T) {
	for _, dir := range []string{
		"internal/interfaces/trip",
		"internal/interfaces/tripmember",
		"internal/interfaces/itineraryday",
		"internal/interfaces/itineraryitem",
	} {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatalf("expected interface directory %s: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %s to be a directory", dir)
		}
	}
}

func TestItineraryDomainsDoNotImportSiblingFeatureDomains(t *testing.T) {
	domainFiles := []string{
		"internal/domain/trip/trip.go",
		"internal/domain/tripmember/trip_member.go",
		"internal/domain/itineraryday/itinerary_day.go",
		"internal/domain/itineraryitem/itinerary_item.go",
	}

	for _, file := range domainFiles {
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		feature := filepath.Base(filepath.Dir(file))
		for _, imported := range parsed.Imports {
			path := strings.Trim(imported.Path.Value, `"`)
			if strings.HasPrefix(path, "yourz-itinerary/internal/domain/") && !strings.HasSuffix(path, "/"+feature) {
				t.Fatalf("%s imports sibling feature domain %s", file, path)
			}
		}
	}
}

func TestItineraryRepoInterfacesDoNotRedeclareGenericMethods(t *testing.T) {
	repoFiles := []string{
		"internal/interfaces/trip/repo.go",
		"internal/interfaces/tripmember/repo.go",
		"internal/interfaces/itineraryday/repo.go",
		"internal/interfaces/itineraryitem/repo.go",
	}

	genericMethods := map[string]struct{}{
		"Store":      {},
		"GetByID":    {},
		"GetAll":     {},
		"Update":     {},
		"Delete":     {},
		"SoftDelete": {},
	}

	for _, file := range repoFiles {
		parsed, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		ast.Inspect(parsed, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}
			if !embedsGenericRepository(interfaceType) {
				return true
			}

			for _, method := range interfaceType.Methods.List {
				if len(method.Names) == 0 {
					continue
				}
				name := method.Names[0].Name
				if _, exists := genericMethods[name]; exists {
					t.Fatalf("%s redeclares generic repository method %s in %s", typeSpec.Name.Name, name, file)
				}
			}
			return true
		})
	}
}

func embedsGenericRepository(interfaceType *ast.InterfaceType) bool {
	for _, method := range interfaceType.Methods.List {
		if len(method.Names) != 0 {
			continue
		}
		if exprContainsGenericRepository(method.Type) {
			return true
		}
	}
	return false
}

func exprContainsGenericRepository(expr ast.Expr) bool {
	switch value := expr.(type) {
	case *ast.SelectorExpr:
		return value.Sel.Name == "GenericRepository"
	case *ast.IndexExpr:
		return exprContainsGenericRepository(value.X)
	case *ast.IndexListExpr:
		return exprContainsGenericRepository(value.X)
	default:
		return false
	}
}
