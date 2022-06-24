package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func download(dlfile *os.File, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, er := io.Copy(dlfile, resp.Body)
	return er
}

// URLパラメタ
// URLでzip,tar.gzを取得
func getZipFileInfo(url string) ([]string, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, url)
	hsUrl := hex.EncodeToString(sha1.Sum(nil))
	tmp, _ := ioutil.TempFile("", hsUrl)
	defer tmp.Close()

	err := download(tmp, url)
	if err != nil {
		panic(err)
	}

	rst := strings.Split(url, "?")
	ps := strings.Split(rst[0], ".")
	ext := ps[len(ps)-1]

	if strings.ToLower(ext) == "gz" {
		return parseTarGz(tmp)
	} else {
		return parsePKZip(tmp)
	}
}

// tar.gzをパース、ファイルパスを取得
func parseTarGz(zfile *os.File) ([]string, error) {
	read, _ := gzip.NewReader(zfile)
	defer read.Close()
	tarReader := tar.NewReader(read)
	pathes := []string{}
	for {
		tHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		pathes = append(pathes, tHeader.Name)
	}

	return pathes, nil
}

// zipをパース、ファイルパスを取得
func parsePKZip(zfile *os.File) ([]string, error) {
	read, err := zip.OpenReader(zfile.Name())
	if err != nil {
		return []string{}, err
	}
	defer read.Close()

	pathes := []string{}
	for _, file := range read.File {
		pathes = append(pathes, file.FileHeader.Name)
	}
	return pathes, nil
}

type Resp struct {
	Status int      `json:"status"`
	Pathes []string `json:"pathes"`
}

func main() {
	fmt.Println("hello.")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/zipinfo", func(c echo.Context) error {
		url := c.QueryParam("url")
		pathes, err := getZipFileInfo(url)
		if err != nil {
			panic(err)
		}
		r := Resp{Status: 200, Pathes: pathes}

		return c.JSON(http.StatusOK, r)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
