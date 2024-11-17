package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 配置加载-初始化
	// configTest := &Config{}
	// 保存默认配置
	// saveConfig(configTest)
	// 创建app
	duc := app.New()
	ducWindow := duc.NewWindow("Ali-DDNS Client")

	// 更新状态标签
	statusLabel := widget.NewLabel("Ready")
	// 更新时间标签
	updateTimeLabel := widget.NewLabel("Last update: ")

	// 更新按钮点击事件
	updateButton := widget.NewButton("Update IP", func() {
		_err := refresh()
		if _err != nil {
			statusLabel.SetText(fmt.Sprintf("Error: %v", _err))
		} else {
			// setLastIP(currentIP, "last_ip.txt")
			// statusLabel.SetText("IP updated successfully.")
			updateTimeLabel.SetText(fmt.Sprintf("Last update: %s", time.Now().Format("2006-01-02 15:04:05")))
		}
	})

	currentIP, _ := getPublicIP()
	statusLabel.SetText(fmt.Sprintf("Current IP: %s", currentIP))

	fmt.Print("go func ")
	// 定时任务：每5分钟检查一次IP变化
	go func() {
		for range time.Tick(5 * time.Minute) {
			fmt.Println("tick")
			currentIP, _ := getPublicIP()
			// lastIP, _ := getLastIP("last_ip.txt")
			// if lastIP != currentIP {
			// updateDDNS(config, currentIP)
			// setLastIP(currentIP, "last_ip.txt")
			// }
			fmt.Println(currentIP)
			_err := refresh()
			if _err != nil {
				statusLabel.SetText(fmt.Sprintf("Error: %v", _err))
			} else {
				// setLastIP(currentIP, "last_ip.txt")
				statusLabel.SetText(fmt.Sprintf("Current IP: %s", currentIP))
				updateTimeLabel.SetText(fmt.Sprintf("Last update: %s", time.Now().Format("2006-01-02 15:04:05")))
			}
			fmt.Println("execute refresh")
		}
	}()

	// 布局
	content := container.NewVBox(
		widget.NewLabel("Click the button to manually update your DDNS record."),
		updateButton,
		statusLabel,
		updateTimeLabel,
	)
	ducWindow.SetContent(content)

	// 显示窗口
	ducWindow.Resize(fyne.NewSize(300, 150))
	ducWindow.ShowAndRun()
}
