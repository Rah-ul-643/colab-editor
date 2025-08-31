package server

import (
	"bufio"
	"colab/common"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Server holds the state for the collaborative editor server.
type Server struct {
	port     int
	protocol string
	Listener net.Listener
	//clients [] *Client
}

// NewServer creates and returns a new server instance.
func NewServer(port int, protocol string) *Server {
	return &Server{
		port:     port,
		protocol: protocol,
	}
}

// ListenAndServe starts the server and begins listening for client connections.
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen(s.protocol, fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	s.Listener = listener
	fmt.Printf("Server is listening at (%s:%d)...\n", s.protocol, s.port)

	s.startAcceptConnectionsLoop()
	return nil
}

// startAcceptConnectionsLoop continuously accepts and handles new connections.
func (s *Server) startAcceptConnectionsLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			// If the listener is closed, we can exit the loop.
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Println("Error accepting connection:", err)
			}
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection manages a single client connection.
func (s *Server) handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("New connection from", remoteAddr)
	defer conn.Close()

	// Initialize the TUI editor only once when the first client connects.

	app := tview.NewApplication()
	textArea := tview.NewTextArea()
	textArea.SetBorder(true).SetTitle(" Collaborative Editor (Server View) ")

	// The server's text area is read-only to prevent local edits.
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Block all input by returning nil.
		return nil
	})

	// Run the TUI application in a separate goroutine.
	go func() {
		if err := app.SetRoot(textArea, true).Run(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(conn)
	for {
		// Read messages from the client, delimited by a newline.
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected.\n", remoteAddr)
				app.Stop()
				app = nil
				break
			}
			fmt.Println("Failed to read data:", err)
			return
		}

		var edit common.Edit
		if err := json.Unmarshal([]byte(message), &edit); err != nil {
			fmt.Println("Failed to unmarshal edit:", err)
			continue
		}

		// Queue an update to the UI to apply the received edit.
		// This is the thread-safe way to interact with tview components.
		app.QueueUpdateDraw(func() {
			s.applyEdit(edit, textArea)
		})
	}
}

// applyEdit modifies the server's text area based on the client's action.
func (s *Server) applyEdit(edit common.Edit, textArea *tview.TextArea) {
	// Get the current text content and split it into lines.
	text := textArea.GetText()
	lines := strings.Split(text, "\n")

	// Ensure the line index is within bounds.
	if edit.Row >= len(lines) {
		return
	}

	switch edit.Operation {
	case "insert":
		// Handle character insertion.
		if edit.Char == '\n' {
			// Handle newline (Enter key).
			line := lines[edit.Row]
			remaining := line[edit.Col:]
			lines[edit.Row] = line[:edit.Col]
			// Insert a new line into the slice.
			lines = append(lines[:edit.Row+1], append([]string{remaining}, lines[edit.Row+1:]...)...)
		} else {
			// Handle regular character insertion.
			line := lines[edit.Row]
			lines[edit.Row] = line[:edit.Col] + string(edit.Char) + line[edit.Col:]
		}

	case "delete":
		// Handle character deletion (Backspace).
		if edit.Col > 0 {
			// Delete character within the line.
			line := lines[edit.Row]
			lines[edit.Row] = line[:edit.Col-1] + line[edit.Col:]
		} else if edit.Row > 0 {
			// If at the start of a line, merge with the previous line.
			prevLine := lines[edit.Row-1]
			currentLine := lines[edit.Row]
			lines[edit.Row-1] = prevLine + currentLine
			// Remove the now-empty current line.
			lines = append(lines[:edit.Row], lines[edit.Row+1:]...)
		}
	}

	// Rejoin the lines and update the text area content.
	// The 'false' argument prevents the cursor from jumping to the end.
	textArea.SetText(strings.Join(lines, "\n"), false)
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	if s.Listener != nil {
		if err := s.Listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
		fmt.Println("Server closed successfully")
	}
	return nil
}
