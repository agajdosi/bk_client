/*##### BEGIN GPL LICENSE BLOCK #####

  This program is free software; you can redistribute it and/or
  modify it under the terms of the GNU General Public License
  as published by the Free Software Foundation; either version 2
  of the License, or (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program; if not, write to the Free Software Foundation,
  Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

##### END GPL LICENSE BLOCK #####*/

package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"fyne.io/systray"
)

// traySupported reports that a system tray icon can be shown on this platform.
// On Windows the systray implementation is pure Go (no cgo), so it builds with
// the same CGO_ENABLED=0 cross-compile used for the shipped binaries.
const traySupported = true

//go:embed icons/blendkit.ico
var trayIcon []byte

// runTray shows the Blendkit-Client system tray icon and blocks until the user
// quits from its menu. It is expected to be called on the main goroutine while
// the HTTP server runs in a separate goroutine.
//
// serverURL is the Blendkit website the tray links to; listenAddr is the local
// address the Client is serving on (shown for reference).
func runTray(serverURL, listenAddr string) {
	onReady := func() {
		systray.SetIcon(trayIcon)
		systray.SetTitle("Blendkit-Client")
		systray.SetTooltip(fmt.Sprintf("Blendkit-Client v%s — %s", ClientVersion, listenAddr))

		mVersion := systray.AddMenuItem(fmt.Sprintf("Blendkit-Client v%s", ClientVersion), "")
		mVersion.Disable()
		mAddr := systray.AddMenuItem("Listening on "+listenAddr, "Local address the Client serves on")
		mAddr.Disable()

		systray.AddSeparator()
		mOpenWeb := systray.AddMenuItem("Open Blendkit.com", "Open the Blendkit website in your browser")
		mOpenDev := systray.AddMenuItem("Open Dev Dashboard", "Open the Client's developer dashboard in your browser")
		// The dev dashboard is only served when BLENDKIT_DEBUG=1; hide the
		// menu entry otherwise so it can't link to a 404.
		if !devDashboardEnabled() {
			mOpenDev.Hide()
		}
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Stop the Blendkit-Client")

		devURL := fmt.Sprintf("http://%s/dev", listenAddr)

		go func() {
			for {
				select {
				case <-mOpenWeb.ClickedCh:
					openInBrowser(serverURL)
				case <-mOpenDev.ClickedCh:
					openInBrowser(devURL)
				case <-mQuit.ClickedCh:
					systray.Quit()
					return
				}
			}
		}()
	}

	onExit := func() {
		BKLog.Printf("%s Tray closed, shutting down.", EmoOK)
		os.Exit(0)
	}

	systray.Run(onReady, onExit)
}

// openInBrowser opens the given URL in the default web browser on Windows.
func openInBrowser(url string) {
	if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start(); err != nil {
		BKLog.Printf("%s Could not open browser for %s: %v", EmoWarning, url, err)
	}
}
