package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/neoreads-backend/go/server/controllers"
	"github.com/neoreads-backend/go/server/repositories"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var db *sqlx.DB

func initDB(dbstring string) *sqlx.DB {
	log.Printf("dbstring: %s\n", dbstring)
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		log.Fatalf("init db failed: %s\n", err)
	}
	return db
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	v1 := r.Group("/api/v1")

	// /api/v1/book
	book := v1.Group("/book")
	{
		repo := repositories.NewBookRepo(db)
		ctrl := controllers.NewBookController(repo)

		book.GET("/:bookid", ctrl.GetBook)
		book.GET("/:bookid/:chapid", ctrl.GetBookChapter)
	}

	return r
}

// Config config object
type Config struct {
	port     string
	dbstring string
}

func initConfig() *Config {
	// read config
	viper.SetDefault("port", ":8080")
	viper.SetDefault("dbstring", "user=postgres dbname=neoreads sslmode=disable password=123456")

	viper.SetConfigType("toml")
	viper.SetConfigName("neoreads-server")
	viper.AddConfigPath("/etc/neoreads/")
	viper.AddConfigPath("$HOME/.neoreads/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error: config file %s not found", err))
	}

	config := &Config{}
	config.port = viper.GetString("port")
	config.dbstring = viper.GetString("dbstring")

	// Pretty print loaded configs
	js, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	if err != nil {
		panic(fmt.Errorf("Bad Config:%s", viper.AllSettings()))
	}
	log.Printf("Loaded Configs:\n%s", string(js))

	return config
}

func main() {

	// init config
	config := initConfig()

	// init database
	db = initDB(config.dbstring)

	// setup router
	r := setupRouter()

	// listen and serve
	r.Run(config.port)
}
