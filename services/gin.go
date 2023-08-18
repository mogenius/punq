package services

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
)

var HtmlDirFs embed.FS

func InitGin() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.StaticFS("/punq", embedFs())

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", utils.CONFIG.Browser.Host, utils.CONFIG.Browser.Port),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		logger.Log.Info("listen: %s\n", err)
	}
}

func embedFs() http.FileSystem {
	sub, err := fs.Sub(HtmlDirFs, "ui/dist")

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

func printPrettyPost(c *gin.Context) {
	var out bytes.Buffer
	body, _ := io.ReadAll(c.Request.Body)
	err := json.Indent(&out, []byte(body), "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out.Bytes()))
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
