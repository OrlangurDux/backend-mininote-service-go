package middlewares

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var Env string

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {
	var param string
	if flag.Lookup("test.v") == nil {
		name := ".env"
		if Env == "development" {
			name = ".env-test"
		}
		err := godotenv.Load(name)
		if err != nil {
			err = godotenv.Load("../" + name)
			if err != nil {
				log.Fatalf("Error loading " + name + " file")
				os.Exit(1)
			}
		}
		param = os.Getenv(key)
	} else {
		path, _ := os.Getwd()
		err := godotenv.Load(strings.Split(path, "tests")[0] + ".env-test")
		if err != nil {
			log.Fatalf("Error loading .env.test file")
		}
		param = os.Getenv(key)
	}
	return param
}
