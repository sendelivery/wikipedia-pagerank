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

func New() *Reporter {
	return &Reporter{make(chan any)}
}

func (r *Reporter) NewWorkInProgress(label string) {
	fmt.Print(label)
	go func() {
		defer fmt.Println()
		for i := 1; ; i++ {
			select {
			case <-r.quit:
				return
			default:
				if i%4 == 0 {
					// Remove the last 3 characters from the terminal.
					fmt.Print("\b\b\b   \b\b\b")
				} else {
					fmt.Print(".")
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (r *Reporter) Stop() {
	r.quit <- struct{}{}
}
