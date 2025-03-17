package pages

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// LoginPage 表示登录页面对象
type LoginPage struct {
	page     playwright.Page
	loginURL string
}

// NewLoginPage 创建一个新的登录页面对象
func NewLoginPage(page playwright.Page) *LoginPage {
	return &LoginPage{
		page:     page,
		loginURL: "http://the-internet.herokuapp.com/login", // 默认URL，将被配置文件中的URL覆盖
	}
}

// SetLoginURL 设置登录URL
func (l *LoginPage) SetLoginURL(url string) {
	l.loginURL = url
}

// Navigate 导航到登录页面
func (l *LoginPage) Navigate() error {
	_, err := l.page.Goto(l.loginURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	return err
}

// Login 执行登录操作
func (l *LoginPage) Login(username, password string) error {
	// 输入用户名
	if err := l.page.Fill("#username", username); err != nil {
		return fmt.Errorf("无法输入用户名: %w", err)
	}

	// 输入密码
	if err := l.page.Fill("#password", password); err != nil {
		return fmt.Errorf("无法输入密码: %w", err)
	}

	// 点击登录按钮
	if err := l.page.Click("button[type=\"submit\"]"); err != nil {
		return fmt.Errorf("无法点击登录按钮: %w", err)
	}

	// 等待页面加载完成
	if err := l.page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		return fmt.Errorf("等待页面加载超时: %w", err)
	}

	return nil
}

// VerifyLoginSuccess 验证登录是否成功
func (l *LoginPage) VerifyLoginSuccess() (bool, error) {
	// 等待成功消息出现
	successLocator := l.page.Locator(".flash.success")
	if err := successLocator.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(5000),
	}); err != nil {
		return false, fmt.Errorf("未找到成功消息: %w", err)
	}

	// 检查是否存在登出按钮
	logoutButton, err := l.page.IsVisible("a[href=\"/logout\"]")
	if err != nil {
		return false, fmt.Errorf("检查登出按钮失败: %w", err)
	}

	return logoutButton, nil
}

// Logout 执行登出操作
func (l *LoginPage) Logout() error {
	// 点击登出按钮
	if err := l.page.Click("a[href=\"/logout\"]"); err != nil {
		return fmt.Errorf("无法点击登出按钮: %w", err)
	}

	// 等待页面加载完成
	if err := l.page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		return fmt.Errorf("等待页面加载超时: %w", err)
	}

	return nil
}

// WaitForTimeout 等待指定时间
func (l *LoginPage) WaitForTimeout(ms int) {
	l.page.WaitForTimeout(float64(ms))
}

// VerifyLoginFailed 验证登录失败场景
func (l *LoginPage) VerifyLoginFailed() (bool, error) {
	// 等待错误消息出现
	errorLocator := l.page.Locator(".flash.error")
	if err := errorLocator.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(5000),
	}); err != nil {
		return false, fmt.Errorf("未找到错误消息: %w", err)
	}

	// 检查是否仍在登录页面（通过登录按钮是否可见来判断）
	loginButton, err := l.page.IsVisible("button[type=\"submit\"]")
	if err != nil {
		return false, fmt.Errorf("检查登录按钮失败: %w", err)
	}

	return loginButton, nil
}
