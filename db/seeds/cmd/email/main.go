package main

import (
	"flag"
	"seeds/email"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/lsdch/biome/db"
)

func main() {
	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(gelcfg.Options{Database: *database})
	email.SetupEmailConfig(client, email.EmailSetupArgs{NoAuto: true})
}
