package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/lichenxin/m3u8/dl"
    "github.com/lichenxin/m3u8/tool"
)

var inputURL, outputName, dataPath string

// ffmpeg -i main.ts -acodec copy -vcodec copy -absf aac_adtstoasc output.mp4
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/m3u8-downloader
func init() {
    flag.StringVar(&inputURL, "u", "", "链接地址输入")
    flag.StringVar(&dataPath, "p", "./data", "输入目录，默认当前目录中")
    flag.StringVar(&outputName, "name", "", "输出名称，默认使用链接地址的MD5值作为名称")
    flag.Parse()

    if !strings.HasPrefix(inputURL, "http") {
        panic("输入正确的链接地址")
    }

    dir, err := filepath.Abs(dataPath)
    if err != nil {
        panic(err)
    }

    dataPath = fmt.Sprintf("/%s/", strings.Trim(dir, "/"))
    if err = os.MkdirAll(dataPath, os.ModePerm); err != nil {
        panic(err)
    }
}

func main() {
    if strings.HasSuffix(inputURL, ".m3u8") {
        downloadTsFileToMP4(inputURL)
    } else {
        b, err := tool.Get("http://211.149.238.132:8885?url=" + inputURL)
        if err != nil {
            panic(err)
        }
        res, err := io.ReadAll(b)
        if err != nil {
            panic(err)
        }
        fmt.Println(string(res))
    }
}

func downloadTsFileToMP4(url string) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println("[error]", fmt.Errorf("%v", err))
        }
    }()

    name := tool.MD5([]byte(url))
    if outputName != "" {
        name = outputName
    }
    mp4File := dataPath + fmt.Sprintf("%s.mp4", name)
    if tool.FileExist(mp4File) {
        fmt.Println("[file]", mp4File, "exists")
        return
    }

    tsFile := fmt.Sprintf("/tmp/m3u8")
    downloader, err := dl.NewTask(tsFile, url)
    if err != nil {
        panic(err)
    }

    if err = downloader.Start(20); err != nil {
        panic(err)
    }

    defer func() {
        if er := tool.FileRemove(tsFile); er != nil {
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
