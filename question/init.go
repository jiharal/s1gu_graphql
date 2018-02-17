package question

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/neelance/graphql-go"

	log "github.com/sirupsen/logrus"
)

var (
	Logger    *log.Logger
	DbPool    *sql.DB
	CachePool *redis.Pool
	Schema    *graphql.Schema
)
