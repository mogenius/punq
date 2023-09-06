package operator

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/docs"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/version"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var HtmlDirFs embed.FS

func InitGin() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.StaticFS("/punq", embedFs())

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "punq API documentation"
	docs.SwaggerInfo.Description = "This is the documentation of all available API calls for the punq UI."
	docs.SwaggerInfo.Version = version.Ver
	//docs.SwaggerInfo.Host = fmt.Sprintf("0.0.0.0:%d", utils.CONFIG.Browser.Port)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", utils.CONFIG.Browser.Host, utils.CONFIG.Browser.Port)),
		ginSwagger.DefaultModelsExpandDepth(5),
		ginSwagger.DocExpansion("none")))

	InitContextRoutes(router)
	InitAuthRoutes(router)
	InitUserRoutes(router)
	InitGeneralRoutes(router)
	InitWorkloadRoutes(router)

	err := router.Run(fmt.Sprintf(":%d", OPERATOR_PORT))
	logger.Log.Errorf("Gin stopped with error: %s", err.Error())
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
