package commands

import (
	"context"

	"github.com/spf13/cobra"
)

// GetKptCommands returns the set of kpt commands to be registered
func GetAppTestCommands(ctx context.Context, name, version string) []*cobra.Command {
	var c []*cobra.Command
	fnCmd := GetFnCommand(ctx, name)

	c = append(c, fnCmd)

	return c
}
