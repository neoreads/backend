package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/neoreads-backend/go/server/controllers"
	"github.com/neoreads-backend/go/server/repositories"
	"github.com/neoreads-backend/go/util"

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

// Config config object
type Config struct {
	Port     string `json:"port"`
	DBString string `json:"dbstring"`
	DataDir  string `json:"datadir"`
	Prod     bool   `json:"prod"`
}

func initConfig() *Config {
	// read config
	viper.SetDefault("port", ":8080")
	viper.SetDefault("dbstring", "user=postgres dbname=neoreads sslmode=disable password=123456")
	viper.SetDefault("datadir", "D:/neoreads/data/")
	viper.SetDefault("prod", "true")

	viper.SetConfigType("toml")
	viper.SetConfigName("neoreads-server")
	viper.AddConfigPath("/etc/neoreads/")
	viper.AddConfigPath("$HOME/.neoreads/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error: config file %s not found", err)
	}

	config := &Config{}
	config.Port = viper.GetString("port")
	config.DBString = viper.GetString("dbstring")
	config.DataDir = viper.GetString("datadir")
	config.Prod = viper.GetBool("prod")

	// Pretty print loaded configs
	js, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Bad Config:%s", viper.AllSettings())
	}
	log.Printf("Loaded Configs:\n%s", string(js))

	return config
}

func initRouter(config *Config) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/index.html", "./public/index.html")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.Static("/css", "./public/css")
	r.Static("/js", "./public/js")
	r.Static("/fonts", "./public/fonts")
	v1 := r.Group("/api/v1")

	// /api/v1/book
	book := v1.Group("/book")
	{
		repo := repositories.NewBookRepo(db, config.DataDir)
		ctrl := controllers.NewBookController(repo)

		book.GET("/:bookid", ctrl.GetBook)
		book.GET("/:bookid/toc", ctrl.GetTOC)
		book.GET("/:bookid/chapter/:chapid", ctrl.GetBookChapter)
	}

	note := v1.Group("/note")
	{
		repo := repositories.NewNoteRepo(db)
		ctrl := controllers.NewNoteController(repo)

		note.POST("/add", ctrl.AddNote)
		note.GET("/remove/:noteid", ctrl.RemoveNote)
		note.GET("/list", ctrl.ListNotes)
	}

	return r
}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	//log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {
	util.InitSeed()

	// init config
	config := initConfig()

	// init database
	db = initDB(config.DBString)

	// init router
	r := initRouter(config)

	// listen and serve
	if config.Prod {
		gin.SetMode(gin.ReleaseMode)
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		log.Fatal(http.ListenAndServeTLS(":443", "certs/cert.pem", "certs/key.pem", r))
	} else {
		r.Run(config.Port)
	}
}
