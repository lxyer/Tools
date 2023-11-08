package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const outDir = "out"

var (
	totalCount     int
	remainingCount int
)

func main() {

	// 获取输入目录
	exePath, _ := os.Executable()
	inputDir := filepath.Dir(exePath)

	// 统计总文件数
	totalCount = 0
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			totalCount++
		}
		return nil
	})

	remainingCount = totalCount

	fmt.Println("开始处理文件")

	// 遍历处理文件
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			processFile(path)
		}
		return nil
	})

	// 完成日志
	log.Printf("总文件数: %d, 剩余文件数: %d", totalCount, remainingCount)
	fmt.Println("处理完成!")
}

func processFile(filePath string) {
	startTime := time.Now()
	// 判断out文件是否存在,存在就删除
	outFileName := filepath.Base(filePath)
	outPath := filepath.Join(outDir, outFileName)

	// 日志打印文件名
	fileName := filepath.Base(filePath)

	if !isMediaFile(filePath) {
		log.Printf("跳过非媒体文件: %s", fileName)
		return
	}

	remainingCount--

	log.Printf("开始处理文件: %s", fileName)

	// ffmpeg处理
	outDir := filepath.Join(filepath.Dir(filePath), "out")
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		return
	}

	outPath = filepath.Join(outDir, filepath.Base(filePath))

	if fileExists(outPath) {
		err := os.Remove(outPath)
		if err != nil {
			return
		}
		log.Printf("删除已存在临时文件: %s", outFileName)
	}

	cmd := exec.Command("ffmpeg", "-i", filePath, "-filter:a", "atempo=2.5", outPath)
	err = cmd.Run()
	if err != nil {
		log.Printf("处理文件错误: %s", err)
		remainingCount++
		return
	}

	// 获取原文件信息
	origStat, err := os.Stat(filePath)
	if err != nil {
		log.Printf("获取文件信息错误: %s", err)
		return
	}
	origSize := origStat.Size()

	// 获取处理后文件信息
	newStat, err := os.Stat(outPath)
	if err != nil {
		log.Printf("获取文件信息错误: %s", err)
		return
	}
	newSize := newStat.Size()

	// 获取原文件大小(MB)
	origSizeMB := float64(origSize) / 1024 / 1024

	// 获取新文件大小(MB)
	newSizeMB := float64(newSize) / 1024 / 1024

	// 计算缩小的MB大小
	sizeReducedMB := origSizeMB - newSizeMB

	// 输出信息
	sizeReduced := 100 - int(newSizeMB/origSizeMB*100)
	log.Printf("文件大小: %.2f MB, 压缩率: %d%%", newSizeMB, sizeReduced)
	log.Printf("文件减小: %.2f MB", sizeReducedMB)
	// 在函数的最后，记录结束时间并计算耗时
	duration := time.Since(startTime)
	minutes := duration / time.Minute
	seconds := duration % time.Minute / time.Second
	log.Printf("文件处理耗时: %d分%d秒", minutes, seconds)
	// 删除原文件
	err = os.Remove(filePath)
	if err != nil {
		remainingCount++
		log.Printf("删除原文件失败: %s", err)
	}

	log.Printf("处理完成,删除文件: %s", fileName)

}

// 判断是否媒体文件
func isMediaFile(path string) bool {

	ext := filepath.Ext(path)

	videoExts := []string{".mp4", ".avi", ".wmv", ".mpg", ".mpeg", ".mkv", ".rmvb", ".flv", ".mov",
		".webm", ".vob", ".m4v", ".mts", ".m2ts", ".ts", ".qt", ".yuv", ".mxf"}

	audioExts := []string{".mp3", ".wav", ".ogg", ".aac", ".flac",
		".ape", ".aiff", ".wma", ".amr", ".m4a"}

	if contains(videoExts, ext) {
		return true
	}

	if contains(audioExts, ext) {
		return true
	}

	return false

}

// 工具函数
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
