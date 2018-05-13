package main

import (
	"os"
)

func init() {
	if os.Getenv("TEST_LIVE") != "" {
		return
	}

	main()
}
