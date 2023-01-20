package middlewares

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Env -> switch config file cli --env variable
var Env string

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {
	var param string
	if flag.Lookup("test.v") == nil {
		name := ".env"
		if Env == "development" {
			name = ".env.test"
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
		config := strings.Split(path, "tests")[0] + ".env.test"
		err := godotenv.Load(config)
		fmt.Println(os.Getenv("TEST_VAR"))
		if err != nil {
			log.Println(err)
			log.Fatalf("Error loading %s file", config)
		}
		param = os.Getenv(key)
	}
	return param
}
