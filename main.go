package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/wan/playwright-go-demo/config"
	"github.com/wan/playwright-go-demo/pages"
	"github.com/wan/playwright-go-demo/utils"
)

func main() {
	// 加载配置文件
	configPath := "./config/config.json"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("无法启动Playwright: %v", err)
	}
	defer pw.Stop()

	// 根据配置选择浏览器类型
	var browserType playwright.BrowserType
	switch cfg.Browser.Type {
	case "firefox":
		browserType = pw.Firefox
	case "webkit":
		browserType = pw.WebKit
	default:
		browserType = pw.Chromium
	}

	// 创建浏览器实例
	browser, err := browserType.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Browser.Headless),
		SlowMo:   playwright.Float(float64(cfg.Browser.SlowMo)),
	})
	if err != nil {
		log.Fatalf("无法启动浏览器: %v", err)
	}
	defer browser.Close()

	// 创建上下文
	contextOptions := playwright.BrowserNewContextOptions{
		RecordVideo: &playwright.RecordVideo{
			Dir: "./videos",
		},
	}

	// 如果配置了最大化，设置视口大小为最大
	if cfg.Browser.Maximized {
		// 设置一个足够大的视口大小来模拟最大化
		contextOptions.Viewport = &playwright.Size{
			Width:  1920,
			Height: 1080,
		}
	}

	context, err := browser.NewContext(contextOptions)
	if err != nil {
		log.Fatalf("无法创建浏览器上下文: %v", err)
	}
	defer context.Close()

	// 创建页面
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("无法创建页面: %v", err)
	}

	// 确保截图目录存在
	screenshotDir := "./screenshots"
	if _, err := os.Stat(screenshotDir); os.IsNotExist(err) {
		os.MkdirAll(screenshotDir, 0755)
	}

	// 确保视频目录存在
	videoDir := "./videos"
	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		os.MkdirAll(videoDir, 0755)
	}

	// 初始化测试报告
	reportManager := utils.NewReportManager("登录测试")
	reportManager.StartTest("登录测试")

	// 执行测试
	testStart := time.Now()

	try := func() bool {
		// 创建登录页面对象
		loginPage := pages.NewLoginPage(page)
		// 设置登录URL
		loginPage.SetLoginURL(cfg.Login.URL)

		// 步骤1: 导航到登录页面
		reportManager.StartStep("导航到登录页面")
		if err := loginPage.Navigate(); err != nil {
			// 失败时截图
			screenshotPath := filepath.Join(screenshotDir, "navigate_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("导航到登录页面失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功导航到登录页面")

		// 步骤2: 执行登录
		reportManager.StartStep("执行登录操作")
		if err := loginPage.Login(cfg.Login.Username, cfg.Login.Password); err != nil {
			// 失败时截图
			screenshotPath := filepath.Join(screenshotDir, "login_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("登录失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功执行登录操作")

		// 步骤3: 验证登录成功
		reportManager.StartStep("验证登录结果")
		if success, err := loginPage.VerifyLoginSuccess(); err != nil || !success {
			// 失败时截图
			screenshotPath := filepath.Join(screenshotDir, "verification_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("验证登录失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功验证登录结果")

		return true
	}

	success := try()
	testDuration := time.Since(testStart)

	// 完成测试报告
	if success {
		reportManager.LogSuccess("登录测试成功", testDuration)
	} else {
		reportManager.LogFailure("登录测试失败", testDuration)
	}

	// 生成测试报告
	reportPath, err := reportManager.GenerateReport()
	if err != nil {
		log.Fatalf("生成测试报告失败: %v", err)
	}

	fmt.Printf("测试完成，报告已生成: %s\n", reportPath)
}
