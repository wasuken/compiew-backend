package zip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
func GetZipFileInfo(url string) ([]string, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, url)
	hsUrl := hex.EncodeToString(sha1.Sum(nil))

	tmp, _ := ioutil.TempFile("", hsUrl)
	defer tmp.Close()

	err := download(tmp, url)
	if err != nil {
		return []string{}, err
	}

	rst := strings.Split(url, "?")
	ps := strings.Split(rst[0], ".")
	ext := ps[len(ps)-1]

	if strings.ToLower(ext) == "gz" {
		tmp.Close()
		return parseTarGz(tmp.Name())
	} else {
		return parsePKZip(tmp)
	}
}

// tar.gzをパース、ファイルパスを取得
// tmpfileの削除もこちらで行う
func parseTarGz(path string) ([]string, error) {
	zfile, errr := os.Open(path)
	if errr != nil {
		os.Remove(path)
		return []string{}, errr
	}
	defer zfile.Close()
	read, er := gzip.NewReader(zfile)
	defer read.Close()
	if er != nil {
		os.Remove(path)
		return []string{}, er
	}

	tarReader := tar.NewReader(read)
	pathes := []string{}
	for {
		tHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		pathes = append(pathes, tHeader.Name)
	}
	// tempfileの削除
	er = os.Remove(path)
	if er != nil {
		// 消せなくても一応返却はする
		return pathes, er
	}

	return pathes, nil
}

// tmpfileの削除もこちらで行う
// zipをパース、ファイルパスを取得
func parsePKZip(zfile *os.File) ([]string, error) {
	read, err := zip.OpenReader(zfile.Name())
	if err != nil {
		name := zfile.Name()
		zfile.Close()
		os.Remove(name)
		return []string{}, err
	}
	defer read.Close()

	pathes := []string{}
	for _, file := range read.File {
		pathes = append(pathes, file.FileHeader.Name)
	}
	err = os.Remove(zfile.Name())
	if err != nil {
		// 消せなくても一応返却はする
		name := zfile.Name()
		zfile.Close()
		os.Remove(name)
		return pathes, err
	}
	return pathes, nil
}
