package authmc

import (
	"database/sql"
	"html/template"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

var (
	Logger        *log.Logger
	DbPool        *sql.DB
	CachePool     *redis.Pool
	LoginTemplate *template.Template
	Gsm           *AUTHMC
)
