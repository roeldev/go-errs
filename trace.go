package errs

import (
	"errors"
	"fmt"
)

// Trace wraps an existing error with information about the stack frame its
// called from. Errors that implement the `ErrorWithStackTrace` interface add
// the frame to the existing stack trace. Other "simple" errors are wrapped
// in a `traceErr` type.
func Trace(err error) error {
	return TraceSkip(err, 1)
}

// TraceSkip, just like Trace(), wraps an existing error with information about
// the stack frame its called from. The stack frame is selected based on the
// skip argument.
func TraceSkip(err error, skip uint) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(ErrorWithFrames); ok {
		e.Frames().Capture(skip + 1)
		return err
	}

	return &traceErr{
		error:  err,
		frames: CaptureFrames(1, skip+2),
	}
}

// traceErr is a wrapper for "primitive" errors that do not have stack trace
// information. It does not contain an error message by itself and always
// displays the message of the underlying wrapped error.
type traceErr struct {
	error
	frames Frames
}

func (t *traceErr) Frames() *Frames { return &t.frames }

func (t *traceErr) Unwrap() error { return t.error }

func (t *traceErr) Is(target error) bool       { return errors.Is(t.error, target) }
func (t *traceErr) As(target interface{}) bool { return errors.As(t.error, target) }

func (t *traceErr) Format(s fmt.State, v rune) { FormatError(t, s, v) }
