#!/bin/bash

# 脚本名称: install_ffmpeg_checked.sh
# 描述: 自动检测操作系统，检查 FFmpeg 是否已安装，并安装或跳过。

# --- 颜色定义 ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}### FFmpeg 跨平台安装脚本 (带检查) ###${NC}"

# 1. 检查 FFmpeg 是否已安装
echo -e "${YELLOW}正在检查 FFmpeg 是否已安装...${NC}"
if command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${GREEN}FFmpeg 已经安装!${NC}"
    echo "版本信息:"
    ffmpeg -version | head -n 1
    echo "跳过安装步骤。"
    exit 0
fi
echo -e "${YELLOW}FFmpeg 未安装。继续执行安装...${NC}"

# 2. 检测操作系统并执行安装
OS=$(uname -s)

if [ "$OS" == "Linux" ]; then
    echo -e "${YELLOW}检测到系统: Linux${NC}"

    # 尝试检测 Linux 发行版使用的包管理器
    if command -v apt >/dev/null; then
        echo "使用 apt (Debian/Ubuntu/Mint) 安装 FFmpeg..."
        sudo apt update || { echo -e "${RED}apt update 失败${NC}"; exit 1; }
        sudo apt install ffmpeg -y
    elif command -v dnf >/dev/null; then
        echo "使用 dnf (Fedora/CentOS Stream) 安装 FFmpeg..."
        # 启用 RPM Fusion 仓库 (推荐方式)
        sudo dnf install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm -y
        sudo dnf update
        sudo dnf install ffmpeg -y
    elif command -v yum >/dev/null; then
        echo "使用 yum (CentOS/RHEL) 安装 FFmpeg..."
        # 启用 EPEL 和 RPM Fusion 仓库
        sudo yum install epel-release -y
        sudo yum localinstall --nogpgcheck https://download1.rpmfusion.org/free/el/rpmfusion-free-release-7.noarch.rpm -y
        sudo yum update
        sudo yum install ffmpeg -y
    elif command -v pacman >/dev/null; then
        echo "使用 pacman (Arch Linux) 安装 FFmpeg..."
        sudo pacman -Syu --noconfirm
        sudo pacman -S ffmpeg --noconfirm
    else
        echo -e "${RED}错误: 未知的 Linux 包管理器。请手动安装 FFmpeg。${NC}"
        exit 1
    fi

elif [ "$OS" == "Darwin" ]; then
    echo -e "${YELLOW}检测到系统: macOS${NC}"

    # 检查 Homebrew 是否安装
    if ! command -v brew >/dev/null; then
        echo -e "${YELLOW}Homebrew 未安装。开始安装 Homebrew...${NC}"
        # Homebrew 官方安装命令
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" || { echo -e "${RED}Homebrew 安装失败${NC}"; exit 1; }
        # 确保 Homebrew 在当前 shell 中可用 (防止在同一个脚本中后续命令找不到 brew)
        if [ -f /opt/homebrew/bin/brew ]; then
            export PATH="/opt/homebrew/bin:$PATH"
        elif [ -f /usr/local/bin/brew ]; then
            export PATH="/usr/local/bin:$PATH"
        fi
    fi

    echo "使用 Homebrew 安装 FFmpeg..."
    brew install ffmpeg

else
    echo -e "${RED}错误: 不支持的操作系统 (${OS})。本脚本仅支持 Linux 和 macOS。${NC}"
    echo -e "${YELLOW}对于 Windows，请使用 winget 或 Chocolatey (choco) 进行安装。${NC}"
    exit 1
fi

# 3. 最终验证安装结果
echo -e "\n${GREEN}--- 最终验证 ---${NC}"
if command -v ffmpeg >/dev/null 2>&1; then
    echo -e "${GREEN}FFmpeg 已成功安装!${NC}"
    ffmpeg -version | head -n 1
else
    echo -e "${RED}FFmpeg 安装命令执行完毕，但 'ffmpeg' 命令未找到。请检查安装过程中的错误提示。${NC}"
    exit 1
fi

exit 0
