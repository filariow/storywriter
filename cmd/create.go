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
	"fmt"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/filariow/storywriter/config"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new file starting from a template",
	Long: `Creates a new file starting from a template.

Configuration's Defaults are also injected leveraging
on Go Template notation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read configuration
		cfg, err := config.ParseDefault()
		if err != nil {
			fmt.Printf("error reading configuration: %v\n", err)
			return nil
		}

		// build output writer
		fw, op, err := buildOutputWriter(cmd, args, cfg)
		if err != nil {
			return err
		}

		// parse and execute template
		ifn := fmt.Sprintf("%s.md", args[0])
		if err := executeTemplate(fw, ifn, cfg); err != nil {
			return err
		}

		// print result
		if op != nil {
			fmt.Println(*op)
		}
		return nil
	},
	Args: cobra.ExactArgs(2),
}

func executeTemplate(fw io.Writer, ifn string, cfg *config.Config) error {
	ip := path.Join(cfg.Templates.Folder, ifn)

	t, err := template.New("").ParseFiles(ip)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(fw, ifn, cfg.Defaults)
}

func buildOutputWriter(cmd *cobra.Command, args []string, cfg *config.Config) (io.Writer, *string, error) {
	// return stdout, if stdout is required
	rstdout, err := cmd.Flags().GetBool("stdout")
	if err != nil {
		return nil, nil, err
	}
	if rstdout {
		return os.Stdout, nil, nil
	}

	// parse override flag
	override, err := cmd.Flags().GetBool("override")
	if err != nil {
		return nil, nil, err
	}

	// create file
	op := path.Join(cfg.Output.Folder, cfg.Output.Typed[args[0]].(string), fmt.Sprintf("%s.md", args[1]))

	// fail if file exists and override flag is not set
	if !override {
		_, err := os.Stat(op)
		if err == nil {
			return nil, nil, fmt.Errorf("file '%s' already exists: use --override to replace it", op)
		}
		if !os.IsNotExist(err) {
			return nil, nil, err
		}
	}

	// create file
	fw, err := os.Create(op)
	if err != nil {
		return nil, nil, err
	}

	return fw, &op, nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().Bool("stdout", false, "won't create a file, but print to standard output")
	createCmd.Flags().Bool("override", false, "will delete target file, if already exists")
}
