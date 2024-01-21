package main

import (
	"flag"
	"fmt"
	blogservice "homework/internal/blog_service"
	"homework/internal/config"

	"github.com/go-playground/validator/v10"
)

func main() {
	configPath := flag.String("cfg", "", "path to file config")
	flag.Parse()

	fmt.Println(*configPath)

	if *configPath == "" {
		panic("nil config file")
	}

	validator := validator.New()

	cfg := config.MustConfig(*configPath, validator)
	blogservice.Start(cfg, validator)
}
