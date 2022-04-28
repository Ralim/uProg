package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"os/exec"
	"sync"
)

type ProgramView struct {
	sync.Mutex
	view       *tview.TextView
	viewWriter io.Writer
	running    bool
	parent     *UI
}

func NewProgramView(parent *UI) *ProgramView {

	UARTStreamView := tview.NewTextView()
	UARTStreamView.SetTextAlign(tview.AlignLeft)
	UARTStreamView.SetDynamicColors(true)
	UARTStreamView.SetChangedFunc(func() {
		parent.app.Draw()
	})
	UARTStreamView.SetWrap(false)
	UARTStreamView.SetScrollable(true)

	UARTStreamView.SetTitle("Programming")
	UARTStreamView.SetBorder(true)
	w := tview.ANSIWriter(UARTStreamView)
	u := &ProgramView{
		view:       UARTStreamView,
		parent:     parent,
		viewWriter: w,
	}
	UARTStreamView.SetInputCapture(u.keyPress)
	defer u.updateTitle()

	return u
}

func (v *ProgramView) keyPress(event *tcell.EventKey) *tcell.EventKey {
	if v.running {
		return nil // block input during running command for now
	}
	if event.Rune() == 'r' {
		v.Run()
	} else {
		//On other keys we drop back to the uart view
		v.parent.ShowUARTLog()
	}
	// Swallow the key press
	return nil

}
func (v *ProgramView) Run() {
	// Spawn worker to run the programming process
	go v.runner()
}
func (v *ProgramView) runner() {
	v.running = true
	v.updateTitle()
	v.view.Clear()
	command := v.parent.config.ProgrammingCommand

	_, _ = v.viewWriter.Write([]byte(fmt.Sprintf("Running programming command %s\r\n", command)))

	// Spawn worker to run the programming process
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = v.viewWriter
	cmd.Stderr = v.viewWriter
	err := cmd.Run()
	v.running = false
	v.updateTitle()
	if err != nil {
		_, _ = v.viewWriter.Write([]byte(fmt.Sprintf("Programmer exited with error %v\r\nPress r to retry, or any other key to exit.\r\n", err)))
		//log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		_, _ = v.viewWriter.Write([]byte("Programmer finished\r\nPress r to re-run, or any other key to exit.\n"))
	}
}

func (v *ProgramView) updateTitle() {
	status := "Not Running"
	if v.running {
		status = "Running"
	}
	v.view.SetTitle(fmt.Sprintf(" Programming | %s ", status))
}
