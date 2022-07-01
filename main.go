package main

import (
	myzip "compiew_api/zip"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Resp struct {
	Status int      `json:"status"`
	Pathes []string `json:"pathes"`
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3333"},
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
		r := Resp{Status: 200, Pathes: pathes}

		return c.JSON(http.StatusOK, r)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
