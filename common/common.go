package common

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Common struct {
	Log *zap.Logger
	Env Env
}

type Env struct {
	API_PORT int

	PGSQL_DATA     string
	PGSQL_HOST     string
	PGSQL_PORT     string
	PGSQL_USER     string
	PGSQL_PASSWORD string
	PGSQL_DBNAME   string
	PGSQL_NAME     string

	REDIS_HOST     string
	REDIS_PASSWORD string
	REDIS_PORT     string
	REDIS_NAME     string
	REDIS_DB       int

	AUTH_URL      string
	AUTH_PORT     string
	AUTH_ENDPOINT string

	JWT_SIGN []byte
}

func NewCommon() *Common {
	l, _ := zap.NewProduction()

	if err := godotenv.Load(); err != nil {
		l.Error(err.Error())
		return &Common{}
	}

	apiPortStr := os.Getenv("API_PORT")
	apiPort, err := strconv.Atoi(apiPortStr)

	if err != nil {
		l.Error(err.Error())
		return &Common{}
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		l.Error(err.Error())
		return &Common{}
	}
	return &Common{
		Log: l,
		Env: Env{
			API_PORT:       apiPort,
			PGSQL_DATA:     os.Getenv("PGSQL_DATA"),
			PGSQL_HOST:     os.Getenv("PGSQL_HOST"),
			PGSQL_PORT:     os.Getenv("PGSQL_PORT"),
			PGSQL_USER:     os.Getenv("PGSQL_USER"),
			PGSQL_PASSWORD: os.Getenv("PGSQL_PASSWORD"),
			PGSQL_DBNAME:   os.Getenv("PGSQL_DBNAME"),
			PGSQL_NAME:     os.Getenv("PGSQL_NAME"),

			REDIS_HOST:     os.Getenv("REDIS_HOST"),
			REDIS_PORT:     os.Getenv("REDIS_PORT"),
			REDIS_NAME:     os.Getenv("REDIS_NAME"),
			REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
			REDIS_DB:       redisDB,

			AUTH_URL:      os.Getenv("AUTH_URL"),
			AUTH_PORT:     os.Getenv("AUTH_PORT"),
			AUTH_ENDPOINT: os.Getenv("AUTH_ENDPOINT"),

			JWT_SIGN: []byte(os.Getenv("JWT_SIGN")),
		},
	}
}
