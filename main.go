package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/lichenxin/m3u8/dl"
)

var dataPath = "/tmp/mp4/"

var port string

// ffmpeg -i 2.ts -acodec copy -vcodec copy -absf aac_adtstoasc output.mp4
// https://43.154.3.196/vtt/movie1/m/05/%E5%9B%9E%E9%AD%82.m3u8
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/m3u8-http-server
func init() {
	flag.StringVar(&port, "port", "5050", "http server port")
	flag.Parse()

	dataPath = fmt.Sprintf("/%s/", strings.Trim(dataPath, "/"))
	fmt.Println(dataPath)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if len(url) > 0 && path.Ext(url) == ".m3u8" {
			go downloadTsFileToMP4(url)
			w.Write([]byte("ok"))
		} else {
			files, err := readFiles(dataPath)
			if err != nil {
				fmt.Println("[error]", err.Error())
			}

			w.Write([]byte(strings.Join(files, "\n")))
		}
	})

	fs := http.FileServer(http.Dir(dataPath))
	http.Handle("/data/", http.StripPrefix("/data/", fs))

	fmt.Println("[start http]", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}

func downloadTsFileToMP4(url string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[error]", fmt.Errorf("%v", err))
		}
	}()

	name := MD5([]byte(url))
	mp4File := dataPath + fmt.Sprintf("%s.mp4", name)
	if fileExist(mp4File) {
		fmt.Println("[file]", mp4File, "exists")
		return
	}

	tsFile := fmt.Sprintf("/tmp/m3u8")
	downloader, err := dl.NewTask(tsFile, url)
	if err != nil {
		panic(err)
	}

	if err := downloader.Start(20); err != nil {
		panic(err)
	}

	defer func() {
		if er := fileRemove(tsFile); er != nil {
			fmt.Println("[error]", err)
		}
	}()

	tsFile = tsFile + "/main.ts"
	fmt.Println("[file]", tsFile, mp4File)
	tsFileMergeMP4(tsFile, mp4File)
}

func tsFileMergeMP4(input, output string) {
	cmd := exec.Command("ffmpeg", "-i", input, "-acodec", "copy", "-vcodec", "copy", "-absf", "aac_adtstoasc", output)
	fmt.Println("[ffmpeg]", cmd.String())

	if err := cmd.Start(); err != nil {
		fmt.Println("[cmd start]", err.Error())
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("[cmd wait]", err.Error())
	}

	fmt.Println("[done]", output)
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func fileRemove(name string) error {
	if len(name) > 0 && fileExist(name) {
		return os.Remove(name)
	}
	return nil
}

func MD5(v []byte) string {
	return fmt.Sprintf("%x", md5.Sum(v))
}

func readFiles(path string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, fmt.Sprintf(`<p><a href="/data/%s">%s</a></p>`, info.Name(), info.Name()))
		}
		return nil
	})
	return files, err
}
