package utils

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestStep 表示测试步骤
type TestStep struct {
	Name      string
	Status    string // "Success", "Failure", "Running"
	Message   string
	Error     error
	Timestamp time.Time
	Screenshot string
}

// Test 表示一个测试
type Test struct {
	Name      string
	Status    string // "Success", "Failure", "Running"
	Message   string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Steps     []TestStep
}

// ReportManager 管理测试报告
type ReportManager struct {
	Title     string
	StartTime time.Time
	Tests     []Test
	currentTest *Test
	currentStep *TestStep
}

// NewReportManager 创建一个新的报告管理器
func NewReportManager(title string) *ReportManager {
	return &ReportManager{
		Title:     title,
		StartTime: time.Now(),
		Tests:     []Test{},
	}
}

// StartTest 开始一个新的测试
func (r *ReportManager) StartTest(name string) {
	test := Test{
		Name:      name,
		Status:    "Running",
		StartTime: time.Now(),
		Steps:     []TestStep{},
	}
	r.Tests = append(r.Tests, test)
	r.currentTest = &r.Tests[len(r.Tests)-1]
}

// StartStep 开始一个新的测试步骤
func (r *ReportManager) StartStep(name string) {
	if r.currentTest == nil {
		return
	}
	step := TestStep{
		Name:      name,
		Status:    "Running",
		Timestamp: time.Now(),
	}
	r.currentTest.Steps = append(r.currentTest.Steps, step)
	r.currentStep = &r.currentTest.Steps[len(r.currentTest.Steps)-1]
}

// EndStepSuccess 标记当前步骤为成功
func (r *ReportManager) EndStepSuccess(message string) {
	if r.currentStep == nil {
		return
	}
	r.currentStep.Status = "Success"
	r.currentStep.Message = message
}

// EndStepFailure 标记当前步骤为失败
func (r *ReportManager) EndStepFailure(message string, err error, screenshot string) {
	if r.currentStep == nil {
		return
	}
	r.currentStep.Status = "Failure"
	r.currentStep.Message = message
	r.currentStep.Error = err
	r.currentStep.Screenshot = screenshot
}

// LogSuccess 标记当前测试为成功
func (r *ReportManager) LogSuccess(message string, duration time.Duration) {
	if r.currentTest == nil {
		return
	}
	r.currentTest.Status = "Success"
	r.currentTest.Message = message
	r.currentTest.EndTime = time.Now()
	r.currentTest.Duration = duration
}

// LogFailure 标记当前测试为失败
func (r *ReportManager) LogFailure(message string, duration time.Duration) {
	if r.currentTest == nil {
		return
	}
	r.currentTest.Status = "Failure"
	r.currentTest.Message = message
	r.currentTest.EndTime = time.Now()
	r.currentTest.Duration = duration
}

// GenerateReport 生成HTML测试报告
func (r *ReportManager) GenerateReport() (string, error) {
	// 创建报告目录
	reportDir := "./reports"
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		os.MkdirAll(reportDir, 0755)
	}

	// 生成报告文件名
	timestamp := time.Now().Format("20060102-150405")
	reportPath := filepath.Join(reportDir, fmt.Sprintf("report-%s.html", timestamp))

	// 创建报告文件
	file, err := os.Create(reportPath)
	if err != nil {
		return "", fmt.Errorf("无法创建报告文件: %w", err)
	}
	defer file.Close()

	// 准备报告模板
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(reportTemplate))

	// 计算统计信息
	totalSteps := 0
	passedSteps := 0
	failedSteps := 0

	for _, test := range r.Tests {
		for _, step := range test.Steps {
			totalSteps++
			if step.Status == "Success" {
				passedSteps++
			} else if step.Status == "Failure" {
				failedSteps++
			}
		}
	}

	// 准备模板数据
	data := struct {
		Title       string
		StartTime   time.Time
		Tests       []Test
		TotalTests  int
		TotalSteps  int
		PassedSteps int
		FailedSteps int
	}{
		Title:       r.Title,
		StartTime:   r.StartTime,
		Tests:       r.Tests,
		TotalTests:  len(r.Tests),
		TotalSteps:  totalSteps,
		PassedSteps: passedSteps,
		FailedSteps: failedSteps,
	}

	// 执行模板
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("无法生成报告: %w", err)
	}

	return reportPath, nil
}

// HTML报告模板
const reportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>测试报告: {{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .summary { margin: 20px 0; padding: 10px; background-color: #f5f5f5; border-radius: 5px; }
        .test-result { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .success { background-color: #dff0d8; border: 1px solid #d6e9c6; }
        .failure { background-color: #f2dede; border: 1px solid #ebccd1; }
        .running { background-color: #d9edf7; border: 1px solid #bce8f1; }
        .details { margin-top: 5px; font-size: 0.9em; }
        .timestamp { color: #777; font-size: 0.8em; }
        .duration { font-weight: bold; }
        .error { color: #a94442; margin-top: 5px; }
        .step { margin-left: 20px; padding: 8px; margin-top: 5px; border-radius: 3px; }
        .screenshot { max-width: 100%; margin-top: 10px; border: 1px solid #ddd; }
        .step-details { margin-top: 5px; }
        .stats { display: flex; justify-content: space-between; margin-bottom: 10px; }
        .stat-box { flex: 1; text-align: center; padding: 10px; margin: 0 5px; border-radius: 5px; }
        .stat-box.total { background-color: #f5f5f5; }
        .stat-box.passed { background-color: #dff0d8; }
        .stat-box.failed { background-color: #f2dede; }
    </style>
</head>
<body>
    <h1>测试报告: {{.Title}}</h1>
    
    <div class="summary">
        <p><strong>开始时间:</strong> {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
        <p><strong>总测试数:</strong> {{.TotalTests}}</p>
        
        <div class="stats">
            <div class="stat-box total">
                <h3>总步骤数</h3>
                <p>{{.TotalSteps}}</p>
            </div>
            <div class="stat-box passed">
                <h3>通过</h3>
                <p>{{.PassedSteps}}</p>
            </div>
            <div class="stat-box failed">
                <h3>失败</h3>
                <p>{{.FailedSteps}}</p>
            </div>
        </div>
    </div>
    
    {{range .Tests}}
    <div class="test-result {{.Status | lower}}">
        <h3>{{.Name}}</h3>
        <p>状态: <strong>{{.Status}}</strong></p>
        <p class="details">{{.Message}}</p>
        <p class="timestamp">时间: {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
        {{if ne .Duration 0}}
        <p class="duration">耗时: {{.Duration}}</p>
        {{end}}
        
        {{if .Steps}}
        <h4>测试步骤:</h4>
        {{range .Steps}}
        <div class="step {{.Status | lower}}">
            <strong>{{.Name}}</strong> - {{.Status}}
            {{if .Message}}
            <p class="step-details">{{.Message}}</p>
            {{end}}
            {{if .Error}}
            <p class="error">错误: {{.Error}}</p>
            {{end}}
            {{if .Screenshot}}
            <p><a href="{{.Screenshot}}" target="_blank">查看截图</a></p>
            <img class="screenshot" src="{{.Screenshot}}" alt="截图">
            {{end}}
        </div>
        {{end}}
        {{end}}
    </div>
    {{end}}
</body>
</html>
`