package client

import (
	"colab/common"
	"encoding/json"
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StartEditorAndSync connects to the server and starts the collaborative TUI.
func StartEditorAndSync(address string) error {
	// Connect to the server.
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server. Starting editor...")

	app := tview.NewApplication()
	textArea := tview.NewTextArea()
	textArea.SetBorder(true).SetTitle(" Collaborative Editor (Client) ")

	// SetInputCapture intercepts all key events.
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get the current cursor position before the event is processed.
		row, col, _, _ := textArea.GetCursor()
		var edit common.Edit

		switch event.Key() {
		case tcell.KeyRune:
			// A standard character was typed.
			edit = common.Edit{
				Operation: "insert",
				Char:      event.Rune(),
				Row:       row,
				Col:       col,
			}
		case tcell.KeyEnter:
			// The Enter key was pressed.
			edit = common.Edit{
				Operation: "insert",
				Char:      '\n',
				Row:       row,
				Col:       col,
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			// The Backspace key was pressed.
			// Note: tcell differentiates between Backspace and Ctrl+H.
			edit = common.Edit{
				Operation: "delete",
				Row:       row,
				Col:       col,
			}
		case tcell.KeyCtrlC:
			// Gracefully exit on Ctrl+C.
			app.Stop()
			return nil

		default:
			// Let the text area handle other keys (like arrow keys) locally.
			return event
		}

		// Send the edit to the server.
		if err := sendEdit(conn, edit); err != nil {
			// In a real app, you might want to handle this more gracefully.
			fmt.Println("Error sending edit:", err)
			app.Stop()
		}

		// Return the event so the local text area is also updated.
		return event
	})

	// Run the TUI application.
	if err := app.SetRoot(textArea, true).Run(); err != nil {
		return fmt.Errorf("error running client application: %w", err)
	}

	return nil
}

// sendEdit serializes an Edit struct to JSON and sends it to the server.
func sendEdit(conn net.Conn, edit common.Edit) error {
	jsonData, err := json.Marshal(edit)
	if err != nil {
		return fmt.Errorf("failed to marshal edit: %w", err)
	}

	// Append a newline character to act as a message delimiter.
	_, err = conn.Write(append(jsonData, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write to connection: %w", err)
	}
	return nil
}
