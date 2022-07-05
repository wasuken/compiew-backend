package zip

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 圧縮ファイル内部コンテンツの抽象化構造体
type ZFileInfo struct {
	Content string
	// 相対パス
	Path  string
	IsDir bool
}

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

	tmpdir := os.Getenv("TEMP_DIR") + "/"
	cur, _ := os.Getwd()
	tmp, _ := os.Create(cur + tmpdir + hsUrl)
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

// 圧縮ファイルをtmpへ吐き出す処理
// zname 圧縮ファイルの名前、この名前でディレクトリを作成する
func dumpZipFile(zname string, savePathes []*ZFileInfo) error {
	tmpdir := os.Getenv("TEMP_DIR") + "/"
	cur, _ := os.Getwd()
	zipDir := cur + tmpdir + zname
	err := os.MkdirAll(zipDir, 0777)
	// 既にする場合もあるので、とりあえずログ吐くだけ
	if err != nil {
		fmt.Println(err)
	}

	var f *os.File
	for _, zfile := range savePathes {
		if zfile.IsDir {
			os.MkdirAll(zipDir+"/"+zfile.Path, 0777)
		} else {
			f, err = os.Create(zipDir + "/" + zfile.Path)
			defer f.Close()
			_, err = f.WriteString(zfile.Content)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		return err
	}
	return nil

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
	savePathes := []*ZFileInfo{}
	pathes := []string{}
	for {
		tHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		pathes = append(pathes, tHeader.Name)
		var content string

		if !tHeader.FileInfo().IsDir() {

			buf := new(bytes.Buffer)
			buf.ReadFrom(tarReader)

			content = buf.String()
		}
		p := &ZFileInfo{Path: tHeader.Name, IsDir: tHeader.FileInfo().IsDir(), Content: content}
		savePathes = append(savePathes, p)
	}
	err := dumpZipFile(filepath.Base(path), savePathes)
	if err != nil {
		return []string{}, err
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

	savePathes := []*ZFileInfo{}
	pathes := []string{}
	for _, file := range read.File {
		var content string
		if !file.Mode().IsDir() {
			fs, errr := file.Open()
			if errr != nil {
				continue
			}
			defer fs.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(fs)

			content = buf.String()

		}
		p := &ZFileInfo{Path: file.FileHeader.Name, IsDir: file.Mode().IsDir(), Content: content}
		savePathes = append(savePathes, p)
		pathes = append(pathes, file.FileHeader.Name)
	}
	err = dumpZipFile(filepath.Base(zfile.Name()), savePathes)
	if err != nil {
		return []string{}, err
	}
	return pathes, nil
}

func GetZipFileContent(path string) (string, error) {
	return path, nil
}
