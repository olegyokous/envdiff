// Package watch — cmd.go wires Watch into the CLI.
package watch

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RunOptions holds CLI-level configuration for the watch command.
type RunOptions struct {
	Paths    []string
	Interval time.Duration
	Out      io.Writer
	OnChange func()
}

// Run starts the file watcher and blocks until SIGINT or SIGTERM.
// onChange is invoked each time a change is detected in Paths.
func Run(opts RunOptions) error {
	if len(opts.Paths) == 0 {
		return fmt.Errorf("watch: no files specified")
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.Interval <= 0 {
		opts.Interval = DefaultOptions().Interval
	}

	wOpts := Options{
		Interval: opts.Interval,
		Out:      opts.Out,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	go func() {
		<-sigs
		fmt.Fprintln(opts.Out, "\nwatch: shutting down")
		close(done)
	}()

	fmt.Fprintf(opts.Out, "watch: monitoring %d file(s) every %s\n", len(opts.Paths), opts.Interval)

	onChange := opts.OnChange
	if onChange == nil {
		onChange = func() {
			fmt.Fprintln(opts.Out, "watch: files changed")
		}
	}

	return Watch(opts.Paths, wOpts, onChange, done)
}
