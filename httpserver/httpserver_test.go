package httpserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/playground/httpserver"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_Serve(test *testing.T) {
	ctx := test.Context()
	httpServer := getTestObject(ctx, test)
	err := httpServer.Serve(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test private methods
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_getRootMux(test *testing.T) {
	ctx := test.Context()
	httpServer := getTestObject(ctx, test)
	rootMux, _ := httpserver.GetRootMux(ctx, httpServer)

	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{name: "Console", path: "/", expectedCode: http.StatusOK},
		{name: "Entity Search", path: "/entity-search/", expectedCode: http.StatusOK},
		{name: "Removed Senzing REST API", path: "/api/heartbeat", expectedCode: http.StatusNotFound},
		{
			name:         "Removed Senzing REST API proxy",
			path:         "/entity-search/api/heartbeat",
			expectedCode: http.StatusNotFound,
		},
		{name: "Removed Swagger UI", path: "/swagger/", expectedCode: http.StatusNotFound},
	}

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			request := httptest.NewRequestWithContext(ctx, http.MethodGet, testCase.path, nil)
			recorder := httptest.NewRecorder()
			rootMux.ServeHTTP(recorder, request)
			require.Equal(test, testCase.expectedCode, recorder.Code)
		})
	}
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *httpserver.BasicHTTPServer {
	t.Helper()

	_ = ctx

	observer1 := &observer.NullObserver{
		ID: "Observer 1",
	}

	logLevelName := "INFO"
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")

	if len(osenvLogLevel) > 0 {
		logLevelName = osenvLogLevel
	}

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &httpserver.BasicHTTPServer{
		AvoidServing:            true,
		EnableAll:               true,
		EntitySearchRoutePrefix: "entity-search",
		JupyterLabRoutePrefix:   "jupyter",
		LogLevelName:            logLevelName,
		ObserverOrigin:          "Test Observer origin",
		Observers:               []observer.Observer{observer1},
		ReadHeaderTimeout:       10 * time.Second,
		SenzingInstanceName:     "Test HTTP Server",
		SenzingSettings:         senzingSettings,
		TtyOnly:                 true,
		XtermURLRoutePrefix:     "xterm",
	}

	return result
}
