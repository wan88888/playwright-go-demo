package utils

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
)

// TakeScreenshot 捕获页面截图并保存到指定路径
func TakeScreenshot(page playwright.Page, path string) error {
	// 设置截图选项
	options := playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	}

	// 捕获截图
	_, err := page.Screenshot(options)
	if err != nil {
		return fmt.Errorf("截图失败: %w", err)
	}

	return nil
}

// TakeElementScreenshot 捕获特定元素的截图并保存到指定路径
func TakeElementScreenshot(page playwright.Page, selector string, path string) error {
	// 查找元素
	element, err := page.QuerySelector(selector)
	if err != nil {
		return fmt.Errorf("找不到元素 '%s': %w", selector, err)
	}
	if element == nil {
		return fmt.Errorf("元素 '%s' 不存在", selector)
	}

	// 设置截图选项
	options := playwright.ElementHandleScreenshotOptions{
		Path: playwright.String(path),
	}

	// 捕获元素截图
	_, err = element.Screenshot(options)
	if err != nil {
		return fmt.Errorf("元素截图失败: %w", err)
	}

	return nil
}