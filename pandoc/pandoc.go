package pandoc

import (
	"context"
	"os"
	"os/exec"
)

const (
	pandocCommand string = "pandoc"
)

func Run(ctx context.Context, filename, oformat string) error {
	cmd := exec.CommandContext(ctx, pandocCommand, "-t", oformat, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
