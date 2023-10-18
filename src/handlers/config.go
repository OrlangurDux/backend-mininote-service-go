package middlewares

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

// Env -> switch config file cli --env variable
var Env string

func LoadDotEnv() {
	if flag.Lookup("test.v") == nil {
		name := ".env"
		if Env == "development" {
			name = ".env.test"
		}
		err := godotenv.Load(name)
		if err != nil {
			err = godotenv.Load("../" + name)
			if err != nil {
				log.Println(err)
				log.Println("Error loading " + name + " file")
				//log.Fatalf("Error loading " + name + " file")
				//os.Exit(1)
			}
		}
	} else {
		path, _ := os.Getwd()
		config := strings.Split(path, "tests")[0] + ".env.test"
		err := godotenv.Load(config)
		if err != nil {
			log.Println(err)
			log.Fatalf("Error loading %s file", config)
		}
	}
}

// DotEnvVariable -> get .env
func DotEnvVariable(key, fallback string) string {
	var param string
	param = os.Getenv(key)
	if param == "" {
		param = fallback
	}
	return param
}
