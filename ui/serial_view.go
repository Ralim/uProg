package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jacobsa/go-serial/serial"
	"github.com/rivo/tview"
	"io"
	"sync"
)

// A Textview mixed serial port
// Uart -> Text view and text view -> uart
// Also listen for menu hotkey, programming hot-key

type SerialView struct {
	sync.Mutex
	view             *tview.TextView
	viewWriter       io.Writer
	serialOptions    *serial.OpenOptions
	menuChordPressed bool
	parent           *UI
	portStream       *io.ReadWriteCloser
	portCloseSignal  chan bool
}

func NewSerialView(parent *UI) *SerialView {
	options := serial.OpenOptions{
		PortName:                parent.config.SerialPort,
		BaudRate:                parent.config.BaudRate,
		DataBits:                8,
		StopBits:                1,
		ParityMode:              0,
		RTSCTSFlowControl:       false,
		InterCharacterTimeout:   50,
		MinimumReadSize:         4,
		Rs485Enable:             false,
		Rs485RtsHighDuringSend:  false,
		Rs485RtsHighAfterSend:   false,
		Rs485RxDuringTx:         false,
		Rs485DelayRtsBeforeSend: 0,
		Rs485DelayRtsAfterSend:  0,
	}
	UARTStreamView := tview.NewTextView()
	UARTStreamView.SetTextAlign(tview.AlignLeft)
	UARTStreamView.SetDynamicColors(true)
	UARTStreamView.SetChangedFunc(func() {
		parent.app.Draw()
	})
	UARTStreamView.SetWrap(false)
	UARTStreamView.SetScrollable(true)

	UARTStreamView.SetMaxLines(4096)
	UARTStreamView.SetBorder(true)
	w := tview.ANSIWriter(UARTStreamView)
	u := &SerialView{
		view:            UARTStreamView,
		serialOptions:   &options,
		portCloseSignal: make(chan bool),
		parent:          parent,
		viewWriter:      w,
	}
	UARTStreamView.SetInputCapture(u.handleUartLogKeyPress)
	defer u.updateTitle()

	return u
}

func (v *SerialView) handleUartLogKeyPress(event *tcell.EventKey) *tcell.EventKey {
	keyCode := event.Key()
	if ((event.Modifiers()&tcell.ModCtrl == tcell.ModCtrl) && event.Rune() == 'k') || keyCode == tcell.KeyCtrlK {
		// Chord start
		defer v.updateTitle()
		if !v.menuChordPressed {
			v.menuChordPressed = true
			return nil // suppress handling in the stack
		}
		v.menuChordPressed = false
	} else if v.menuChordPressed {
		defer v.updateTitle()
		switch event.Rune() {
		case 'o':
			v.OpenPort()
		case 'c':
			v.ClosePort()
		case 'p':
			// We want to trigger the programming command to be executed in its owm view
			// This is raised up the stack
			v.ClosePort()
			v.parent.ShowProgrammerView()
		case 'l':
			v.view.Clear()
		}
		v.menuChordPressed = false
		return nil
	}
	encoded := []byte(string(event.Rune()))
	//Special handlers
	if keyCode == tcell.KeyEnter {
		encoded = []byte("\r\n")
	}
	//Send to the uart
	if v.portStream != nil {
		_, _ = (*v.portStream).Write(encoded)
		if v.parent.config.LocalEcho {
			_, _ = v.view.Write(encoded)
		}
	}
	// Swallow the key press
	return nil

}
func (v *SerialView) ClosePort() {
	if v.portStream != nil {
		(*v.portStream).Close()
		v.portStream = nil
		_, _ = v.view.Write([]byte("Port Closed\r\n"))

	}
}

func (v *SerialView) OpenPort() {
	if v.portStream == nil {
		port, err := serial.Open(*v.serialOptions)
		if err == nil {
			v.portStream = &port
			go func(r io.Reader, w io.Writer) {
				_, _ = io.Copy(w, r)

			}(port, v.viewWriter)
			_, _ = v.view.Write([]byte("Port Opened\r\n"))
		} else {
			_, _ = v.view.Write([]byte(fmt.Sprintf("Could not open port ~> %v\r\n", err)))
		}
	}
}

func (v *SerialView) updateTitle() {
	chord := ""
	status := "Closed"
	if v.portStream != nil {
		status = "Open"
	}
	if v.menuChordPressed {
		chord = "| Chord pressed"
	}
	v.view.SetTitle(fmt.Sprintf("Uart Stream | Port: %s %s", status, chord))
}
