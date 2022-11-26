package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

//go:embed frontend/dist/*
var FS embed.FS

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	} else {
		executable, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(executable)
		filename := uuid.New().String()
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fullPath := path.Join("uploads", filename+".txt")
		fmt.Println(fullPath)
		err = os.WriteFile(filepath.Join(dir, fullPath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullPath})
	}
}

func main() {

	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.POST("/api/v1/texts", TextsController)
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			p := c.Request.URL.Path
			if strings.HasPrefix(p, "/static") {
				file, err := staticFiles.Open("index.html")
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
	signal.Notify(chSingal, syscall.SIGTERM)

	select {
	case <-chSingal:
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal(err)
		}
	}
}
