package main

import (
	"flag"
	"seeds/email"

	"github.com/lsdch/biome/db"

	"github.com/edgedb/edgedb-go"
)

func main() {
	database := flag.String("db", "", "The name of the database to seed")
	flag.Parse()

	client := db.Connect(edgedb.Options{Database: *database})
	email.SetupEmailConfig(client, email.EmailSetupArgs{NoAuto: true})
}
