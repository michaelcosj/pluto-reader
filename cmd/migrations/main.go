package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/michaelcosj/pluto-reader/db/migrations"
	"github.com/pressly/goose/v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("pgx", os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("unable to connect to database: %w", err)
	}

	goose.SetBaseFS(migrations.Embed)

	if err := goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	flag.Parse()
	args := flag.Args()

	switch args[0] {
	case "up":
		log.Printf("running up migration\n")
		if err := goose.Up(db, "."); err != nil {
			panic(err)
		}
	case "down":
		log.Printf("running down migration\n")
		if err := goose.Down(db, "."); err != nil {
			panic(err)
		}
	case "redo":
		log.Printf("running redo migration\n")
		if err := goose.Redo(db, "."); err != nil {
			panic(err)
		}
	case "reset":
		log.Printf("running reset migration\n")
		if err := goose.Reset(db, "."); err != nil {
			panic(err)
		}
    case "status":
		if err := goose.Status(db, "."); err != nil {
			panic(err)
		}
	default:
		log.Fatal("invalid migration argument: must be one of 'up', 'down', 'redo' and 'reset'\n")
	}

}
