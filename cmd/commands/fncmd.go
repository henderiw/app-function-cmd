package commands

import (
	"context"

	"github.com/henderiw/app-function-cmd/internal/cmdrender"
	"github.com/spf13/cobra"
)

func GetFnCommand(ctx context.Context, name string) *cobra.Command {
	functions := &cobra.Command{
		Use:     "fn",
		Short:   "execute functions",
		Long:    "execute functions",
		Aliases: []string{"functions"},
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := cmd.Flags().GetBool("help")
			if err != nil {
				return err
			}
			if h {
				return cmd.Help()
			}
			return cmd.Usage()
		},
	}

	functions.AddCommand(
		cmdrender.NewCommand(ctx, name),
	)
	return functions
}
