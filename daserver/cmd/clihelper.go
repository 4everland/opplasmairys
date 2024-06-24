package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
)

// DefaultInterruptSignals is a set of default interrupt signals.
var DefaultInterruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

type interruptContextKeyType struct{}

var blockerContextKey = interruptContextKeyType{}

type interruptCatcher struct {
	incoming chan os.Signal
}

// Block blocks until either an interrupt signal is received, or the context is cancelled.
// No error is returned on interrupt.
func (c *interruptCatcher) Block(ctx context.Context) {
	select {
	case <-c.incoming:
	case <-ctx.Done():
	}
}

// BlockFn simply blocks until the implementation of theBlockFn blocker interrupts it, or till the given context is cancelled.
type BlockFn func(ctx context.Context)

// WithInterruptBlocker attaches an interrupt handler to the context,
// which continues to receive signals after every block.
// This helps functions block on individual consecutive interrupts.
func WithInterruptBlocker(ctx context.Context) context.Context {
	if ctx.Value(blockerContextKey) != nil { // already has an interrupt handler
		return ctx
	}
	catcher := &interruptCatcher{
		incoming: make(chan os.Signal, 10),
	}
	signal.Notify(catcher.incoming, DefaultInterruptSignals...)

	return context.WithValue(ctx, blockerContextKey, BlockFn(catcher.Block))
}

// BlockOnInterrupts blocks until a SIGTERM is received.
// Passing in signals will override the default signals.
func BlockOnInterrupts(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultInterruptSignals
	}
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signals...)
	<-interruptChannel
}

func FormatVersion(version string, gitCommit string, gitDate string, meta string) string {
	v := version
	if gitCommit != "" {
		if len(gitCommit) >= 8 {
			v += "-" + gitCommit[:8]
		} else {
			v += "-" + gitCommit
		}
	}
	if gitDate != "" {
		v += "-" + gitDate
	}
	if meta != "" {
		v += "-" + meta
	}
	return v
}

type CloneableGeneric interface {
	cli.Generic
	Clone() any
}

// ProtectFlags ensures that no flags are safe to Apply() flag sets to without accidental flag-value mutations.
// ProtectFlags panics if any of the flag definitions cannot be protected.
func ProtectFlags(flags []cli.Flag) []cli.Flag {
	out := make([]cli.Flag, 0, len(flags))
	for _, f := range flags {
		fCopy, err := cloneFlag(f)
		if err != nil {
			panic(fmt.Errorf("failed to clone flag %q: %w", f.Names()[0], err))
		}
		out = append(out, fCopy)
	}
	return out
}

func cloneFlag(f cli.Flag) (cli.Flag, error) {
	switch typedFlag := f.(type) {
	case *cli.GenericFlag:
		// We have to clone Generic, since it's an interface,
		// and setting it causes the next use of the flag to have a different default value.
		if genValue, ok := typedFlag.Value.(CloneableGeneric); ok {
			cpy := *typedFlag
			cpyVal, ok := genValue.Clone().(cli.Generic)
			if !ok {
				return nil, fmt.Errorf("cloned Generic value is not Generic: %T", typedFlag)
			}
			cpy.Value = cpyVal
			return &cpy, nil
		} else {
			return nil, fmt.Errorf("cannot clone Generic value: %T", typedFlag)
		}
	default:
		// Other flag types are safe to re-use, although not strictly safe for concurrent use.
		// urfave v3 hopefully fixes this.
		return f, nil
	}
}
