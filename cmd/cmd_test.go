package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving", "--tty-only"}
	Execute()
}

func Test_Execute_completion(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "completion"}
	Execute()
}

func Test_Execute_docs(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "docs"}
	Execute()
}

func Test_Execute_help(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--help"}
	Execute()
}

func Test_PreRun(test *testing.T) {
	_ = test
	args := []string{"command-name", "--help"}
	PreRun(RootCmd, args)
}

func Test_RunE(test *testing.T) {
	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")
	err := RunE(RootCmd, []string{})
	require.NoError(test, err)
}

func Test_RunE_badGrpcURL(test *testing.T) {
	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")
	test.Setenv("SENZING_TOOLS_GRPC_URL", "grpc://bad")
	err := RunE(RootCmd, []string{})
	require.NoError(test, err)
}

func Test_RootCmd(test *testing.T) {
	_ = test
	err := RootCmd.Execute()
	require.NoError(test, err)
	err = RootCmd.RunE(RootCmd, []string{})
	require.NoError(test, err)
}

func Test_completionCmd(test *testing.T) {
	_ = test
	err := CompletionCmd.Execute()
	require.NoError(test, err)
	err = CompletionCmd.RunE(CompletionCmd, []string{})
	require.NoError(test, err)
}

func Test_docsCmd(test *testing.T) {
	_ = test
	err := DocsCmd.Execute()
	require.NoError(test, err)
	err = DocsCmd.RunE(DocsCmd, []string{})
	require.NoError(test, err)
}

func Test_docsAction_badDir(test *testing.T) {
	var buffer bytes.Buffer
	badDir := "/tmp/no/directory/exists"
	err := DocsAction(&buffer, badDir)
	require.Error(test, err)
}
