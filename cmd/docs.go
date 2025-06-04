/*
 */
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// DocsCmd represents the docs command.
var DocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for the command",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = args
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return wraperror.Errorf(err, "getting 'dir' value")
		}
		if dir == "" {
			if dir, err = os.MkdirTemp("", "playground"); err != nil {
				return wraperror.Errorf(err, "constructing cobra.Command")
			}
		}

		return DocsAction(os.Stdout, dir)
	},
}

func init() {
	RootCmd.AddCommand(DocsCmd)
	DocsCmd.Flags().StringP("dir", "d", "", "Destination directory for docs")
}

func DocsAction(out io.Writer, dir string) error {
	if err := doc.GenMarkdownTree(RootCmd, dir); err != nil {
		return wraperror.Errorf(err, "DocsAction")
	}

	if _, err := fmt.Fprintf(out, "Documentation successfully created in %s\n", dir); err != nil {
		return wraperror.Errorf(err, "printing succsss")
	}

	return nil
}
