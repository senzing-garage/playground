package httpserver

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/docktermj/cloudshell/xtermservice"
	"github.com/pkg/browser"
	"github.com/senzing-garage/demo-entity-search/entitysearchservice"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	AvoidServing              bool
	EnableAll                 bool
	EnableEntitySearch        bool
	EnableJupyterLab          bool
	EnableXterm               bool
	EntitySearchRoutePrefix   string // IMPROVE: Only works with "entity-search"
	GrpcDialOptions           []grpc.DialOption
	GrpcTarget                string
	IsInDevelopment           bool
	JupyterLabRoutePrefix     string // IMPROVE: Only works with "jupyter"
	LogLevelName              string
	ObserverOrigin            string
	Observers                 []observer.Observer
	ReadHeaderTimeout         time.Duration
	SenzingInstanceName       string
	SenzingSettings           string
	SenzingVerboseLogging     int64
	ServerAddress             string
	ServerPort                int
	TtyOnly                   bool
	XtermAllowedHostnames     []string
	XtermArguments            []string
	XtermCommand              string
	XtermConnectionErrorLimit int
	XtermKeepalivePingTimeout int
	XtermMaxBufferSizeBytes   int
	XtermURLRoutePrefix       string // IMPROVE: Only works with "xterm"
}

type TemplateVariables struct {
	BasicHTTPServer
	EntitySearchStatus string
	EntitySearchURL    string
	HTMLTitle          string
	JupyterLabStatus   string
	JupyterLabURL      string
	XtermStatus        string
	XtermURL           string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method simply prints the 'Something' value in the type-struct.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	var err error

	rootMux, userMessages := httpServer.getRootMux(ctx)

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessages = append(userMessages, fmt.Sprintf("Starting server on interface:port '%s'...", listenOnAddress))

	for _, userMessage := range userMessages {
		outputln(userMessage)
	}

	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}

	// Start a web browser.  Unless disabled.

	if !httpServer.TtyOnly {
		_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", httpServer.ServerPort))
	}

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) getRootMux(ctx context.Context) (*http.ServeMux, []string) {
	var userMessages []string

	rootMux := http.NewServeMux()

	userMessages = append(userMessages, httpServer.addEntitySearchToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addJupiterToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addXtermToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addSiteToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addExamplesToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addStaticToMux(ctx, rootMux)...)

	return rootMux, userMessages
}

func (httpServer *BasicHTTPServer) addEntitySearchToMux(ctx context.Context, rootMux *http.ServeMux) []string {
	var result []string

	if httpServer.EnableAll || httpServer.EnableEntitySearch {
		entitySearchMux := httpServer.getEntitySearchMux(ctx)
		rootMux.Handle(
			fmt.Sprintf("/%s/", httpServer.EntitySearchRoutePrefix),
			http.StripPrefix("/entity-search", entitySearchMux),
		)

		result = append(result, fmt.Sprintf(
			"Serving Entity Search at    http://localhost:%d/%s",
			httpServer.ServerPort,
			httpServer.EntitySearchRoutePrefix))
	}

	return result
}

func (httpServer *BasicHTTPServer) addJupiterToMux(ctx context.Context, rootMux *http.ServeMux) []string {
	var result []string

	_ = ctx

	if httpServer.EnableAll || httpServer.EnableJupyterLab {
		proxy, err := newReverseProxy("http://localhost:8888")
		if err != nil {
			panic(err)
		}

		rootMux.HandleFunc(fmt.Sprintf("/%s/", httpServer.JupyterLabRoutePrefix), reverseProxyRequestHandler(proxy))
		result = append(result, fmt.Sprintf(
			"Serving JupyterLab at       http://localhost:%d/%s",
			httpServer.ServerPort,
			httpServer.JupyterLabRoutePrefix))
	}

	return result
}

func (httpServer *BasicHTTPServer) addXtermToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	if httpServer.EnableAll || httpServer.EnableXterm {
		err := os.Setenv("SENZING_ENGINE_CONFIGURATION_JSON", httpServer.SenzingSettings)
		if err != nil {
			panic(err)
		}

		xtermMux := httpServer.getXtermMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.XtermURLRoutePrefix), http.StripPrefix("/xterm", xtermMux))
		result = append(result, fmt.Sprintf(
			"Serving XTerm at            http://localhost:%d/%s",
			httpServer.ServerPort,
			httpServer.XtermURLRoutePrefix))
	}

	return result
}

func (httpServer *BasicHTTPServer) addSiteToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	_ = ctx

	rootMux.HandleFunc("/site/", httpServer.handleFuncForSite)
	result = append(result, fmt.Sprintf("Serving Console at          http://localhost:%d", httpServer.ServerPort))

	return result
}

func (httpServer *BasicHTTPServer) addStaticToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	result := []string{}

	_ = ctx

	rootDir, err := fs.Sub(httpServer.getStatic(), "static/root")
	if err != nil {
		panic(err)
	}

	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	return result
}

func (httpServer *BasicHTTPServer) addExamplesToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	result := []string{}
	_ = ctx

	rootMux.Handle("/examples/", http.StripPrefix("/examples", http.FileServer(http.Dir("/examples"))))

	return result
}

func (httpServer *BasicHTTPServer) getServerStatus(active bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}

	if active {
		result = "green"
	}

	return result
}

func (httpServer *BasicHTTPServer) getServerURL(active bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}

	if active {
		result = url
	}

	return result
}

func (httpServer *BasicHTTPServer) getStatic() fs.FS {
	if httpServer.IsInDevelopment {
		return os.DirFS("httpserver/")
	}

	return static
}

func (httpServer *BasicHTTPServer) populateStaticTemplate(
	responseWriter http.ResponseWriter,
	request *http.Request,
	filepath string,
	templateVariables TemplateVariables,
) {
	_ = request
	// templateBytes, err := static.ReadFile(filepath)
	templateBytes, err := fs.ReadFile(httpServer.getStatic(), filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error()) // #nosec G705 -- writing to stderr, not HTTP response
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error()) // #nosec G705 -- writing to stderr, not HTTP response
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	err = templateParsed.Execute(responseWriter, templateVariables)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error()) // #nosec G705 -- writing to stderr, not HTTP response
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *BasicHTTPServer) getEntitySearchMux(ctx context.Context) *http.ServeMux {
	service := &entitysearchservice.BasicHTTPService{}

	return service.Handler(ctx)
}

func (httpServer *BasicHTTPServer) getXtermMux(ctx context.Context) *http.ServeMux {
	xtermService := &xtermservice.XtermServiceImpl{
		AllowedHostnames:     httpServer.XtermAllowedHostnames,
		Arguments:            httpServer.XtermArguments,
		Command:              httpServer.XtermCommand,
		ConnectionErrorLimit: httpServer.XtermConnectionErrorLimit,
		KeepalivePingTimeout: httpServer.XtermKeepalivePingTimeout,
		MaxBufferSizeBytes:   httpServer.XtermMaxBufferSizeBytes,
		UrlRoutePrefix:       httpServer.XtermURLRoutePrefix,
	}

	return xtermService.Handler(ctx)
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *BasicHTTPServer) handleFuncForSite(writer http.ResponseWriter, request *http.Request) {
	templateVariables := TemplateVariables{
		BasicHTTPServer:    *httpServer,
		EntitySearchStatus: httpServer.getServerStatus(httpServer.EnableEntitySearch),
		EntitySearchURL: httpServer.getServerURL(
			httpServer.EnableEntitySearch,
			fmt.Sprintf("http://%s/entity-search", request.Host),
		),
		HTMLTitle:        "Senzing Quickstart",
		JupyterLabStatus: httpServer.getServerStatus(httpServer.EnableJupyterLab),
		JupyterLabURL: httpServer.getServerURL(
			httpServer.EnableJupyterLab,
			fmt.Sprintf("http://%s/jupyter", request.Host),
		),
		XtermStatus: httpServer.getServerStatus(httpServer.EnableXterm),
		XtermURL:    httpServer.getServerURL(httpServer.EnableXterm, fmt.Sprintf("http://%s/xterm", request.Host)),
	}

	writer.Header().Set("Content-Type", "text/html")

	filePath := "static/templates" + request.RequestURI
	httpServer.populateStaticTemplate(writer, request, filePath, templateVariables)
}

// newReverseProxy takes target host and creates a reverse proxy.
func newReverseProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, wraperror.Errorf(err, "url.Parse: %s", targetHost)
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// reverseProxyRequestHandler handles the http request using proxy.
func reverseProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r) // #nosec G704 -- proxy target is internally configured
	}
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}
