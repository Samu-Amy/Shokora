package main

import (
	"fmt"
	"log"

	"github.com/Samu-Amy/Shokora/internal/db"
	"github.com/Samu-Amy/Shokora/internal/env"
	"github.com/Samu-Amy/Shokora/internal/store"
)

func main() {
	env.LoadEnv() //! - Dev Only (use file .env) - !
	addr := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_USER", ""), env.GetString("POSTGRES_PASSWORD", ""), env.GetString("POSTGRES_DB", ""))

	db_conn, err := db.New(addr, 3, 3, "15m", true)
	if err != nil {
		log.Fatal(err)
	}

	defer db_conn.Close()

	store := store.NewPostgresStorage(db_conn)

	db.Seed(store, db_conn)
}
