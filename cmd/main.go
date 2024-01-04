package main

import (
	"github.com/vadimpk/url-pinger/config"
	app "github.com/vadimpk/url-pinger/internal"
)

func main() {
	cfg := config.Get()

	app.Run(cfg)
}
