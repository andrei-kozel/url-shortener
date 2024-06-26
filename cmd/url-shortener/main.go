package main

import (
	"fmt"

	"github.com/andrei-kozel/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Print(cfg)
}
