package reporter

// TODO: Update this package so that successive calls to `NewWorkInProgress` cancel previous
// calls. In other words, No more than one `NewWorkInProgress`, without having to explicitly
// call `Stop`.

import (
	"fmt"
	"time"
)

type Reporter struct {
	quit chan any
}

func New() Reporter {
	return Reporter{make(chan any)}
}

func (r *Reporter) NewWorkInProgress(label string) {
	fmt.Print(label)
	go func() {
		defer fmt.Println()
		for {
			select {
			case <-r.quit:
				return
			default:
				fmt.Print(".")
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (r *Reporter) Stop() {
	r.quit <- struct{}{}
}
