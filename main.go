package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {

	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static") {
				file, err := staticFiles.Open("index/html")
				if err != nil {
					log.Fatal(err)
				}
				defer func(file fs.File) {
					err := file.Close()
					if err != nil {

					}
				}(file)
				stat, err := file.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html", file, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		err := router.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()

	cmd := exec.Command(
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"--app=http:127.0.0.1:8080/static/index.html",
		"--user-data-dir=test-user-data",
	)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	chSingal := make(chan os.Signal, 1)
	signal.Notify(chSingal, os.Interrupt)

	select {
	case <-chSingal:
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal(err)
		}
	}
}
