package main

import (
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
