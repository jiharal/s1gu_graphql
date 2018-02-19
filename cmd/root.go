package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/s1gu/s1gu-lib/cache"
	"github.com/s1gu/s1gu-lib/db"

	"github.com/s1gu/s1gu_graphql/question"
	"github.com/s1gu/s1gu_graphql/starwars"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger    *log.Logger
	dbPool    *sql.DB
	cachePool *redis.Pool

	cfgFile string
)

var RootCmd = &cobra.Command{
	Use:   "AuthS Api",
	Short: "AuthS Api",
	Long:  `AuthS website API, Provide data for AuthS FrontEnd`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
			AUTH SECURE SERVER
			`)
		fmt.Println("AuthS running")
		fmt.Println("Version:", viper.GetString("app.version"))
		fmt.Println("App port:", viper.GetString("app.port"))
		fmt.Println("Host :", viper.GetString("database.host"))
		fmt.Println("DBName :", viper.GetString("database.name"))
		go initDB()
		go initCache()
	},
	Run: func(cmd *cobra.Command, args []string) {

		// handler graphql server for client
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(page)
		}))

		http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next := &relay.Handler{Schema: question.Schema}
			authorization := r.Header.Get("Authorization")
			token := strings.Replace(authorization, "Bearer ", "", 1)
			ctx := context.WithValue(r.Context(), "AuthorizationToken", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}))

		http.Handle("/star", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next := &relay.Handler{Schema: starwars.Schema}
			authorization := r.Header.Get("Authorization")
			token := strings.Replace(authorization, "Bearer ", "", 1)
			ctx := context.WithValue(r.Context(), "AuthorizationToken", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}))

		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), nil)
	},
}

// Add graphql client
var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger, initGraphQLserver)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.s1gu.config.toml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/Users/zzz/go/src/github.com/s1gu/exp-modem2phone")
		viper.SetConfigName(".s1gu.config")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initDB() {
	dbOptions := db.DBOptions{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.name"),
		SSLMode:  viper.GetString("database.sslmode"),
	}

	fmt.Println("DBptions:", dbOptions)

	dbConn, err := db.Connect(dbOptions)
	if err != nil {
		fmt.Println("Gagal konek", err)
		panic(err)
	}
	dbPool = dbConn
	question.DbPool = dbPool
}

func initCache() {
	cacheOptions := cache.CacheOptions{
		Host:        viper.GetString("cache.host"),
		Port:        viper.GetInt("cache.port"),
		MaxIdle:     viper.GetInt("cache.max_idle"),
		IdleTimeout: viper.GetInt("cache.idle_timeout"),
		Enabled:     viper.GetBool("cache.enabled"),
	}
	pool := cache.Connect(cacheOptions)
	cachePool = pool
	question.CachePool = cachePool
}

func initLogger() {
	logger = log.New()
	logger.Formatter = &log.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	question.Logger = logger
}

func initGraphQLserver() {
	// Schema for question
	if QuestionSchema, err := ioutil.ReadFile("question/schema.graphql"); err != nil {
		panic(err)
	} else {
		question.Schema = graphql.MustParseSchema(string(QuestionSchema), &question.Resolver{})
	}

	// Scema for starwars
	if StarwarsSchema, err := ioutil.ReadFile("starwars/schema.graphql"); err != nil {
		panic(err)
	} else {
		starwars.Schema = graphql.MustParseSchema(string(StarwarsSchema), &starwars.Resolver{})
	}

}
