// ! - Dev Only (use file .env) - !
package env

import "github.com/joho/godotenv"

func LoadDevEnv() {
	godotenv.Load("../.env")
}

func LoadTestEnv() {
	godotenv.Load("../../.env")
}
