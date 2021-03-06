package errs

import (
	"sync"
)

// A WaitGroup is a collection of goroutines working on subtasks that are part
// of the same overall task. Unlike `errgroup.Group`, it collects possible
// errors returned from the subtasks and does not cancel the group when an
// error is encountered.
type WaitGroup struct {
	wg   sync.WaitGroup
	list List
}

// ErrList returns a List of collected errors from the called goroutines.
func (g *WaitGroup) ErrList() *List { return &g.list }

// Wait blocks until all function calls from the Go method have returned, then
// returns all collected errors as a combined (multi) error.
func (g *WaitGroup) Wait() error {
	g.wg.Wait()
	return g.list.Combine()
}

// Go calls the given function in a new goroutine. Errors from all calls are
// collected, combined and returned by Wait.
func (g *WaitGroup) Go(fn func() error) {
	g.wg.Add(1)

	go func() {
		// check of er een slice word aangemaakt in List
		// zoniet, dan dit in append/prepend doen
		g.list.Append(fn())
		g.wg.Done()
	}()
}
