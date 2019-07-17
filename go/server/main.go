package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"

	"github.com/spf13/viper"

	"github.com/neoreads-backend/go/server/controllers"
	"github.com/neoreads-backend/go/server/models"
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
	JWTKey   string `json:"jwtkey"`
}

func initConfig() *Config {
	// read config
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
	viper.SetDefault("port", ":8080")
	config.Port = viper.GetString("port")

	viper.SetDefault("dbstring", "user=postgres dbname=neoreads sslmode=disable password=123456")
	config.DBString = viper.GetString("dbstring")

	viper.SetDefault("datadir", "D:/neoreads/data/")
	config.DataDir = viper.GetString("datadir")

	viper.SetDefault("prod", "false")
	config.Prod = viper.GetBool("prod")

	viper.SetDefault("jwtkey", "jwtkey secret very secret")
	config.JWTKey = viper.GetString("jwtkey")

	// Pretty print loaded configs
	js, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Bad Config:%s", viper.AllSettings())
	}
	if !config.Prod {
		log.Printf("Loaded Configs:\n%s", string(js))
	}

	return config
}

var identityKey = "id"

func initAuth(config *Config, userRepo *repositories.UserRepo) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.JWTKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.Credential
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if userRepo.CheckLogin(userID, password) {
				if user, found := userRepo.GetUser(userID); found {
					return &user, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatalf("JWT error: %s\n", err)
	}
	return authMiddleware
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

	userRepo := repositories.NewUserRepo(db)
	authMiddleware := initAuth(config, userRepo)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

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
	note.Use(authMiddleware.MiddlewareFunc())
	{
		repo := repositories.NewNoteRepo(db)
		ctrl := controllers.NewNoteController(repo)

		note.POST("/add", ctrl.AddNote)
		note.GET("/remove/:noteid", ctrl.RemoveNote)
		note.GET("/list", ctrl.ListNotes)
	}

	user := v1.Group("/user")
	{
		repo := repositories.NewUserRepo(db)
		ctrl := controllers.NewUserController(repo)
		user.POST("/register", ctrl.RegisterUser)

	}

	token := v1.Group("/token")
	{
		token.POST("/login", authMiddleware.LoginHandler)
		token.GET("/refresh", authMiddleware.RefreshHandler)
	}

	secure := v1.Group("/secure")
	secure.Use(authMiddleware.MiddlewareFunc())
	{
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
