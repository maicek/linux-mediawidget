package main

import (
	"fmt"
	"time"

	"maicek.dev/linux-mediawidget/internal/gui"
	"maicek.dev/linux-mediawidget/internal/mpris"
)

func main() {
	fmt.Println("Starting popupd daemon...")

	ui := gui.New()

	go ui.Run()

	mpris := mpris.New()
	mpris.Start()

	mpris.OnMetadataChange = func(title string, artist string) {
		fmt.Printf("Now playing: %s — %s\n", artist, title)
		ui.ShowPopup(fmt.Sprintf("Now playing: %s — %s", artist, title), 2*time.Second)
	}

	// wait for exit signal
	<-make(chan struct{})
}
