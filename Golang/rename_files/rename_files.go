package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 获取命令行参数
	folder := os.Args[1]

	// 读取文件夹下所有文件
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历所有文件,去掉前缀
	for _, f := range files {
		fileName := f.Name()
		newName := strings.TrimPrefix(fileName, fileName[:6])
		oldPath := filepath.Join(folder, fileName)
		newPath := filepath.Join(folder, newName)

		// 重命名文件
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("去除前缀完成!")
}
