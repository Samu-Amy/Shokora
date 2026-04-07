package main

import (
	"database/sql"
)

// Clear db
func clearTestDB(db *sql.DB) {
	tables := []string{
		"users",
		"user_settings",
		"user_sessions",
		"refresh_tokens",
		"reset_session_tokens",
		"verification_tokens",
		"products",
	}

	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
		if err != nil {
			panic("failed to clean table " + table + ": " + err.Error())
		}
	}
}
