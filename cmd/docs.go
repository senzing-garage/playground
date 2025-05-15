/*
 */
package cmd

import (
	"fmt"
	"io"
	"os"

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
			return fmt.Errorf("getting 'dir' value: %w", err)
		}
		if dir == "" {
			if dir, err = os.MkdirTemp("", "playground"); err != nil {
				return fmt.Errorf("constructing cobra.Command: %w", err)
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
		return fmt.Errorf("DocsAction: %w", err)
	}

	if _, err := fmt.Fprintf(out, "Documentation successfully created in %s\n", dir); err != nil {
		return fmt.Errorf("printing succsss: %w", err)
	}

	return nil
}
