package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DYNAMODB string
}

var env *Env

func GetEnv() *Env {
	if env != nil {
		return env
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	env = &Env{
		DYNAMODB: os.Getenv("DYNAMODB"),
	}

	return env
}
