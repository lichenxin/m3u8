# M3U8 视频下载器

一个基于 Go 语言开发的 M3U8 视频下载和合并工具，支持加密视频的解密和合并为 MP4 格式。

## 功能特性

- 支持 M3U8 视频流下载
- 支持 AES-128 加密视频解密
- 多线程并发下载，提高下载速度
- 自动合并 TS 片段为 MP4 文件
- 支持自定义输出目录和文件名

## 依赖环境

### FFmpeg

本项目依赖 [FFmpeg](https://ffmpeg.org/) 进行视频格式转换，需要先安装 FFmpeg。

#### 自动安装（推荐）

项目提供了自动安装脚本，支持 macOS 和 Linux 系统：

```bash
./install_ffmpeg.sh
```

该脚本会：
1. 检查系统是否已安装 FFmpeg
2. 如果未安装，则根据操作系统自动选择合适的安装方式
3. 对于 macOS，会自动安装 Homebrew（如果未安装）并使用 Homebrew 安装 FFmpeg
4. 对于 Linux，会根据发行版使用相应的包管理器安装 FFmpeg

#### 手动安装

您也可以手动安装 FFmpeg：

- **macOS**: `brew install ffmpeg`
- **Ubuntu/Debian**: `sudo apt update && sudo apt install ffmpeg`
- **CentOS/RHEL**: `sudo yum install epel-release && sudo yum install ffmpeg`
- **Fedora**: `sudo dnf install ffmpeg`

安装完成后，可以通过以下命令验证：

```bash
ffmpeg -version
```

## 安装

```bash
go build -o m3u8 cmd/main.go
```

## 使用方法

### 基本用法

```bash
./m3u8 -u [M3U8_URL] [-p 输出目录] [-name 输出文件名]
```

### 参数说明

- `-u`: M3U8 视频链接地址（必需）
- `-p`: 输出目录，默认为当前目录下的 `data` 文件夹
- `-name`: 输出文件名，默认使用链接地址的 MD5 值作为文件名

### 示例

```bash
# 解析地址，如腾讯，爱奇艺播放地址
./m3u8 -u "https://www.iqiyi.com/v_1xkc9zgg4to.html"

# 指定输出目录和文件名
./m3u8 -u "https://example.com/video/index.m3u8" -p "./downloads" -name "my_video"

# 使用默认目录，自定义文件名
./m3u8 -u "https://example.com/video/index.m3u8" -name "course_video"
```

### 输出说明

下载完成后，程序会自动将 TS 片段合并为 MP4 文件，输出路径为：
```
[输出目录]/[文件名].mp4
```

## 工作原理

1. 解析 M3U8 文件，获取所有 TS 片段信息
2. 多线程并发下载所有 TS 片段
3. 对加密的 TS 片段进行 AES-128 解密
4. 将所有 TS 片段按顺序合并为一个完整的 TS 文件
5. 使用 FFmpeg 将合并后的 TS 文件转换为 MP4 格式

## 注意事项

1. 确保网络连接稳定，下载过程中断可能导致文件不完整
2. 对于加密视频，程序会自动处理解密过程
3. 下载速度取决于网络带宽和服务器响应速度
4. 大文件下载可能需要较长时间，请耐心等待