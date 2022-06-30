package argutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/henderiw/app-function-cmd/internal/printer"
)

// ResolveSymlink returns the resolved symlink path for the input path
func ResolveSymlink(ctx context.Context, path string) (string, error) {
	isSymlink := false
	f, err := os.Lstat(path)
	if err == nil {
		// this step only helps with printing WARN message by checking if the input
		// path has symlink, so do not error out at this phase and let
		// filepath.EvalSymlinks(path) handle the cases
		if f.Mode().Type() == os.ModeSymlink {
			isSymlink = true
		}
	}
	rp, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	if isSymlink {
		fmt.Fprintf(printer.FromContextOrDie(ctx).ErrStream(), "[WARN] resolved symlink %q to %q, please note that the symlinks within the package are ignored\n", path, rp)
	}
	return rp, nil
}
