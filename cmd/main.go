package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"krakjam2022_scoreboard/pkg/database"
	"krakjam2022_scoreboard/pkg/rest"
	"log"
	"os"
)

func run() error {
	pg, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_CONN_STRING")), &gorm.Config{})
	if err != nil {
		return err
	}

	db, err := database.NewDB(pg)
	if err != nil {
		return err
	}

	srv := rest.New(db, os.Getenv("SECRET"))
	log.Fatalln(srv.Run(os.Getenv("ADDR")))

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
