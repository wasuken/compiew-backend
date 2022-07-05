package main

import (
	myzip "compiew_api/zip"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3333", "http://wasu-arch:3333"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/zipinfo", func(c echo.Context) error {
		url := c.QueryParam("url")
		pathes, err := myzip.GetZipFileInfo(url)
		if err != nil {
			panic(err)
		}
		r := ZipInfoResp{Status: 200, Pathes: pathes}

		return c.JSON(http.StatusOK, r)
	})
	e.GET("/zip/content", func(c echo.Context) error {
		path := c.QueryParam("path")
		content, err := myzip.GetZipFileContent(path)
		if err != nil {
			panic(err)
		}
		r := ZipContentResp{Status: 200, Content: content}

		return c.JSON(http.StatusOK, r)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
