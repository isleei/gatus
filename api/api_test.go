package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/ui"
	"github.com/TwiN/gatus/v5/security"
	"github.com/gofiber/fiber/v2"
)

func TestNew(t *testing.T) {
	type Scenario struct {
		Name         string
		Method       string
		Path         string
		ExpectedCode int
		Gzip         bool
		WithSecurity bool
	}
	scenarios := []Scenario{
		{
			Name:         "health",
			Method:       http.MethodGet,
			Path:         "/health",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "custom.css",
			Method:       http.MethodGet,
			Path:         "/css/custom.css",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "custom.css-gzipped",
			Method:       http.MethodGet,
			Path:         "/css/custom.css",
			ExpectedCode: fiber.StatusOK,
			Gzip:         true,
		},
		{
			Name:         "metrics",
			Method:       http.MethodGet,
			Path:         "/metrics",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "favicon.ico",
			Method:       http.MethodGet,
			Path:         "/favicon.ico",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "app.js",
			Method:       http.MethodGet,
			Path:         "/js/app.js",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "app.js-gzipped",
			Method:       http.MethodGet,
			Path:         "/js/app.js",
			ExpectedCode: fiber.StatusOK,
			Gzip:         true,
		},
		{
			Name:         "chunk-vendors.js",
			Method:       http.MethodGet,
			Path:         "/js/chunk-vendors.js",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "chunk-vendors.js-gzipped",
			Method:       http.MethodGet,
			Path:         "/js/chunk-vendors.js",
			ExpectedCode: fiber.StatusOK,
			Gzip:         true,
		},
		{
			Name:         "index",
			Method:       http.MethodGet,
			Path:         "/",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "admin",
			Method:       http.MethodGet,
			Path:         "/admin",
			ExpectedCode: fiber.StatusOK,
		},
		{
			Name:         "index-html-redirect",
			Method:       http.MethodGet,
			Path:         "/index.html",
			ExpectedCode: fiber.StatusMovedPermanently,
		},
		{
			Name:         "index-should-return-200-even-if-not-authenticated",
			Method:       http.MethodGet,
			Path:         "/",
			ExpectedCode: fiber.StatusOK,
			WithSecurity: true,
		},
		{
			Name:         "admin-should-return-401-if-not-authenticated",
			Method:       http.MethodGet,
			Path:         "/admin",
			ExpectedCode: fiber.StatusUnauthorized,
			WithSecurity: true,
		},
		{
			Name:         "endpoints-should-return-401-if-not-authenticated",
			Method:       http.MethodGet,
			Path:         "/api/v1/endpoints/statuses",
			ExpectedCode: fiber.StatusUnauthorized,
			WithSecurity: true,
		},
		{
			Name:         "admin-notifications-should-return-401-if-not-authenticated",
			Method:       http.MethodGet,
			Path:         "/api/v1/admin/notifications",
			ExpectedCode: fiber.StatusUnauthorized,
			WithSecurity: true,
		},
		{
			Name:         "admin-reload-should-return-401-if-not-authenticated",
			Method:       http.MethodPost,
			Path:         "/api/v1/admin/reload",
			ExpectedCode: fiber.StatusUnauthorized,
			WithSecurity: true,
		},
		{
			Name:         "config-should-return-200-even-if-not-authenticated",
			Method:       http.MethodGet,
			Path:         "/api/v1/config",
			ExpectedCode: fiber.StatusOK,
			WithSecurity: true,
		},
		{
			Name:         "admin-notifications-should-return-200-without-security",
			Method:       http.MethodGet,
			Path:         "/api/v1/admin/notifications",
			ExpectedCode: fiber.StatusOK,
			WithSecurity: false,
		},
		{
			Name:         "admin-reload-should-return-202-without-security",
			Method:       http.MethodPost,
			Path:         "/api/v1/admin/reload",
			ExpectedCode: fiber.StatusAccepted,
			WithSecurity: false,
		},
		{
			Name:         "config-should-always-return-200",
			Method:       http.MethodGet,
			Path:         "/api/v1/config",
			ExpectedCode: fiber.StatusOK,
			WithSecurity: false,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			cfg := &config.Config{Metrics: true, UI: &ui.Config{}}
			if scenario.WithSecurity {
				cfg.Security = &security.Config{
					Basic: &security.BasicConfig{
						Username:                        "john.doe",
						PasswordBcryptHashBase64Encoded: "JDJhJDA4JDFoRnpPY1hnaFl1OC9ISlFsa21VS09wOGlPU1ZOTDlHZG1qeTFvb3dIckRBUnlHUmNIRWlT",
					},
				}
			}
			api := New(cfg)
			router := api.Router()
			method := scenario.Method
			if len(method) == 0 {
				method = http.MethodGet
			}
			request := httptest.NewRequest(method, scenario.Path, http.NoBody)
			if scenario.Gzip {
				request.Header.Set("Accept-Encoding", "gzip")
			}
			response, err := router.Test(request)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != scenario.ExpectedCode {
				t.Errorf("%s %s should have returned %d, but returned %d instead", request.Method, request.URL, scenario.ExpectedCode, response.StatusCode)
			}
		})
	}
}
