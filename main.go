package main

import (
	myzip "compiew_api/zip"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type ZipInfoResp struct {
	Status int      `json:"status"`
	Pathes []string `json:"pathes"`
}

type ZipContentResp struct {
	Status  int    `json:"status"`
	Content string `json:"content"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	origins := os.Getenv("ORIGINS")
	oris := strings.Split(origins, ",")
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: oris,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.GET("/api", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/api/zipinfo", func(c echo.Context) error {
		url := c.QueryParam("url")
		pathes, err := myzip.GetZipFileInfo(url)
		if err != nil {
			e := errors.Wrap(err, "error")
			log.Printf("%+v", e)
		}
		r := ZipInfoResp{Status: 200, Pathes: pathes}

		return c.JSON(http.StatusOK, r)
	})
	e.GET("/api/zip/content", func(c echo.Context) error {
		path := c.QueryParam("path")
		path = strings.ReplaceAll(path, "'", "")
		url := c.QueryParam("url")
		url = strings.ReplaceAll(url, "'", "")
		status := 200
		content, err := myzip.GetZipFileContent(url, path)
		if err != nil {
			e := errors.Wrap(err, "error")
			log.Printf("%+v", e)
			status = 500
			content = "error"
		}
		r := ZipContentResp{Status: status, Content: content}

		return c.JSON(http.StatusOK, r)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
