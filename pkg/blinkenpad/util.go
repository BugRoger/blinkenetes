package blinkenpad

import (
	"fmt"
	"os"
)

func (b *Blinkenpad) handleError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("An error occured: %v\n", err))
		b.Stop()
		os.Exit(-1)
	}
}
