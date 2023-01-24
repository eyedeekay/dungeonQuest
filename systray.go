package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/atotto/clipboard"
	goi2pbrowser "github.com/eyedeekay/go-i2pbrowser"
)

func startState(dir string) {
	if _, err := os.Stat(filepath.Join(dir, "running")); err == nil {
		os.Exit(0)
	}
}

func stopState(dir string) {
	os.RemoveAll(filepath.Join(dir, "running"))
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("mooz")
	systray.SetTooltip("Video calls over I2P")
	mBrowse := systray.AddMenuItem("Join game", "Open a game window")
	mCopy := systray.AddMenuItem("Copy Game URL", "Copy the URL of your game")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		for {
			select {
			case <-mBrowse.ClickedCh:
				k, _ := garlic.Keys()
				v := k.Address.Base32()
				addr := strings.TrimSpace(fmt.Sprintf("https://%s/game/client/index.html", v))
				h := *clientDir
				profileDir := filepath.Join(h, "dungeonCrawler")
				go goi2pbrowser.BrowseApp(profileDir, addr)
			case <-mCopy.ClickedCh:
				k, _ := garlic.Keys()
				v := k.Address.Base32()
				addr := strings.TrimSpace(fmt.Sprintf("https://%s/game/client/index.html", v))
				clipboard.WriteAll(addr)
			}
			time.Sleep(time.Second)
		}
	}()
	// Sets the icon of a menu item.
	mQuit.SetIcon(icon.Data)
}

func onExit() {
	// clean up here
	if runtime.GOOS == "windows" {
		syscall.Kill(syscall.Getpid(), syscall.SIGKILL)
	} else {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
}
