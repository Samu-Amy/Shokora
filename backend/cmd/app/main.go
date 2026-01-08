package main

import (
	"log"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/env"
)

//TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato (?)

// DB Connection string
// connStr := "user=${DEV_POSTGRES_USER} dbname=${DEV_POSTGRES_DB} password=${DEV_POSTGRES_PASSWORD} host=localhost port=5432 sslmode=disable" // TODO: usa ssl in prod ("verify-full"?) (?)

func main() {
	env.LoadEnv() //! - Dev Only (use file .env) - !

	config := api.Config{
		Addr: env.GetString("SERVER_PORT", ":8080"),
	}

	app := api.NewApp(config)

	err := app.Run()

	if err != nil {
		log.Println(err) // TODO: sistema
	}
}
