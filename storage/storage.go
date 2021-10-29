package storage

import (
	"fmt"

	"github.com/fiuskylab/grpc-example/common"
	redis "github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	PGSQL  *gorm.DB
	Redis  *redis.Client
	Common *common.Common
}

func NewStorage(c *common.Common) (*Storage, error) {
	s := Storage{
		Common: c,
	}

	var pgsqlURL string
	url := `host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`

	pgsqlURL = fmt.Sprintf(
		url,
		c.Env.PGSQL_HOST,
		c.Env.PGSQL_PORT,
		c.Env.PGSQL_USER,
		c.Env.PGSQL_PASSWORD,
		c.Env.PGSQL_DBNAME,
	)

	db, err := gorm.Open(postgres.Open(pgsqlURL), &gorm.Config{})

	if err != nil {
		c.Log.Error(err.Error())
		return &s, err
	}
	s.PGSQL = db

	s.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Env.REDIS_HOST + ":" + c.Env.REDIS_PORT,
		Password: c.Env.REDIS_PASSWORD,
		DB:       c.Env.REDIS_DB,
	})

	return &s, nil
}
