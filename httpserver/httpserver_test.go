package httpserver_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
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
		APIUrlRoutePrefix:        "api",
		AvoidServing:             true,
		EnableAll:                true,
		EntitySearchRoutePrefix:  "entity-search",
		JupyterLabRoutePrefix:    "jupyter",
		LogLevelName:             logLevelName,
		ObserverOrigin:           "Test Observer origin",
		Observers:                []observer.Observer{observer1},
		OpenAPISpecificationRest: senzingrestservice.OpenAPISpecificationJSON,
		ReadHeaderTimeout:        10 * time.Second,
		SenzingInstanceName:      "Test HTTP Server",
		SenzingSettings:          senzingSettings,
		SwaggerURLRoutePrefix:    "swagger",
		TtyOnly:                  true,
		XtermURLRoutePrefix:      "xterm",
	}

	return result
}
