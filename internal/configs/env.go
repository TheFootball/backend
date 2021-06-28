package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DYNAMODB   string
	REDIS_ADDR string
	REDIS_PW   string
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
		DYNAMODB:   os.Getenv("DYNAMODB"),
		REDIS_ADDR: os.Getenv("REDIS_ADDRESS"),
		REDIS_PW:   os.Getenv("REDIS_PASSWORD"),
	}

	return env
}
