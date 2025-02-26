# 🚀 Echo WebSocket Chat Application

A real-time chat application built with Go's Echo framework and JavaScript WebSockets.

## ✨ Features

- 💬 Real-time messaging
- 👤 Username identification
- 🔄 Instant message broadcasting
- 🔔 System notifications (user join/leave)
- 🎨 Clean, responsive UI
- 🛡️ Secure WebSocket connection

## 🏗️ Technology Stack

- **Backend**:

  - Go 1.20+
  - Echo Web Framework v4
  - Gorilla WebSocket

- **Frontend**:
  - HTML5
  - CSS3
  - Vanilla JavaScript

## 📦 Installation

### Prerequisites

- Go 1.20 or higher
- Git

### Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/websocket-chat.git
cd websocket-chat
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Run the server**

```bash
go run main.go
```

4. **Access the application**

Open your browser and navigate to:

```
http://localhost:8080
```

## 🔧 Project Structure

```
websocket-chat/
├── main.go             # Server code with WebSocket implementation
├── go.mod              # Go module dependencies
├── public/             # Static files directory
│   └── index.html      # Chat UI and client-side WebSocket code
└── README.md           # This file
```

## 💻 Usage

1. Open the application in your browser
2. Enter your desired username
3. Start chatting in real-time!
4. Open multiple browser windows to simulate different users

## 🧩 How It Works

### Server Side

The Go server uses Echo framework for HTTP routing and serves the static files. It implements a WebSocket hub that:

- Manages active client connections
- Handles client registration and disconnection
- Broadcasts messages to all connected clients

### Client Side

The frontend establishes a WebSocket connection to the server and:

- Sends messages to the server when the user submits the form
- Listens for incoming messages and updates the UI
- Displays system notifications
- Manages the chat interface

## 🔍 Code Highlights

- Concurrent message handling with Goroutines
- Clean separation between client and hub logic
- Proper WebSocket connection lifecycle management
- Responsive UI with minimal dependencies

## 🚀 Future Enhancements

- 🔐 User authentication
- 💾 Message persistence
- 🏠 Multiple chat rooms
- 📱 Mobile-optimized experience
- 📁 File sharing capabilities
- ✏️ Typing indicators

## 📄 License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0) - see the LICENSE file for details.

This means:

- ✅ You are free to use, modify, and distribute this software
- ✅ You can use this software for commercial purposes
- ✅ You must disclose source code of your modifications
- ✅ You must license your modifications under the same license
- ✅ You must state changes made to the code

## 👏 Acknowledgements

- [Echo Framework](https://echo.labstack.com/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

---

Built with ❤️ using Go and JavaScript
