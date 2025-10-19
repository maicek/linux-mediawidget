package gui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"net/http"
)

var WindowSize = 70
var WindowWidth = 320

type MusicMetadata struct {
	Title      string
	Artist     string
	Album      string
	AlbumCover string
}

type MetadataRenderElements struct {
	title      *gtk.Label
	artist     *gtk.Label
	albumCover *gtk.Image
	// album  *gtk.Label
}

type GUI struct {
	visible          bool
	app              *gtk.Application
	win              *gtk.ApplicationWindow
	renderUntil      time.Time
	musicMetadata    MusicMetadata
	metadataElements MetadataRenderElements
}

func (gui *GUI) getCoverPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".popupd/album-cover.png")
}

func (gui *GUI) createMetadataElements(bg *gtk.Box) {
	fixed := gtk.NewFixed()
	fixed.SetCanFocus(false)
	fixed.SetCSSClasses([]string{"MediaInfo"})
	bg.Append(fixed)

	title := *gtk.NewLabel("The whole music title")
	title.SetCSSClasses([]string{"MediaInfo__title"})
	fixed.Put(&title, 0, 0)
	gui.metadataElements.title = &title

	artist := *gtk.NewLabel("The whole artist")
	artist.SetCSSClasses([]string{"MediaInfo__artist"})
	fixed.Put(&artist, 0, 20)
	gui.metadataElements.artist = &artist

	albumCover := *gtk.NewImage()
	albumCover.SetCSSClasses([]string{"MediaInfo__album-cover"})
	albumCover.SetPixelSize(50)
	albumCover.SetFromFile(gui.getCoverPath())

	fixed.Put(&albumCover, float64(WindowWidth-40), 0)

	gui.metadataElements.albumCover = &albumCover
}

func (gui *GUI) DownloadAlbumCover() error {
	if gui.musicMetadata.AlbumCover == "" {
		return nil
	}

	resp, err := http.Get(gui.musicMetadata.AlbumCover)
	if err != nil {
		fmt.Println("Failed to download album cover:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read album cover:", err)
		return err
	}

	// Save to file at ~/.popupd/album-cover.png
	os.WriteFile(gui.getCoverPath(), body, 0644)
	return nil
}

func (gui *GUI) updateAlbumCover() {
	err := gui.DownloadAlbumCover()
	if err != nil {
		fmt.Println("Failed to update album cover:", err)
		gui.metadataElements.albumCover.SetVisible(false)
		return
	}
	gui.metadataElements.albumCover.SetVisible(true)
	gui.metadataElements.albumCover.SetFromFile(gui.getCoverPath())
}

func New() *GUI {
	app := gtk.NewApplication("dev.maicek.popupd", 0)

	gui := &GUI{
		app: app,
	}

	home, _ := os.UserHomeDir()
	cssPath := filepath.Join(home, ".popupd/style.css")

	app.Connect("activate", func() {
		win := gtk.NewApplicationWindow(app)
		win.SetTitle("Popupd")
		win.SetDefaultSize(320, WindowSize)
		win.SetDecorated(false)
		win.SetCSSClasses([]string{"popupd-bg"})
		win.SetOpacity(0.8)
		win.SetResizable(false)
		gui.win = win

		display := gdk.DisplayGetDefault()
		provider := gtk.NewCSSProvider()
		provider.LoadFromPath(cssPath)
		gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_USER)

		bg := gtk.NewBox(gtk.OrientationVertical, 6)
		bg.SetCanFocus(false)

		gui.createMetadataElements(bg)

		win.SetChild(bg)

		if gtk4layershell.IsSupported() {
			fmt.Println("Layer shell supported")
			gtk4layershell.InitForWindow(&win.Window)
			gtk4layershell.SetLayer(&win.Window, gtk4layershell.LayerShellLayerTop)
			gtk4layershell.SetExclusiveZone(&win.Window, 0)
			gtk4layershell.SetAnchor(&win.Window, gtk4layershell.LayerShellEdgeTop, true)
			gtk4layershell.SetMargin(&win.Window, gtk4layershell.LayerShellEdgeTop, -WindowSize-20)

			win.SetVisible(true)

		} else {
			fmt.Println("Layer shell not supported")
			panic("Layer shell not supported")
		}
	})

	targetMargin := -WindowSize - 20
	currentMargin := -WindowSize - 20

	go func() {
		for {
			if targetMargin != currentMargin {
				if targetMargin > currentMargin {
					currentMargin++
				} else {
					currentMargin--
				}
				gtk4layershell.SetMargin(&gui.win.Window, gtk4layershell.LayerShellEdgeTop, int(currentMargin))
			}

			time.Sleep(1 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if gui.win != nil {
				if gui.visible && time.Now().After(gui.renderUntil) {
					gui.visible = false
					targetMargin = -WindowSize - 20

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

func (g *GUI) Run() {
	g.app.Run(nil)
}

func (g *GUI) UpdatePopupMusicMetadata(metadata MusicMetadata) {
	g.metadataElements.title.SetText(metadata.Title)
	g.metadataElements.artist.SetText(metadata.Artist)
	g.musicMetadata.AlbumCover = metadata.AlbumCover
	go g.updateAlbumCover()
}

func (g *GUI) ShowPopup(timeout time.Duration) {
	g.renderUntil = time.Now().Add(timeout + 100*time.Millisecond)
}

func (g *GUI) HidePopup() {
	g.renderUntil = time.Now()
}
