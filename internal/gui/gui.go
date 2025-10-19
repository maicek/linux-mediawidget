package gui

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MusicMetadata struct {
	Title  string
	Artist string
}

type GUI struct {
	visible          bool
	app              *gtk.Application
	win              *gtk.ApplicationWindow
	renderUntil      time.Time
	musicMetadata    MusicMetadata
	metadataElements map[string]*gtk.Label
}

func New() *GUI {
	app := gtk.NewApplication("dev.maicek.popupd", 0)

	gui := &GUI{
		app:              app,
		metadataElements: make(map[string]*gtk.Label),
	}

	app.Connect("activate", func() {
		win := gtk.NewApplicationWindow(app)
		win.SetTitle("Popupd")
		win.SetDefaultSize(400, 80)
		win.SetResizable(false)
		gui.win = win

		bg := gtk.NewBox(gtk.OrientationVertical, 6)
		bg.SetMarginStart(6)
		bg.SetMarginEnd(6)
		bg.SetMarginTop(0)
		bg.SetMarginBottom(0)
		bg.SetCanFocus(false)

		label := gtk.NewLabel("Title")
		bg.Append(label)

		gui.metadataElements["title"] = label

		label = gtk.NewLabel("Artist")
		bg.Append(label)
		gui.metadataElements["artist"] = label

		win.SetChild(bg)

		if gtk4layershell.IsSupported() {
			fmt.Println("Layer shell supported")
			gtk4layershell.InitForWindow(&win.Window)
			gtk4layershell.SetLayer(&win.Window, gtk4layershell.LayerShellLayerTop)
			gtk4layershell.SetExclusiveZone(&win.Window, 0)
			gtk4layershell.SetAnchor(&win.Window, gtk4layershell.LayerShellEdgeTop, true)
			gtk4layershell.SetMargin(&win.Window, gtk4layershell.LayerShellEdgeTop, -80)

			win.SetVisible(true)

		} else {
			fmt.Println("Layer shell not supported")
		}
	})

	targetMargin := -80
	currentMargin := -80

	go func() {
		for {
			if targetMargin != currentMargin {
				if targetMargin > currentMargin {
					currentMargin++
				} else {
					currentMargin--
				}
				gtk4layershell.SetMargin(&gui.win.Window, gtk4layershell.LayerShellEdgeTop, currentMargin)
			}

			time.Sleep(1 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if gui.win != nil {
				if gui.visible && time.Now().After(gui.renderUntil) {
					gui.visible = false
					targetMargin = -80

					fmt.Println("Hiding window")
				} else if !gui.visible && time.Now().Before(gui.renderUntil) {
					gui.visible = true
					targetMargin = 20

					fmt.Println("Showing window")
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	return gui
}

func (g *GUI) UpdateMusicMetadata() {
	g.metadataElements["title"].SetText(g.musicMetadata.Title)
	g.metadataElements["artist"].SetText(g.musicMetadata.Artist)
}

func (g *GUI) Run() {
	g.app.Run(nil)
}

func (g *GUI) ShowPopup(title string, timeout time.Duration) {
	g.renderUntil = time.Now().Add(timeout + 100*time.Millisecond)

	g.musicMetadata = MusicMetadata{
		Title:  title,
		Artist: "",
	}
	g.UpdateMusicMetadata()

	g.win.SetTitle(title)

}
