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

	mpris.OnMetadataChange = func(metadata gui.MusicMetadata) {
		fmt.Printf("Now playing: %s — %s\n", metadata.Artist, metadata.Title)
		ui.ShowPopup(metadata, 3*time.Second)
	}

	// wait for exit signal
	<-make(chan struct{})
}
