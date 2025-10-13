package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jljl1337/xpense/internal/env"
)

func main() {
	env.SetConstants()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/health", env.Port))
	if err != nil || resp.StatusCode != 200 {
		os.Exit(1)
	}

	os.Exit(0)
}
