# popupd

popupd is a simple daemon that shows a popup window with the current MPRIS metadata.

Development phase 1.

TODO:

- Improve GUI Code
- Better Handler (use channels instead of polling)
- Fix mpris events handling to prevent random popup show.

## Requires:
- gtk4
- gtk4-layer-shell

Run with

```
GDK_BACKEND=wayland go run main.go
```
