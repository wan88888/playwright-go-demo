# Playwright-Go 自动化测试框架

这是一个基于 Playwright-Go 的 Web 自动化测试框架，提供了多浏览器支持、页面对象模型(POM)、测试报告生成、截图和视频录制等功能。

## 项目概述

本框架使用 Go 语言和 Playwright 实现，旨在提供一个简单易用、功能完善的 Web 自动化测试解决方案。框架采用页面对象模型(POM)设计模式，使测试代码更加模块化和可维护。

### 主要特性

- **多浏览器支持**：支持 Chromium、Firefox 和 WebKit 浏览器
- **页面对象模型(POM)**：将页面元素和操作封装，提高代码复用性
- **自动截图**：测试失败时自动截图，方便问题定位
- **视频录制**：自动录制测试过程，便于回放分析
- **HTML 测试报告**：生成美观、详细的测试报告
- **配置化**：通过 JSON 配置文件灵活设置测试参数
- **自动清理**：自动清理旧的测试结果，保持工作目录整洁

## 架构设计

### 目录结构

```
.
├── config/             # 配置文件目录
│   ├── config.go      # 配置加载和处理逻辑
│   └── config.json    # 测试配置文件
├── pages/             # 页面对象模型目录
│   └── login_page.go  # 登录页面对象
├── utils/             # 工具函数目录
│   ├── cleanup.go     # 清理旧测试结果
│   ├── report_manager.go # 测试报告生成
│   └── screenshot.go  # 截图工具
├── main.go            # 主程序入口
├── go.mod             # Go 模块定义
└── go.sum             # 依赖版本锁定
```

### 核心组件

1. **配置管理**：通过 `config` 包加载和管理测试配置
2. **页面对象**：在 `pages` 包中定义页面元素和操作
3. **测试报告**：使用 `utils/report_manager.go` 生成 HTML 测试报告
4. **截图工具**：使用 `utils/screenshot.go` 在测试失败时捕获截图
5. **清理工具**：使用 `utils/cleanup.go` 清理旧的测试结果

## 功能详解

### 1. 多浏览器支持

框架支持在 Chromium、Firefox 和 WebKit 浏览器上运行测试，可以通过配置文件指定要使用的浏览器类型、无头模式、慢动作模式等参数。

```go
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
```

### 2. 页面对象模型(POM)

框架采用页面对象模型设计模式，将页面元素和操作封装在特定的页面对象中，提高代码的可读性和可维护性。

```go
// LoginPage 表示登录页面对象
type LoginPage struct {
    page     playwright.Page
    loginURL string
}

// Login 执行登录操作
func (l *LoginPage) Login(username, password string) error {
    // 输入用户名
    if err := l.page.Fill("#username", username); err != nil {
        return fmt.Errorf("无法输入用户名: %w", err)
    }
    // ...
}
```

### 3. 自动截图

在测试失败时，框架会自动捕获页面截图，帮助快速定位问题。

```go
if err := loginPage.Navigate(); err != nil {
    // 失败时截图
    screenshotPath := filepath.Join(browserScreenshotDir, "navigate_failure.png")
    utils.TakeScreenshot(page, screenshotPath)
    reportManager.EndStepFailure("导航到登录页面失败", err, screenshotPath)
    return false
}
```

### 4. 视频录制

框架支持自动录制测试过程，生成视频文件，便于回放分析。

```go
contextOptions := playwright.BrowserNewContextOptions{
    RecordVideo: &playwright.RecordVideo{
        Dir: filepath.Join(videoDir, browserConfig.Type),
    },
}
```

### 5. HTML 测试报告

框架会生成美观、详细的 HTML 测试报告，包含测试步骤、状态、截图、错误信息等内容。

```go
// 生成测试报告
reportPath, err := reportManager.GenerateReport()
if err != nil {
    log.Fatalf("生成测试报告失败: %v", err)
}
```

### 6. 配置化

通过 JSON 配置文件，可以灵活设置浏览器类型、登录信息等测试参数。

```json
{
  "browsers": [
    {
      "type": "chromium",
      "headless": true,
      "slowMo": 0,
      "maximized": true
    }
  ],
  "login": {
    "username": "tomsmith",
    "password": "SuperSecretPassword!",
    "url": "http://the-internet.herokuapp.com/login"
  }
}
```

### 7. 自动清理

框架会自动清理旧的测试结果，只保留最新的几个文件，保持工作目录整洁。

```go
// 清理旧的测试结果
if err := utils.CleanupOldTestResults(); err != nil {
    log.Printf("警告: 清理旧测试结果失败: %v", err)
}
```

## 配置说明

配置文件位于 `config/config.json`，包含以下主要配置项：

### 浏览器配置

```json
"browsers": [
  {
    "type": "chromium",  // 浏览器类型：chromium, firefox, webkit
    "headless": true,    // 是否无头模式
    "slowMo": 0,         // 慢动作模式，毫秒
    "maximized": true    // 是否最大化
  }
]
```

### 登录配置

```json
"login": {
  "username": "tomsmith",           // 用户名
  "password": "SuperSecretPassword!", // 密码
  "url": "http://the-internet.herokuapp.com/login", // 登录URL
  "invalid_username": "invaliduser",  // 无效用户名
  "invalid_password": "invalidpass"   // 无效密码
}
```

## 使用方法

### 前置条件

1. 安装 Go 1.18 或更高版本
2. 安装 Playwright 依赖

```bash
go mod download
go run github.com/playwright-community/playwright-go/cmd/playwright install
```

### 运行测试

```bash
go run main.go
```

### 查看测试报告

测试完成后，可以在 `reports` 目录中找到生成的 HTML 测试报告。

### 自定义测试

1. 修改 `config/config.json` 配置文件，设置浏览器类型、登录信息等
2. 在 `pages` 目录中添加新的页面对象
3. 在 `main.go` 中添加新的测试场景

## 示例

### 添加新的页面对象

```go
// pages/dashboard_page.go
package pages

import (
    "fmt"
    "github.com/playwright-community/playwright-go"
)

type DashboardPage struct {
    page playwright.Page
}

func NewDashboardPage(page playwright.Page) *DashboardPage {
    return &DashboardPage{page: page}
}

func (d *DashboardPage) VerifyDashboardLoaded() (bool, error) {
    // 验证仪表盘页面是否加载成功
    return d.page.IsVisible(".dashboard-header"), nil
}
```

### 添加新的测试场景

```go
// 测试场景4: 登录后验证仪表盘
reportManager.StartStep("验证仪表盘")
if err := loginPage.Login(loginConfig.Username, loginConfig.Password); err != nil {
    // 处理错误...
}

dashboardPage := pages.NewDashboardPage(page)
if loaded, err := dashboardPage.VerifyDashboardLoaded(); err != nil || !loaded {
    // 处理错误...
}
reportManager.EndStepSuccess("成功验证仪表盘加载")
```

## 贡献

欢迎提交 Issue 和 Pull Request 来完善本框架。

## 许可证

MIT