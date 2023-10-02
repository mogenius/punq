package operator

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var HtmlDirFs embed.FS

func InitFrontend() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}

	router.Use(cors.New(config))
	router.Use(CreateLogger("ANGULAR"))

	router.StaticFS("/", embedFs())

	router.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", getIndexHtml())
	})

	err := router.Run(fmt.Sprintf(":%d", utils.CONFIG.Frontend.Port))
	logger.Log.Errorf("Frontend (gin) stopped with error: %s", err.Error())
}

func InitBackend() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}

	router.Use(cors.New(config))
	router.Use(CreateLogger("BACKEND"))

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "punq API documentation"
	docs.SwaggerInfo.Description = "This is the documentation of all available API calls for the punq UI."
	docs.SwaggerInfo.Version = version.Ver
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DefaultModelsExpandDepth(3), ginSwagger.DocExpansion("none")))

	InitContextRoutes(router)
	InitAuthRoutes(router)
	InitUserRoutes(router)
	InitGeneralRoutes(router)
	InitWorkloadRoutes(router)

	err := router.Run(fmt.Sprintf(":%d", utils.CONFIG.Backend.Port))
	logger.Log.Errorf("Operator (gin) stopped with error: %s", err.Error())
}

func embedFs() http.FileSystem {
	sub, err := fs.Sub(HtmlDirFs, "ui/dist")
	if err != nil {
		logger.Log.Fatalf("Cannot load ui/dist from filesystem.")
	}

	dirContent, err := getAllFilenames(&HtmlDirFs, "")
	if err != nil {
		panic(err)
	}

	if len(dirContent) <= 0 {
		panic("dist folder empty. Cannnot serve site. FATAL.")
	} else {
		logger.Log.Noticef("Loaded %d static files from embed.", len(dirContent))
	}
	return http.FS(sub)
}

func getIndexHtml() []byte {
	data, err := HtmlDirFs.ReadFile("ui/dist/index.html")
	if err != nil {
		logger.Log.Fatal("Cannot load index.html from filesystem.")
		return []byte{}
	}
	return data
}

func getAllFilenames(fs *embed.FS, dir string) (out []string, err error) {
	if len(dir) == 0 {
		dir = "."
	}

	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fp := path.Join(dir, entry.Name())
		if entry.IsDir() {
			res, err := getAllFilenames(fs, fp)
			if err != nil {
				return nil, err
			}

			out = append(out, res...)

			continue
		}

		out = append(out, fp)
	}
	return
}
