package main

import (
	"context"
	"flag"
	"os"

	"github.com/henderiw/app-function-cmd/cmd/run"
	"github.com/spf13/cobra"
	"k8s.io/component-base/cli"
	"k8s.io/klog"
	k8scmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func main() {
	// Handle all setup in the runMain function so os.Exit doesn't interfere
	// with defer.
	os.Exit(runMain())
}

// runMain does the initial setup in order to run kpt. The return value from
// this function will be the exit code when kpt terminates.
func runMain() int {
	var logFlags flag.FlagSet
	var err error

	ctx := context.Background()

	// runs the command
	cmd := run.GetMain(ctx)

	// Enable commandline flags for klog.
	// logging will help in collecting debugging information from users
	klog.InitFlags(&logFlags)
	// By default klog v1 logs to stderr, switch that off
	_ = logFlags.Set("logtostderr", "false")
	_ = logFlags.Set("alsologtostderr", "true")
	cmd.Flags().AddGoFlagSet(&logFlags)

	err = cli.RunNoErrOutput(cmd)
	if err != nil {
		return handleErr(cmd, err)
	}
	return 0
}

// handleErr takes care of printing an error message for a given error.
func handleErr(cmd *cobra.Command, err error) int {

	// Finally just let the error handler for kubectl handle it. This handles
	// printing of several error types used in kubectl
	// TODO: See if we can handle this in kpt and get a uniform experience
	// across all of kpt.
	k8scmdutil.CheckErr(err)
	return 1
}
