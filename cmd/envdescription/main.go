package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"messagio_assignment/internal/config"
)

func main() {
	var cfg config.Config

	help, err := cleanenv.GetDescription(&cfg, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(help)
}
