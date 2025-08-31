# Go Collaborative Text Editor

This project is a simple, terminal-based collaborative text editor built with Go. It allows one user (the client) to type and edit text, with all changes instantly reflected on a second user's terminal (the server) in real-time. The application uses a TUI (Terminal User Interface) for a clean editing experience directly in the console.

### Features
* **Real-time Synchronization**: Edits made by the client are sent to the server and displayed immediately.
* **Terminal-Based UI**: A shared text editor environment for both users, built with the `tview` library.
* **Precise Positional Editing**: The system transmits not just keystrokes, but also cursor positions, allowing for accurate insertions and deletions anywhere in the document.
* **Client-Server Architecture**: A clear distinction between the host (server) and the joiner (client).

## Setup and Installation

Follow these steps to get the application running on your local machine.

#### Prerequisites
* Go (version 1.18 or later) installed on your system.

#### Installation
1.  Clone or download the project files into a directory.
2.  Open your terminal and navigate into the project directory.
3.  Install the necessary dependencies by running:
    ```bash
    go mod tidy
    ```

## How to Run

The application operates in two modes: `server` and `client`.

#### 1. Start the Server (User A)
The first user must start the application in `server` mode. This will start a listener and wait for a client to connect.

* **To run the server on the default port (8080):**
    ```bash
    go run . server
    ```
* **To run the server on a custom port (e.g., 9999):**
    ```bash
    go run . server 9999
    ```
    The server terminal will then wait. The editor UI for the server will only launch *after* a client successfully connects.

#### 2. Connect the Client (User B)
The second user connects to the running server by starting the application in `client` mode and providing the server's address.

* **To connect to a server running on `localhost:8080`:**
    ```bash
    go run . client
    ```
* **To connect to a server at a specific IP address and port:**
    ```bash
    go run . client <server_ip_address>:<port>
    
    # Example:
    go run . client 192.168.1.5:9999
    ```
Once the client connects, the TUI editor will launch on both users' terminals simultaneously. The client can start typing, and the changes will appear on the server's screen.

## How It Works

The primary ingenuity lies in abstracting raw keystrokes into structured JSON "edit events," which encode the precise operation, character, and cursor position for every change. This event-driven model enables an asymmetrical architecture where the joiner (B) broadcasts local edits to the host (A), who applies them in a thread-safe manner. The application cleverly determines a "host" or "joiner" role from startup arguments, facilitating a peer-to-peer connection without requiring a separate dedicated server. This design transforms a simple byte stream into a robust editing session that correctly synchronizes insertions, deletions, and newlines anywhere in the document.
