package render

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/henderiw/app-function-cmd/internal/fnruntime"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/fn/runtime/runtimeutil"
)

// Renderer hydrates a given pkg by running the functions in the input pipeline
type Renderer struct {
	// InputFile is the absolute path to the input file
	InputFile string

	// Container Image that runs the apptest function
	ContainerImage string

	// Runtime knows how to pick a function runner for a given function
	//Runtime fn.FunctionRuntime

	// ResultsDirPath is absolute path to the directory to write results
	ResultsDirPath string

	// fnResultsList is the list of results from the pipeline execution
	//fnResultsList *fnresult.ResultList

	// Output is the writer to which the output resources are written
	Output io.Writer

	// ImagePullPolicy pull image before running the container.
	// It must be one of fnruntime.AlwaysPull, fnruntime.IfNotPresentPull, fnruntime.NeverPull
	ImagePullPolicy fnruntime.ImagePullPolicy

	// AllowExec allow binary executable to be run during pipeline execution
	AllowExec bool

	// FileSystem is the input filesystem to operate on
	FileSystem filesys.FileSystem
}

// Execute runs a pipeline.
func (e *Renderer) Execute(ctx context.Context) error {
	fmt.Printf("InputFile: %s, ContainerImage: %s Output: %s\n", e.InputFile, e.ContainerImage, e.Output)

	cfn := &fnruntime.ContainerFn{
		Ctx:             ctx,
		InputFile:       e.InputFile,
		Image:           e.ContainerImage,
		ImagePullPolicy: e.ImagePullPolicy,
		//Timeout: ,
		Perm: fnruntime.ContainerFnPermission{},
		//UIDGID: ,
		StorageMounts: []runtimeutil.StorageMount{},
		Env:           []string{},
	}

	r := strings.NewReader(e.InputFile)
	

	if err := cfn.Run(r, e.Output); err != nil {
		fmt.Printf("cfn run error: %s\n", err)
	}
	return nil
}
