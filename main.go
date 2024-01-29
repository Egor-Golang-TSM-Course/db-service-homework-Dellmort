package main

import (
	"flag"
	"fmt"
	blogservice "homework/internal/blog_service"
	"homework/internal/config"

	"github.com/go-playground/validator/v10"
)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDYxNTA2MzQsImlkIjozfQ.eu-7Fcbgnedj6Hp2bUTmR47uNd269DsuyBj7t1578vI

// TODO: add migrations
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
