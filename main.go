package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/webkit"
	"github.com/therecipe/qt/widgets"
)

var (
	ap           *Application
	NavImage     *gui.QMovie
	zoomLevels   = []float64{30, 50, 67, 80, 90, 100, 110, 120, 133, 150, 170, 200, 240, 300}
	zoomState    = 5
	current_page *webkit.QWebPage
)

type Application struct {
	*widgets.QApplication
	WebView   *webkit.QWebView
	Window    *widgets.QMainWindow
	Statusbar *widgets.QStatusBar
	Urlbar    *widgets.QLineEdit
}

func main() {
	ap = &Application{}
	ap.QApplication = widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	ap.Window = window
	ap.Window.SetWindowTitle("Webkit")

	ap.Statusbar = widgets.NewQStatusBar(window)
	ap.Window.SetStatusBar(ap.Statusbar)

	wv := webkit.NewQWebView(window)
	ap.WebView = wv

	wv.ConnectLoadStarted(load_started)

	wv.ConnectLoadFinished(load_finished)
	wv.ConnectUrlChanged(url_changed)
	wv.ConnectWheelEvent(wheel_event)
	wv.ConnectTitleChanged(title_changed)

	vlayout := widgets.NewQVBoxLayout()
	hlayout := widgets.NewQHBoxLayout()

	vlayout.AddLayout(hlayout, 0)
	vlayout.AddWidget(wv, 0, 0)

	mw := widgets.NewQWidget(window, core.Qt__Widget)
	mw.SetLayout(vlayout)
	window.SetCentralWidget(mw)

	homebutton_click()

	widgets.QApplication_SetStyle2("fusion")
	ap.Window.ShowMaximized()
	widgets.QApplication_Exec()
}

func title_changed(title string) {
	ap.Window.SetWindowTitle(title)
}

func homebutton_click() {
	url := `https://parallelcoin.info/`
	qurl := core.NewQUrl3(url, core.QUrl__TolerantMode)
	ap.WebView.Load(qurl)
}

func stopbutton_click() {
	ap.WebView.Stop()
}

func urlbar_edited() {
	url := ap.Urlbar.Text()
	qurl := core.NewQUrl3(url, core.QUrl__TolerantMode)
	ap.WebView.Load(qurl)
}

func url_changed(url *core.QUrl) {
	ap.Urlbar.SetText(url.ToString(core.QUrl__None))
}

func load_started() {
	NavImage.Start()
}

func load_finished(ok bool) {
	NavImage.Stop()
	NavImage.JumpToFrame(0)

	current_page = ap.WebView.Page()
	current_page.ConnectLinkHovered(link_hovered)
}

func link_hovered(link string, title string, textContent string) {
	ap.Statusbar.ShowMessage(link, 0)
}

func backbutton_click() {
	ap.WebView.Back()
}

func forwardbutton_click() {
	ap.WebView.Forward()
}

func reloadbutton_click() {
	ap.WebView.Reload()
}

func wheel_event(e *gui.QWheelEvent) { // alternately: scale+=(event->delta()/120)
	if e.Modifiers() == core.Qt__ControlModifier {
		if e.AngleDelta().Y() > 0 {
			if zoomState != len(zoomLevels)-1 {
				zoomState += 1
			}
			ap.WebView.SetZoomFactor(zoomLevels[zoomState] / 100)
		} else {
			if zoomState != 0 {
				zoomState -= 1
			}
			ap.WebView.SetZoomFactor(zoomLevels[zoomState] / 100)
		}
	} else {
		ap.WebView.WheelEventDefault(e)
	}
}
