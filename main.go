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
	// 清理旧的测试结果
	if err := utils.CleanupOldTestResults(); err != nil {
		log.Printf("警告: 清理旧测试结果失败: %v", err)
	}

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

	// 遍历所有配置的浏览器，分别执行测试
	for _, browserConfig := range cfg.Browsers {
		// 为每个浏览器创建单独的测试报告
		reportManager := utils.NewReportManager(fmt.Sprintf("%s浏览器登录测试", browserConfig.Type))

		// 执行特定浏览器的测试
		runTestWithBrowser(pw, browserConfig, cfg.Login, screenshotDir, videoDir, reportManager)
	}
}

// runTestWithBrowser 使用特定浏览器执行测试
func runTestWithBrowser(pw *playwright.Playwright, browserConfig config.BrowserConfig, loginConfig config.LoginConfig, screenshotDir, videoDir string, reportManager *utils.ReportManager) {
	// 根据配置选择浏览器类型
	var browserType playwright.BrowserType
	switch browserConfig.Type {
	case "firefox":
		browserType = pw.Firefox
	case "webkit":
		browserType = pw.WebKit
	default:
		browserType = pw.Chromium
	}

	fmt.Printf("开始使用 %s 浏览器执行测试\n", browserConfig.Type)

	// 创建浏览器实例
	browser, err := browserType.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(browserConfig.Headless),
		SlowMo:   playwright.Float(float64(browserConfig.SlowMo)),
	})
	if err != nil {
		log.Printf("无法启动 %s 浏览器: %v", browserConfig.Type, err)
		return
	}
	defer browser.Close()

	// 创建上下文
	contextOptions := playwright.BrowserNewContextOptions{
		RecordVideo: &playwright.RecordVideo{
			Dir: filepath.Join(videoDir, browserConfig.Type), // 为每个浏览器创建单独的视频目录
		},
	}

	// 如果配置了最大化，设置视口大小为最大
	if browserConfig.Maximized {
		// 设置一个足够大的视口大小来模拟最大化
		contextOptions.Viewport = &playwright.Size{
			Width:  1920,
			Height: 1080,
		}
	}

	// 确保浏览器特定的视频目录存在
	browserVideoDir := filepath.Join(videoDir, browserConfig.Type)
	if _, err := os.Stat(browserVideoDir); os.IsNotExist(err) {
		os.MkdirAll(browserVideoDir, 0755)
	}

	// 确保浏览器特定的截图目录存在
	browserScreenshotDir := filepath.Join(screenshotDir, browserConfig.Type)
	if _, err := os.Stat(browserScreenshotDir); os.IsNotExist(err) {
		os.MkdirAll(browserScreenshotDir, 0755)
	}

	context, err := browser.NewContext(contextOptions)
	if err != nil {
		log.Printf("无法创建 %s 浏览器上下文: %v", browserConfig.Type, err)
		return
	}
	defer context.Close()

	// 创建页面
	page, err := context.NewPage()
	if err != nil {
		log.Printf("无法创建 %s 浏览器页面: %v", browserConfig.Type, err)
		return
	}
	reportManager.StartTest("登录测试")

	// 执行测试
	testStart := time.Now()

	try := func() bool {
		// 创建登录页面对象
		loginPage := pages.NewLoginPage(page)
		// 设置登录URL
		loginPage.SetLoginURL(loginConfig.URL)

		// 步骤1: 导航到登录页面
		reportManager.StartStep("导航到登录页面")
		if err := loginPage.Navigate(); err != nil {
			// 失败时截图
			screenshotPath := filepath.Join(browserScreenshotDir, "navigate_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("导航到登录页面失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功导航到登录页面")

		// 测试场景1: 使用错误的用户名登录
		reportManager.StartStep("测试错误用户名登录")
		if err := loginPage.Login("wrong_username", loginConfig.Password); err != nil {
			screenshotPath := filepath.Join(browserScreenshotDir, "wrong_username_input_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("输入错误用户名失败", err, screenshotPath)
			return false
		}
		if failed, err := loginPage.VerifyLoginFailed(); err != nil || !failed {
			screenshotPath := filepath.Join(browserScreenshotDir, "wrong_username_verify_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("验证错误用户名失败场景失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功验证错误用户名登录失败场景")

		// 测试场景2: 使用错误的密码登录
		reportManager.StartStep("测试错误密码登录")
		if err := loginPage.Login(loginConfig.Username, "wrong_password"); err != nil {
			screenshotPath := filepath.Join(browserScreenshotDir, "wrong_password_input_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("输入错误密码失败", err, screenshotPath)
			return false
		}
		if failed, err := loginPage.VerifyLoginFailed(); err != nil || !failed {
			screenshotPath := filepath.Join(browserScreenshotDir, "wrong_password_verify_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("验证错误密码失败场景失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功验证错误密码登录失败场景")

		// 测试场景3: 使用正确的凭据登录
		reportManager.StartStep("测试正确凭据登录")
		if err := loginPage.Login(loginConfig.Username, loginConfig.Password); err != nil {
			screenshotPath := filepath.Join(browserScreenshotDir, "login_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("登录失败", err, screenshotPath)
			return false
		}
		if success, err := loginPage.VerifyLoginSuccess(); err != nil || !success {
			screenshotPath := filepath.Join(browserScreenshotDir, "verification_failure.png")
			utils.TakeScreenshot(page, screenshotPath)
			reportManager.EndStepFailure("验证登录失败", err, screenshotPath)
			return false
		}
		reportManager.EndStepSuccess("成功验证正确凭据登录")

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
