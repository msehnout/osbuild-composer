// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// RequestJobJSONBody defines parameters for RequestJob.
type RequestJobJSONBody struct {
	Arch  string   `json:"arch"`
	Types []string `json:"types"`
}

// UpdateJobJSONBody defines parameters for UpdateJob.
type UpdateJobJSONBody struct {
	Result interface{} `json:"result"`
	Status string      `json:"status"`
}

// RequestJobRequestBody defines body for RequestJob for application/json ContentType.
type RequestJobJSONRequestBody RequestJobJSONBody

// UpdateJobRequestBody defines body for UpdateJob for application/json ContentType.
type UpdateJobJSONRequestBody UpdateJobJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Request a job
	// (POST /jobs)
	RequestJob(ctx echo.Context) error
	// Get running job
	// (GET /jobs/{token})
	GetJob(ctx echo.Context, token string) error
	// Update a running job
	// (PATCH /jobs/{token})
	UpdateJob(ctx echo.Context, token string) error
	// Upload an artifact
	// (PUT /jobs/{token}/artifacts/{name})
	UploadJobArtifact(ctx echo.Context, token string, name string) error
	// status
	// (GET /status)
	GetStatus(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// RequestJob converts echo context to params.
func (w *ServerInterfaceWrapper) RequestJob(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.RequestJob(ctx)
	return err
}

// GetJob converts echo context to params.
func (w *ServerInterfaceWrapper) GetJob(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "token" -------------
	var token string

	err = runtime.BindStyledParameter("simple", false, "token", ctx.Param("token"), &token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter token: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetJob(ctx, token)
	return err
}

// UpdateJob converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateJob(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "token" -------------
	var token string

	err = runtime.BindStyledParameter("simple", false, "token", ctx.Param("token"), &token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter token: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateJob(ctx, token)
	return err
}

// UploadJobArtifact converts echo context to params.
func (w *ServerInterfaceWrapper) UploadJobArtifact(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "token" -------------
	var token string

	err = runtime.BindStyledParameter("simple", false, "token", ctx.Param("token"), &token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter token: %s", err))
	}

	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameter("simple", false, "name", ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UploadJobArtifact(ctx, token, name)
	return err
}

// GetStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetStatus(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStatus(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/jobs", wrapper.RequestJob)
	router.GET("/jobs/:token", wrapper.GetJob)
	router.PATCH("/jobs/:token", wrapper.UpdateJob)
	router.PUT("/jobs/:token/artifacts/:name", wrapper.UploadJobArtifact)
	router.GET("/status", wrapper.GetStatus)

}
