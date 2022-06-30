package run

import (
	"context"

	"github.com/henderiw/app-function-cmd/cmd/commands"
	"github.com/henderiw/app-function-cmd/internal/printer"
	"github.com/spf13/cobra"
)

var version = "unknown"

func GetMain(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "apptest",
		Short:        "apptest for testing function programs in apps",
		Long:         "apptest for testing function programs in apps",
		SilenceUsage: true,
		// We handle all errors in main after return from cobra so we can
		// adjust the error message coming from libraries
		SilenceErrors: true,
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
	cmd.PersistentFlags().BoolVar(&printer.TruncateOutput, "truncate-output", true,
		"Enable the truncation for output")
	// wire the global printer
	pr := printer.New(cmd.OutOrStdout(), cmd.ErrOrStderr())

	// create context with associated printer
	ctx = printer.WithContext(ctx, pr)

	cmd.AddCommand(commands.GetAppTestCommands(ctx, "apptest", version)...)
	return cmd
}
