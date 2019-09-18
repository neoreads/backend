package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"

	"github.com/spf13/viper"

	"github.com/neoreads/backend/server/controllers"
	"github.com/neoreads/backend/server/models"
	"github.com/neoreads/backend/server/repositories"
	"github.com/neoreads/backend/util"

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

var identityKey = "jwtuser"

func initAuth(config *Config, userRepo *repositories.UserRepo) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.JWTKey),
		Timeout:     1 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":       v.Id,
					"pid":      v.Pid,
					"username": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Id:       claims["id"].(string),
				Pid:      claims["pid"].(string),
				Username: claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.Credential
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			if userRepo.CheckLogin(username, password) {
				//return &loginVals, nil
				if user, found := userRepo.GetUser(username); found {
					return &user, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// TODO: implement authorizator
			if _, ok := data.(*models.User); ok {
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
	r.Static("/res/img", "./upload/img")
	userRepo := repositories.NewUserRepo(db)
	authMiddleware := initAuth(config, userRepo)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	v1 := r.Group("/api/v1")

	upload := v1.Group("/upload")
	upload.Use(authMiddleware.MiddlewareFunc())
	{
		imgIDGen := util.NewN64Generator(8)
		upload.POST("/img", func(ctx *gin.Context) {
			file, _ := ctx.FormFile("file")
			// TODO: check file name not to contain any path seperator to prevent security flaw
			if strings.Contains(file.Filename, "\\") || strings.Contains(file.Filename, "/") {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "cause": "file path insecure"})
				return
			}
			imgid := imgIDGen.Next()
			log.Printf("receiving img upload: %v\n", file.Filename)
			imgName := imgid + "_" + file.Filename
			fpath := path.Join("./upload/img/", imgName)
			err := ctx.SaveUploadedFile(file, fpath)
			if err != nil {
				log.Println(err)
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "ok", "imgid": imgid})
		})
	}

	// /api/v1/book
	book := v1.Group("/book")
	books := v1.Group("/books")
	{
		repo := repositories.NewBookRepo(db, config.DataDir)
		ctrl := controllers.NewBookController(repo)

		book.GET("/:bookid", ctrl.GetBook)
		book.GET("/:bookid/toc", ctrl.GetTOC)
		book.GET("/:bookid/chapter/:chapid", ctrl.GetBookChapter)

		books.GET("/hotlist", ctrl.HotList)

		books.Use(authMiddleware.MiddlewareFunc())
		books.POST("/add", ctrl.AddBook)
		books.POST("/modify", ctrl.ModifyBook)
		books.GET("/mine", ctrl.ListMyBooks)
		books.GET("/get/:bookid", ctrl.GetBook)
		books.GET("/remove/:bookid", ctrl.RemoveBook)
		books.GET("/chapter/get/:bookid/:chapid", ctrl.GetBookChapter)
		books.POST("/chapter/modify", ctrl.ModifyChapter)
	}

	note := v1.Group("/note")
	note.Use(authMiddleware.MiddlewareFunc())
	{
		repo := repositories.NewNoteRepo(db)
		ctrl := controllers.NewNoteController(repo)

		note.POST("/add", ctrl.AddNote)
		note.POST("/modify", ctrl.ModifyNote)
		note.GET("/remove/:noteid", ctrl.RemoveNote)
		note.GET("/list", ctrl.ListNotes)
	}

	reviews := v1.Group("/reviews")
	reviews.Use(authMiddleware.MiddlewareFunc())
	{
		repo := repositories.NewReviewRepo(db)
		ctrl := controllers.NewReviewController(repo)
		reviews.GET("/notes/:bookid/:chapid", ctrl.ListReviewNotes)
	}

	article := v1.Group("articles")
	article.Use(authMiddleware.MiddlewareFunc())
	{
		repo := repositories.NewArticleRepo(db)
		ctrl := controllers.NewArticleController(repo)
		article.GET("/list", ctrl.ListArticles)
		article.GET("/get/:artid", ctrl.GetArticle)
		article.GET("/remove/:artid", ctrl.RemoveArticle)
		article.POST("/add", ctrl.AddArticle)
		article.POST("/modify", ctrl.ModifyArticle)

		article.GET("/collection/:colid", ctrl.ListArticlesInCollection)
	}

	collections := v1.Group("collections")
	collections.Use(authMiddleware.MiddlewareFunc())
	{
		repo := repositories.NewCollectionRepo(db)
		ctrl := controllers.NewCollectionController(repo)

		collections.GET("/get/:colid", ctrl.GetCollection)
		collections.GET("/list", ctrl.ListCollections)
		collections.GET("/remove/:colid", ctrl.RemoveCollection)
		collections.POST("/add", ctrl.AddCollection)
		collections.POST("/modify", ctrl.ModifyCollection)
	}

	news := v1.Group("news")
	{
		repo := repositories.NewNewsRepo(db)
		ctrl := controllers.NewNewsController(repo)

		news.GET("/list", ctrl.ListNews)

		news.Use(authMiddleware.MiddlewareFunc())
		news.POST("/add", ctrl.AddNews)
	}

	tags := v1.Group("tags")
	{
		repo := repositories.NewTagRepo(db)
		ctrl := controllers.NewTagController(repo)

		tags.GET("/news/list", ctrl.ListNewsTags)
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
