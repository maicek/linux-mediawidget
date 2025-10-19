package mpris

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"maicek.dev/linux-mediawidget/internal/gui"
)

type Mpris struct {
	conn             *dbus.Conn
	OnMetadataChange func(metadata gui.MusicMetadata)
}

func New() *Mpris {
	return &Mpris{}
}

func (m *Mpris) Start() *Mpris {
	fmt.Println("Starting MPRIS daemon...")

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	// Do not close the connection here; we need it alive for the lifetime of the app
	// defer conn.Close()

	rule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)
	if call.Err != nil {
		panic(call.Err)
	}

	ch := make(chan *dbus.Signal, 10)
	conn.Signal(ch)

	fmt.Println("MPRIS daemon started.")

	go func() {
		for sig := range ch {
			if len(sig.Body) < 2 {
				continue
			}

			if iface, ok := sig.Body[0].(string); ok && iface == "org.mpris.MediaPlayer2.Player" {
				changed, ok := sig.Body[1].(map[string]dbus.Variant)
				if !ok {
					continue
				}

				if md, ok := changed["Metadata"]; ok {
					metadata := md.Value().(map[string]dbus.Variant)

					title := "Unknown"
					if t, ok := metadata["xesam:title"].Value().(string); ok {
						title = t
					}

					artist := "Unknown"
					if a, ok := metadata["xesam:artist"].Value().([]string); ok && len(a) > 0 {
						artist = a[0]
					}

					album := "Unknown"
					if al, ok := metadata["xesam:album"].Value().(string); ok {
						album = al
					}

					albumCover := ""
					if ac, ok := metadata["mpris:artUrl"].Value().(string); ok {
						albumCover = ac
					}

					m.OnMetadataChange(gui.MusicMetadata{
						Title:      title,
						Artist:     artist,
						Album:      album,
						AlbumCover: albumCover,
					})
				}
			}
		}
	}()

	m.conn = conn

	return m
}
