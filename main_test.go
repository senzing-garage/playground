package main

import (
	"os"
	"testing"
)

func TestMain(test *testing.T) {
	_ = test
	os.Args = []string{"command-name", "--avoid-serving"}

	main()
}
