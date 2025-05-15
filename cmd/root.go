/*
 */
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/option/optiontype"
	"github.com/senzing-garage/go-cmdhelping/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"github.com/senzing-garage/playground/httpserver"
	"github.com/senzing-garage/serve-grpc/grpcserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "HTTP/gRPC server supporting various services"
	Use   string = "playground"
	Long  string = `
A server supporting the following services:
    - HTTP: Senzing API server
    - HTTP: Swagger UI
    - HTTP: Xterm
    - gRPC:
    `
)

var isInDevelopment = option.ContextVariable{
	Arg:     "is-in-development",
	Default: option.OsLookupEnvBool("SENZING_TOOLS_IS_IN_DEVELOPMENT", false),
	Envar:   "SENZING_TOOLS_IS_IN_DEVELOPMENT",
	Help:    "For testing only. [%s]",
	Type:    optiontype.Bool,
}

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	isInDevelopment,
	option.AvoidServe,
	option.Configuration,
	option.DatabaseURL,
	option.EngineInstanceName,
	option.EngineLogLevel,
	option.EngineSettings,
	option.GrpcPort,
	option.HTTPPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverURL,
	option.ServerAddress,
	option.TtyOnly,
	option.XtermAllowedHostnames.SetDefault(getDefaultAllowedHostnames()),
	option.XtermArguments,
	option.XtermCommand,
	option.XtermConnectionErrorLimit,
	option.XtermKeepalivePingTimeout,
	option.XtermMaxBufferSizeBytes,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command.
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command.
func RunE(_ *cobra.Command, _ []string) error {
	var err error
	ctx := context.Background()

	// Set default value for SENZING_TOOLS_DATABASE_URL.

	_, isSet := os.LookupEnv("SENZING_TOOLS_DATABASE_URL")
	if !isSet {
		err = os.Setenv("SENZING_TOOLS_DATABASE_URL", SenzingToolsDatabaseURL)
		if err != nil {
			return err
		}
	}

	// Build configuration for Senzing engine.

	senzingSettings, err := settings.BuildAndVerifySettings(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Build observers.

	observers := []observer.Observer{}

	// Setup gRPC server

	grpcserver := &grpcserver.BasicGrpcServer{
		AvoidServing:          viper.GetBool(option.AvoidServe.Arg),
		EnableAll:             true,
		EnableSzConfig:        viper.GetBool(option.EnableSzConfig.Arg),
		EnableSzConfigManager: viper.GetBool(option.EnableSzConfigManager.Arg),
		EnableSzDiagnostic:    viper.GetBool(option.EnableSzDiagnostic.Arg),
		EnableSzEngine:        viper.GetBool(option.EnableSzEngine.Arg),
		EnableSzProduct:       viper.GetBool(option.EnableSzProduct.Arg),
		LogLevelName:          viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:        viper.GetString(option.ObserverOrigin.Arg),
		ObserverURL:           viper.GetString(option.ObserverURL.Arg),
		Port:                  viper.GetInt(option.GrpcPort.Arg),
		SenzingSettings:       senzingSettings,
		SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
		SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
	}

	// Create object and Serve.

	httpServer := &httpserver.BasicHTTPServer{
		APIUrlRoutePrefix:         "api",
		AvoidServing:              viper.GetBool(option.AvoidServe.Arg),
		EnableAll:                 true,
		EntitySearchRoutePrefix:   "entity-search",
		IsInDevelopment:           viper.GetBool(isInDevelopment.Arg),
		JupyterLabRoutePrefix:     "jupyter",
		LogLevelName:              viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:            viper.GetString(option.ObserverOrigin.Arg),
		Observers:                 observers,
		OpenAPISpecificationRest:  senzingrestservice.OpenAPISpecificationJSON,
		ReadHeaderTimeout:         60 * time.Second,
		SenzingInstanceName:       viper.GetString(option.EngineInstanceName.Arg),
		SenzingSettings:           senzingSettings,
		SenzingVerboseLogging:     viper.GetInt64(option.EngineLogLevel.Arg),
		ServerAddress:             viper.GetString(option.ServerAddress.Arg),
		ServerPort:                viper.GetInt(option.HTTPPort.Arg),
		SwaggerURLRoutePrefix:     "swagger",
		TtyOnly:                   viper.GetBool(option.TtyOnly.Arg),
		XtermAllowedHostnames:     viper.GetStringSlice(option.XtermAllowedHostnames.Arg),
		XtermArguments:            viper.GetStringSlice(option.XtermArguments.Arg),
		XtermCommand:              viper.GetString(option.XtermCommand.Arg),
		XtermConnectionErrorLimit: viper.GetInt(option.XtermConnectionErrorLimit.Arg),
		XtermKeepalivePingTimeout: viper.GetInt(option.XtermKeepalivePingTimeout.Arg),
		XtermMaxBufferSizeBytes:   viper.GetInt(option.XtermMaxBufferSizeBytes.Arg),
		XtermURLRoutePrefix:       "xterm",
	}

	// Start servers.

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		err = httpServer.Serve(ctx)
		if err != nil {
			fmt.Printf("Error: httpServer - %v\n", err)
		}
	}()

	go func() {
		defer waitGroup.Done()
		err = grpcserver.Serve(ctx)
		if err != nil {
			fmt.Printf("Error: grpcServer - %v\n", err)
		}
	}()

	waitGroup.Wait()

	return nil
}

// Used in construction of cobra.Command.
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// --- Networking -------------------------------------------------------------

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func getDefaultAllowedHostnames() []string {
	result := []string{"localhost"}
	outboundIPAddress := getOutboundIP().String()
	if len(outboundIPAddress) > 0 {
		result = append(result, outboundIPAddress)
	}
	return result
}
