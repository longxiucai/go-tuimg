package main

import (
	"fmt"
	"os"

	"github.com/longxiucai/go-tuimg/viewer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "用法: %s <图片路径>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "快捷键:\n")
		fmt.Fprintf(os.Stderr, "  +/-     : 放大/缩小图片\n")
		fmt.Fprintf(os.Stderr, "  方向键  : 上下左右移动图片\n")
		fmt.Fprintf(os.Stderr, "  0       : 重置缩放和位置\n")
		fmt.Fprintf(os.Stderr, "  Q       : 退出\n")
		fmt.Fprintf(os.Stderr, "鼠标:\n")
		fmt.Fprintf(os.Stderr, "  拖动    : 移动图片\n")
		fmt.Fprintf(os.Stderr, "  滚轮    : 放大/缩小图片\n")
		os.Exit(1)
	}

	path := os.Args[1]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "文件不存在: %s\n", path)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "正在加载图片: %s\n", path)
	viewer.New().Run(path)
}