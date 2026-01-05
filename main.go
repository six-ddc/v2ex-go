package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/six-ddc/v2ex-go/internal/app"
)

func main() {
	// 创建应用 Model
	model := app.NewModel()

	// 创建 Bubble Tea 程序
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // 使用备用屏幕缓冲区
		tea.WithMouseCellMotion(), // 启用鼠标支持
	)

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
