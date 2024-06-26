/*
Copyright © 2024 Francesco Ilario

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
  "os"
	"os/exec"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Opens an editor to edit the story",
  Long: `Opens an editor to edit the story.

It uses the EDITOR environment variable to determine the editor to use.
If not found, 'xdg-open' is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
    ss := func() []string {
      if editor, ok := os.LookupEnv("EDITOR"); ok {
        return []string{editor, args[0]}
      }
      return []string{"xdg-open", args[0]}
    }()

    c := exec.CommandContext(cmd.Context(), ss[0], ss[1:]...)
    c.Stdin = os.Stdin
    c.Stderr = os.Stderr
    c.Stdout = os.Stdout
    return c.Run()
	},
  Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(editCmd)
}
