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
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>测试报告: {{.Title}}</title>
    <style>
        :root {
            --success-color: #28a745;
            --failure-color: #dc3545;
            --running-color: #17a2b8;
            --neutral-color: #6c757d;
            --light-bg: #f8f9fa;
            --border-radius: 8px;
            --box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            --transition: all 0.3s ease;
        }
        
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #fff;
            margin: 0;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: white;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
        }
        
        header {
            text-align: center;
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 1px solid #eee;
        }
        
        h1 {
            color: #333;
            margin-bottom: 10px;
        }
        
        h2 {
            color: #444;
            margin: 25px 0 15px;
        }
        
        h3 {
            color: #555;
            margin: 20px 0 10px;
        }
        
        .summary {
            background-color: var(--light-bg);
            border-radius: var(--border-radius);
            padding: 20px;
            margin-bottom: 30px;
            box-shadow: var(--box-shadow);
        }
        
        .summary-details {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin-top: 15px;
        }
        
        .summary-item {
            flex: 1;
            min-width: 200px;
            padding: 10px;
            background-color: white;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
        }
        
        .stats {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin: 20px 0;
        }
        
        .stat-box {
            flex: 1;
            min-width: 150px;
            text-align: center;
            padding: 20px;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
            transition: var(--transition);
        }
        
        .stat-box:hover {
            transform: translateY(-5px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .stat-box h3 {
            margin-top: 0;
            font-size: 1.2em;
        }
        
        .stat-box p {
            font-size: 2em;
            font-weight: bold;
            margin: 10px 0;
        }
        
        .stat-box.total {
            background-color: var(--light-bg);
            color: var(--neutral-color);
        }
        
        .stat-box.passed {
            background-color: rgba(40, 167, 69, 0.1);
            color: var(--success-color);
        }
        
        .stat-box.failed {
            background-color: rgba(220, 53, 69, 0.1);
            color: var(--failure-color);
        }
        
        .progress-container {
            margin: 15px 0;
            background-color: #e9ecef;
            border-radius: 10px;
            height: 10px;
            overflow: hidden;
        }
        
        .progress-bar {
            height: 100%;
            background-color: var(--success-color);
            border-radius: 10px;
        }
        
        .test-results {
            margin-top: 30px;
        }
        
        .test-result {
            margin-bottom: 25px;
            padding: 20px;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
            transition: var(--transition);
        }
        
        .test-result:hover {
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .test-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding-bottom: 10px;
            border-bottom: 1px solid #eee;
        }
        
        .test-title {
            font-size: 1.2em;
            font-weight: bold;
        }
        
        .test-status {
            padding: 5px 10px;
            border-radius: 20px;
            font-weight: bold;
            font-size: 0.9em;
        }
        
        .test-info {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin-bottom: 15px;
        }
        
        .test-info-item {
            flex: 1;
            min-width: 150px;
        }
        
        .test-steps {
            margin-top: 20px;
        }
        
        .step {
            margin: 10px 0;
            padding: 15px;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
            transition: var(--transition);
        }
        
        .step:hover {
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .step-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        
        .step-name {
            font-weight: bold;
        }
        
        .step-status {
            padding: 3px 8px;
            border-radius: 15px;
            font-size: 0.8em;
            font-weight: bold;
        }
        
        .step-details {
            margin-top: 10px;
        }
        
        .error {
            background-color: rgba(220, 53, 69, 0.1);
            color: var(--failure-color);
            padding: 10px;
            border-radius: var(--border-radius);
            margin-top: 10px;
            font-family: monospace;
            white-space: pre-wrap;
        }
        
        .screenshot-container {
            margin-top: 15px;
            text-align: center;
        }
        
        .screenshot {
            max-width: 100%;
            max-height: 300px;
            border: 1px solid #ddd;
            border-radius: var(--border-radius);
            box-shadow: var(--box-shadow);
            cursor: pointer;
            transition: var(--transition);
        }
        
        .screenshot:hover {
            transform: scale(1.02);
        }
        
        .screenshot-modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.9);
            overflow: auto;
        }
        
        .modal-content {
            margin: auto;
            display: block;
            max-width: 90%;
            max-height: 90%;
        }
        
        .close {
            position: absolute;
            top: 15px;
            right: 35px;
            color: #f1f1f1;
            font-size: 40px;
            font-weight: bold;
            cursor: pointer;
        }
        
        .timestamp {
            color: var(--neutral-color);
            font-size: 0.9em;
        }
        
        .duration {
            font-weight: bold;
            color: var(--neutral-color);
        }
        
        .success {
            background-color: rgba(40, 167, 69, 0.1);
            border-left: 4px solid var(--success-color);
        }
        
        .failure {
            background-color: rgba(220, 53, 69, 0.1);
            border-left: 4px solid var(--failure-color);
        }
        
        .running {
            background-color: rgba(23, 162, 184, 0.1);
            border-left: 4px solid var(--running-color);
        }
        
        .status-success {
            background-color: var(--success-color);
            color: white;
        }
        
        .status-failure {
            background-color: var(--failure-color);
            color: white;
        }
        
        .status-running {
            background-color: var(--running-color);
            color: white;
        }
        
        .collapsible {
            cursor: pointer;
        }
        
        .collapsible:after {
            content: ' ▼';
            font-size: 0.8em;
            margin-left: 5px;
        }
        
        .collapsed:after {
            content: ' ▶';
        }
        
        .content {
            max-height: 1000px;
            overflow: hidden;
            transition: max-height 0.3s ease;
        }
        
        .collapsed + .content {
            max-height: 0;
        }
        
        @media (max-width: 768px) {
            .stats, .summary-details, .test-info {
                flex-direction: column;
            }
            
            .stat-box, .summary-item, .test-info-item {
                min-width: 100%;
            }
        }
    </style>
    <script>
        // 页面加载完成后执行
        document.addEventListener('DOMContentLoaded', function() {
            // 初始化可折叠元素
            initCollapsible();
            
            // 初始化截图模态框
            initScreenshotModal();
            
            // 计算通过率并更新进度条
            updateProgressBar();
        });
        
        // 初始化可折叠元素
        function initCollapsible() {
            const collapsibles = document.querySelectorAll('.collapsible');
            
            collapsibles.forEach(function(collapsible) {
                collapsible.addEventListener('click', function() {
                    this.classList.toggle('collapsed');
                });
            });
        }
        
        // 初始化截图模态框
        function initScreenshotModal() {
            // 获取模态框元素
            const modal = document.getElementById('screenshotModal');
            const modalImg = document.getElementById('modalImage');
            const closeBtn = document.getElementsByClassName('close')[0];
            
            // 为所有截图添加点击事件
            const screenshots = document.querySelectorAll('.screenshot');
            screenshots.forEach(function(img) {
                img.onclick = function() {
                    modal.style.display = 'flex';
                    modalImg.src = this.src;
                }
            });
            
            // 关闭模态框
            closeBtn.onclick = function() {
                modal.style.display = 'none';
            }
            
            // 点击模态框外部关闭
            window.onclick = function(event) {
                if (event.target == modal) {
                    modal.style.display = 'none';
                }
            }
        }
        
        // 更新进度条
        function updateProgressBar() {
            const totalSteps = {{.TotalSteps}};
            const passedSteps = {{.PassedSteps}};
            
            if (totalSteps > 0) {
                const passRate = (passedSteps / totalSteps) * 100;
                const progressBar = document.querySelector('.progress-bar');
                progressBar.style.width = passRate + '%';
            }
        }
    </script>
</head>
<body>
    <!-- 截图模态框 -->
    <div id="screenshotModal" class="screenshot-modal">
        <span class="close">&times;</span>
        <img class="modal-content" id="modalImage">
    </div>

    <div class="container">
        <header>
            <h1>测试报告: {{.Title}}</h1>
            <p class="timestamp">生成时间: {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
        </header>
        
        <section class="summary">
            <h2>测试摘要</h2>
            
            <div class="summary-details">
                <div class="summary-item">
                    <p><strong>开始时间:</strong> {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
                </div>
                <div class="summary-item">
                    <p><strong>总测试数:</strong> {{.TotalTests}}</p>
                </div>
                <div class="summary-item">
                    <p><strong>测试环境:</strong> Playwright</p>
                </div>
            </div>
            
            <h3>测试统计</h3>
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
            
            <h3>通过率</h3>
            <div class="progress-container">
                <div class="progress-bar"></div>
            </div>
        </section>
        
        <section class="test-results">
            <h2>测试详情</h2>
            
            {{range .Tests}}
            <div class="test-result {{.Status | lower}}">
                <div class="test-header">
                    <div class="test-title">{{.Name}}</div>
                    <div class="test-status status-{{.Status | lower}}">{{.Status}}</div>
                </div>
                
                <div class="test-info">
                    <div class="test-info-item">
                        <p><strong>开始时间:</strong> {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
                    </div>
                    {{if ne .Duration 0}}
                    <div class="test-info-item">
                        <p><strong>耗时:</strong> <span class="duration">{{.Duration}}</span></p>
                    </div>
                    {{end}}
                    <div class="test-info-item">
                        <p><strong>结果:</strong> {{.Message}}</p>
                    </div>
                </div>
                
                {{if .Steps}}
                <div class="test-steps">
                    <h3 class="collapsible">测试步骤 ({{len .Steps}})</h3>
                    <div class="content">
                        {{range .Steps}}
                        <div class="step {{.Status | lower}}">
                            <div class="step-header">
                                <div class="step-name">{{.Name}}</div>
                                <div class="step-status status-{{.Status | lower}}">{{.Status}}</div>
                            </div>
                            
                            {{if .Message}}
                            <div class="step-details">{{.Message}}</div>
                            {{end}}
                            
                            {{if .Error}}
                            <div class="error">{{.Error}}</div>
                            {{end}}
                            
                            {{if .Screenshot}}
                            <div class="screenshot-container">
                                <p><a href="{{.Screenshot}}" target="_blank">在新窗口中查看截图</a></p>
                                <img class="screenshot" src="{{.Screenshot}}" alt="测试截图">
                            </div>
                            {{end}}
                        </div>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>
            {{end}}
        </section>
    </div>
</body>
</html>
`