package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Len(t, Filter(), 0)
	})
	t.Run("with nils", func(t *testing.T) {
		assert.Len(t, Filter(nil, nil), 0)
	})
	t.Run("with errors and nils", func(t *testing.T) {
		err1 := errors.New("some err")
		err2 := New("", "")

		f := Filter(err1, nil, nil, err2, nil)
		assert.Equal(t, []error{err1, err2}, f)
	})
}

func TestCombine(t *testing.T) {
	t.Run("with empty and nil", func(t *testing.T) {
		assert.Nil(t, Combine())
		assert.Nil(t, Combine(nil))
	})
	t.Run("with error", func(t *testing.T) {
		err := errors.New("some error")
		have := Combine(err)
		want := Trace(err).(*traceErr)
		want.frames = *GetFrames(have)

		assert.Exactly(t, want, have, "should add frame trace one single error")
	})
	t.Run("with errors", func(t *testing.T) {
		err1 := errors.New("first error")
		err2 := Newf(UnknownKind, "err with trace")
		multi := Combine(err1, err2).(*multiErr)

		assert.Exactly(t, []error{err1, err2}, multi.Errors())
	})
}

func TestAppend(t *testing.T) {
	t.Run("panic on nil dest ptr", func(t *testing.T) {
		assert.PanicsWithValue(t, panicAppendNilPtr, func() {
			Append(nil, New("foo", "bar"))
		})
	})
	t.Run("with nil", func(t *testing.T) {
		want := New("nice", "err")
		assert.Same(t, want, Append(&want, nil))
	})
	t.Run("with error", func(t *testing.T) {
		var have error
		want := errors.New("foobar")
		assert.Same(t, want, Append(&have, want))
		assert.Same(t, want, have)
	})
	t.Run("with errors", func(t *testing.T) {
		var have error
		list := []error{
			New("nice", "err"),
			errors.New("whoops"),
			fmt.Errorf("another %s", "error"),
		}

		Append(&have, list[0]) // set value to *have
		Append(&have, list[1]) // create multi error from errors 0 and 1
		Append(&have, list[2]) // append error 2 to multi error

		assert.IsType(t, new(multiErr), have)

		multi := have.(*multiErr)
		assert.Exactly(t, list, multi.Errors())
		assert.Contains(t, multi.Frames().String(), "multi_test.go:74")
		assert.Equal(t, len(multi.frames), 1)
	})
}
