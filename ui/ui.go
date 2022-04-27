package ui

// Spawns a minimal terminal user interface
import (
	"github.com/ralim/uprog/config"
	"github.com/rivo/tview"
)

type UI struct {
	UARTStreamView *SerialView
	ProgrammerView *ProgramView
	app            *tview.Application
	pages          *tview.Pages
	config         *config.Config
}

func NewUI(config *config.Config) *UI {
	return &UI{
		config: config,
	}
}
func (u *UI) RunUI() {
	u.UARTStreamView = NewSerialView(u)
	u.ProgrammerView = NewProgramView(u)
	u.pages = tview.NewPages()
	u.pages.AddPage("Serial",
		u.UARTStreamView.view,
		true, true)
	u.pages.AddPage("Programmer",
		u.ProgrammerView.view,
		true, true)
	u.app = tview.NewApplication()
	u.app.SetRoot(u.pages, true)
	u.app.SetFocus(u.pages)
	u.app.EnableMouse(true)
	u.ShowUARTLog()
	u.UARTStreamView.OpenPort()
	u.app.Run()
}

func (u *UI) ShowUARTLog() {
	if u == nil || u.pages == nil {
		return
	}
	u.pages.SwitchToPage("Serial")
	u.UARTStreamView.OpenPort()
}

func (u *UI) ShowProgrammerView() {
	if u == nil || u.pages == nil {
		return
	}
	u.pages.SwitchToPage("Programmer")
	u.ProgrammerView.Run()
}
