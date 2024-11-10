package reporter

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
			time.Sleep(500 * time.Millisecond)
			select {
			case <-r.quit:
				return
			default:
				fmt.Print(".")
			}
		}
	}()
}

func (r *Reporter) Stop() {
	r.quit <- struct{}{}
}
