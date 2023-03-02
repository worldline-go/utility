package swagger

import (
	"os"
	"testing"

	"github.com/swaggo/swag"
)

func TestSetInfo(t *testing.T) {
	oauthTest, err := os.ReadFile("testdata/oauth.json")
	if err != nil {
		t.Fatal(err)
	}

	swaggerInfo := &swag.Spec{
		Version:          "",
		Host:             "",
		BasePath:         "",
		Schemes:          []string{},
		Title:            "Auth Test API",
		Description:      "This is a sample server for out Auth library.",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(oauthTest),
	}

	// init function registers the swaggerInfo instance with the docs package
	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)

	// SetInfo modifies the swaggerInfo instance
	if err := SetInfo(
		WithVersion("1.0.0"),
		WithHost("localhost:8080"),
		WithBasePath("/api/v1"),
		WithSchemes("http", "https"),
		WithTitle("Auth Test API modified"),
		WithDescription("This is a sample server for out Auth library. modified"),
		WithCustom(map[string]interface{}{
			"authUrl":  "https://example.com/oauth2/authorize",
			"tokenUrl": "https://example.com/oauth2/token",
		}),
	); err != nil {
		t.Fatal(err)
	}

	// check that the swaggerInfo instance was modified
	if swaggerInfo.Version != "1.0.0" {
		t.Errorf("expected Version to be 1.0.0, got %s", swaggerInfo.Version)
	}
	if swaggerInfo.Host != "localhost:8080" {
		t.Errorf("expected Host to be localhost:8080, got %s", swaggerInfo.Host)
	}

	if swaggerInfo.BasePath != "/api/v1" {
		t.Errorf("expected BasePath to be /api/v1, got %s", swaggerInfo.BasePath)
	}

	if len(swaggerInfo.Schemes) != 2 {
		t.Errorf("expected Schemes to be 2, got %d", len(swaggerInfo.Schemes))
	}

	if swaggerInfo.Title != "Auth Test API modified" {
		t.Errorf("expected Title to be Auth Test API modified, got %s", swaggerInfo.Title)
	}

	if swaggerInfo.Description != "This is a sample server for out Auth library. modified" {
		t.Errorf("expected Description to be This is a sample server for out Auth library. modified, got %s", swaggerInfo.Description)
	}

	// compare result to expected
	expected, err := os.ReadFile("testdata/oauth.expected.json")
	if err != nil {
		t.Fatal(err)
	}

	if swaggerInfo.SwaggerTemplate != string(expected) {
		t.Errorf("expected swagger template is different, got:\n%s", swaggerInfo.SwaggerTemplate)
	}
}
