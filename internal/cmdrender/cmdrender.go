package cmdrender

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/henderiw/app-function-cmd/internal/fnruntime"
	"github.com/henderiw/app-function-cmd/internal/printer"
	"github.com/henderiw/app-function-cmd/internal/util/argutil"
	"github.com/henderiw/app-function-cmd/internal/util/cmdutil"
	"github.com/henderiw/app-function-cmd/internal/util/pathutil"
	"github.com/henderiw/app-function-cmd/internal/util/render"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

// NewRunner returns a command runner
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{ctx: ctx}
	c := &cobra.Command{
		Use:     "render [INPUT_FILE] [CONTAINER_IMAGE] [flags]",
		Short:   "render function",
		Long:    "render function",
		Example: "render function",
		RunE:    r.runE,
		PreRunE: r.preRunE,
	}
	c.Flags().StringVar(&r.resultsDirPath, "results-dir", "",
		"path to a directory to save function results")
	c.Flags().StringVarP(&r.dest, "output", "o", "",
		fmt.Sprintf("output resources are written to provided location. Allowed values: %s|%s|<OUT_DIR_PATH>", cmdutil.Stdout, cmdutil.Unwrap))
	c.Flags().StringVar(&r.imagePullPolicy, "image-pull-policy", string(fnruntime.IfNotPresentPull),
		fmt.Sprintf("pull image before running the container. It must be one of %s, %s and %s.", fnruntime.AlwaysPull, fnruntime.IfNotPresentPull, fnruntime.NeverPull))
	c.Flags().BoolVar(&r.allowExec, "allow-exec", false,
		"allow binary executable to be run during pipeline execution.")
	r.Command = c
	return r
}

// Runner contains the run function pipeline run command
type Runner struct {
	inputFile       string
	containerImage  string
	resultsDirPath  string
	imagePullPolicy string
	allowExec       bool
	dest            string
	Command         *cobra.Command
	ctx             context.Context
}

func (r *Runner) preRunE(c *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("inputFile and containerImage expected")
	} else {
		// resolve and validate the provided path
		r.inputFile = args[0]
		r.containerImage = args[1]
	}
	var err error
	r.inputFile, err = argutil.ResolveSymlink(r.ctx, r.inputFile)
	if err != nil {
		return err
	}
	if r.dest != "" && r.dest != cmdutil.Stdout && r.dest != cmdutil.Unwrap {
		if err := cmdutil.CheckDirectoryNotPresent(r.dest); err != nil {
			return err
		}
	}
	if r.resultsDirPath != "" {
		err := os.MkdirAll(r.resultsDirPath, 0755)
		if err != nil {
			return fmt.Errorf("cannot read or create results dir %q: %w", r.resultsDirPath, err)
		}
	}
	return cmdutil.ValidateImagePullPolicyValue(r.imagePullPolicy)
}

func (r *Runner) runE(c *cobra.Command, _ []string) error {
	var output io.Writer
	outContent := bytes.Buffer{}
	if r.dest != "" {
		// this means the output should be written to another destination
		// capture the content to be written
		output = &outContent
	}
	absInputFile, _, err := pathutil.ResolveAbsAndRelPaths(r.inputFile)
	if err != nil {
		return err
	}

	executor := render.Renderer{
		InputFile:       absInputFile,
		ContainerImage:  r.containerImage,
		ResultsDirPath:  r.resultsDirPath,
		Output:          output,
		ImagePullPolicy: cmdutil.StringToImagePullPolicy(r.imagePullPolicy),
		AllowExec:       r.allowExec,
		FileSystem:      filesys.FileSystemOrOnDisk{},
	}
	if err := executor.Execute(r.ctx); err != nil {
		return err
	}

	return cmdutil.WriteFnOutput(r.dest, outContent.String(), false, printer.FromContextOrDie(r.ctx).OutStream())
}
