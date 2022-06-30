package printer

import (
	"context"
	"fmt"
	"io"
	"os"
)

// TruncateOutput defines should output be truncated
var TruncateOutput bool

// Printer defines capabilities to display content in kpt CLI.
// The main intention, at the moment, is to abstract away printing
// output in the CLI so that we can evolve the kpt CLI UX.
type Printer interface {
	Printf(format string, args ...interface{})
	OutStream() io.Writer
	ErrStream() io.Writer
}

// New returns an instance of Printer.
func New(outStream, errStream io.Writer) Printer {
	if outStream == nil {
		outStream = os.Stdout
	}
	if errStream == nil {
		errStream = os.Stderr
	}
	return &printer{
		outStream: outStream,
		errStream: errStream,
	}
}

// printer implements default Printer to be used.
type printer struct {
	outStream io.Writer
	errStream io.Writer
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// printerKey is the context key for the printer.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const printerKey contextKey = 0

// OutStream returns the StdOut stream, this can be used by callers to print
// command output to stdout, do not print error/debug logs to this stream
func (pr *printer) OutStream() io.Writer {
	return pr.outStream
}

// ErrStream returns the StdErr stream, this can be used by callers to print
// command output to stderr, print only error/debug/info logs to this stream
func (pr *printer) ErrStream() io.Writer {
	return pr.errStream
}

// Printf is the wrapper over fmt.Printf that displays the output.
// this will print messages to stderr stream
func (pr *printer) Printf(format string, args ...interface{}) {
	fmt.Fprintf(pr.errStream, format, args...)
}

// FromContext returns printer instance associated with the context.
func FromContextOrDie(ctx context.Context) Printer {
	pr, ok := ctx.Value(printerKey).(Printer)
	if ok {
		return pr
	}
	panic("printer missing in context")
}

// WithContext creates new context from the given parent context
// by setting the printer instance.
func WithContext(ctx context.Context, pr Printer) context.Context {
	return context.WithValue(ctx, printerKey, pr)
}
