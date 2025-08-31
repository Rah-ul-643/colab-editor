package common

// Edit represents a change made in the text editor.
// It's serialized to JSON and sent from the client to the server.

type Edit struct {
	// Operation defines the type of action, e.g., "insert" or "delete".
	Operation string `json:"op"`
	// Char is the character being inserted. It's not used for deletion.
	Char rune `json:"char"`
	// Row is the cursor's row position for the edit.
	Row int `json:"row"`
	// Col is the cursor's column position for the edit.
	Col int `json:"col"`
}
