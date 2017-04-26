package image

import (
	"fmt"
	"os"
)

func basicErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", e)
		os.Exit(-1)
	}
}
