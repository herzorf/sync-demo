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
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {

	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			_, err := c.Writer.Write([]byte("hello"))
			if err != nil {
				log.Fatal(err)
			}
		})
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		err := router.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()

	cmd := exec.Command(
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"--app=http:127.0.0.1:8080/",
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
