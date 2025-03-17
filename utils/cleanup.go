package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CleanupOldTestResults 清理旧的测试报告、截图和视频，只保留最新的文件
func CleanupOldTestResults() error {
	fmt.Println("开始清理旧的测试结果...")

	// 获取当前工作目录的绝对路径
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}

	// 清理报告目录
	reportsDir := filepath.Join(cwd, "reports")
	if err := cleanupDirectory(reportsDir, ".html", 1); err != nil {
		return fmt.Errorf("清理报告目录失败: %w", err)
	}

	// 清理截图目录
	screenshotsDir := filepath.Join(cwd, "screenshots")
	if err := cleanupDirectory(screenshotsDir, ".png", 3); err != nil {
		return fmt.Errorf("清理截图目录失败: %w", err)
	}

	// 清理视频目录
	videosDir := filepath.Join(cwd, "videos")
	if err := cleanupDirectory(videosDir, ".webm", 1); err != nil {
		return fmt.Errorf("清理视频目录失败: %w", err)
	}

	fmt.Println("清理完成，只保留最新的测试结果")
	return nil
}

// cleanupDirectory 清理指定目录中的文件，只保留指定数量的最新文件
func cleanupDirectory(dirPath string, fileExt string, keepCount int) error {
	// 设置最大重试次数和重试间隔
	const maxRetries = 3
	retryInterval := time.Second * 2
	// 打印清理信息
	fmt.Printf("清理目录 %s，保留 %d 个最新的%s文件\n", dirPath, keepCount, fileExt)
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建它
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dirPath, err)
		}
		return nil // 新目录，没有文件需要清理
	}

	// 获取目录中的所有文件
	var files []fs.DirEntry
	var err error
	if files, err = os.ReadDir(dirPath); err != nil {
		return fmt.Errorf("读取目录 %s 失败: %w", dirPath, err)
	}

	// 筛选出符合扩展名的文件并获取它们的信息
	type fileInfo struct {
		path    string
		modTime time.Time
	}
	var fileInfos []fileInfo

	for _, file := range files {
		if file.IsDir() {
			continue // 跳过子目录
		}

		fileName := file.Name()
		if !strings.HasSuffix(strings.ToLower(fileName), fileExt) {
			continue // 跳过不符合扩展名的文件
		}

		filePath := filepath.Join(dirPath, fileName)
		info, err := os.Stat(filePath)
		if err != nil {
			continue // 跳过无法获取信息的文件
		}

		fileInfos = append(fileInfos, fileInfo{
			path:    filePath,
			modTime: info.ModTime(),
		})
	}

	// 按修改时间排序（最新的在前面）
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].modTime.After(fileInfos[j].modTime)
	})

	// 删除多余的旧文件
	if len(fileInfos) > keepCount {
		for i := keepCount; i < len(fileInfos); i++ {
			filePath := fileInfos[i].path
			fmt.Printf("尝试删除旧文件: %s\n", filePath)

			// 尝试多次删除文件
			var deleteErr error
			for retry := 0; retry < maxRetries; retry++ {
				// 尝试打开文件以确保它没有被其他进程使用
				file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
				if err != nil {
					if os.IsNotExist(err) {
						// 文件已经不存在，视为删除成功
						deleteErr = nil
						break
					}
					// 如果文件被占用，等待后重试
					fmt.Printf("文件 %s 可能被占用，等待重试 (%d/%d)\n", filePath, retry+1, maxRetries)
					time.Sleep(retryInterval)
					continue
				}
				file.Close()

				// 尝试删除文件
				if err := os.Remove(filePath); err != nil {
					deleteErr = err
					fmt.Printf("删除失败，等待重试 (%d/%d): %v\n", retry+1, maxRetries, err)
					time.Sleep(retryInterval)
					continue
				}

				// 验证文件是否确实被删除
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					fmt.Printf("成功删除文件: %s\n", filePath)
					deleteErr = nil
					break
				}
			}

			// 如果所有重试都失败，记录错误
			if deleteErr != nil {
				fmt.Printf("警告: 经过多次尝试后仍无法删除文件 %s: %v\n", filePath, deleteErr)
			}
		}
	} else {
		fmt.Printf("目录 %s 中文件数量(%d)未超过保留数量(%d)，无需清理\n", dirPath, len(fileInfos), keepCount)
	}

	return nil
}
