package httpserver

import (
	"context"
	"net/http"
)

// GetRootMux exposes getRootMux to the external httpserver_test package.
func GetRootMux(ctx context.Context, httpServer *BasicHTTPServer) (*http.ServeMux, []string) {
	return httpServer.getRootMux(ctx)
}
