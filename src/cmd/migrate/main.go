package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"strings"
)

// Env -> switch config file MODE variable
var Env string

// Version -> print version with start program
var Version = "development"

func init() {
	mode := os.Getenv("MODE")
	if mode == "production" || mode == "development" {
		Env = mode
	} else {
		Env = "production"
	}
}

func main() {
	color.Blue("Version:\t" + Version)
	if len(os.Args) == 1 {
		log.Fatalln("Missing options: up or down")
	}
	option := os.Args[1]
	clientOptions := options.Client().ApplyURI(DotEnvVariable("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("close error:", err)
		}
	}()

	driver, err := mongodb.WithInstance(client, &mongodb.Config{
		MigrationsCollection: "migrations",
		DatabaseName:         "notes",
	})
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"notes",
		driver)
	if err != nil {
		log.Fatalln(err)
	}
	switch option {
	case "new":
		if len(os.Args) != 3 {
			log.Fatalln("should be: new description-of-migration")
		}
		files, err := os.ReadDir("./migrations")
		if err != nil {
			log.Fatal(err)
		}
		tNum := 0
		for _, f := range files {
			split := strings.Split(f.Name(), "_")
			num, _ := strconv.Atoi(split[0])
			if num > tNum {
				tNum = num
			}
		}
		tNum = tNum + 1
		mNum := fmt.Sprintf("%03d", tNum)
		fNameUp := fmt.Sprintf("./migrations/%s_%s.up.json", mNum, os.Args[2])
		fNameDown := fmt.Sprintf("./migrations/%s_%s.down.json", mNum, os.Args[2])
		f, err := os.Create(fNameUp)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(f)
		f, err = os.Create(fNameDown)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(f)
		color.Green("Migrate template created")
	case "up":
		if len(os.Args) == 3 {
			mNum, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalln("Incorrect arg for number migration")
			}
			err = m.Steps(mNum)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			err = m.Up()
			if err != nil {
				log.Fatalln(err)
			}
		}
		color.Green("Migrate up complete")
	case "down":
		if len(os.Args) == 3 {
			mNum, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalln("Incorrect arg for number migration")
			}
			mNum = -1 * mNum
			err = m.Steps(mNum)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			err = m.Down()
			if err != nil {
				log.Fatalln(err)
			}
		}
		color.Green("Migrate down complete")
	}
}

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {
	if flag.Lookup("test.v") == nil {
		name := ".env"
		if Env == "development" {
			name = ".env.test"
		}
		err := godotenv.Load("../" + name)
		if err != nil {
			err = godotenv.Load(name)
			if err != nil {
				log.Fatalf("Error loading " + name + " file")
				os.Exit(1)
			}
		}
	} else {
		path, _ := os.Getwd()
		err := godotenv.Load(strings.Split(path, "src")[0] + "/" + ".env.test")
		fmt.Println(path)
		if err != nil {
			log.Fatalf("Error loading .env.test file")
		}
	}
	return os.Getenv(key)
}
