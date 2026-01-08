// ! - Dev Only (use file .env) - !
package env

import "github.com/joho/godotenv"

func LoadEnv() {
	godotenv.Load("../.env")
}
