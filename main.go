package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 配置加载-初始化
	configTest := &Config{}
	// 保存默认配置
	saveConfig(configTest)
	// 创建app
	duc := app.New()
	ducWindow := duc.NewWindow("Ali-DDNS Client")

	// 更新状态标签
	statusLabel := widget.NewLabel("Ready")

	// 更新按钮点击事件
	updateButton := widget.NewButton("Update IP", func() {
		_err := refresh()
		if _err != nil {
			statusLabel.SetText(fmt.Sprintf("Error: %v", _err))
		} else {
			// setLastIP(currentIP, "last_ip.txt")
			statusLabel.SetText("IP updated successfully.")
		}
	})

	currentIP, _ := getPublicIP()
	statusLabel.SetText(fmt.Sprintf("Current IP: %s", currentIP))

	// 定时任务：每5分钟检查一次IP变化
	go func() {
		// for range time.Tick(5 * time.Minute) {
		// 	currentIP, _ := getPublicIP()
		// 	lastIP, _ := getLastIP("last_ip.txt")
		// 	if lastIP != currentIP {
		// 		updateDDNS(config, currentIP)
		// 		setLastIP(currentIP, "last_ip.txt")
		// 	}
		// }
	}()

	// 布局
	content := container.NewVBox(
		widget.NewLabel("Click the button to manually update your DDNS record."),
		updateButton,
		statusLabel,
	)
	ducWindow.SetContent(content)

	// 显示窗口
	ducWindow.Resize(fyne.NewSize(300, 150))
	ducWindow.ShowAndRun()
}
