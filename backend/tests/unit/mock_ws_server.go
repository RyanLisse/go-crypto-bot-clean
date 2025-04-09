package unit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

// MockWsServer represents a mock WebSocket server for testing
type MockWsServer struct {
	Server   *httptest.Server
	URL      string
	Upgrader websocket.Upgrader
	Conn     *websocket.Conn // The connection to the client
	T        *testing.T
}

// NewMockWsServer creates a new mock WebSocket server for testing
func NewMockWsServer(t *testing.T) *MockWsServer {
	mock := &MockWsServer{
		T: t,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Create a test server
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := mock.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
			return
		}

		mock.Conn = conn

		// Simple echo server for testing
		go func() {
			defer conn.Close()
			for {
				messageType, p, err := conn.ReadMessage()
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") ||
						strings.Contains(err.Error(), "websocket: close") {
						return // Normal close
					}
					t.Logf("Read error: %v", err)
					return
				}

				// Log the received message
				t.Logf("Mock server received: %s", string(p))

				// Handle ping messages
				if strings.Contains(string(p), "ping") {
					// Send pong response
					pongMsg := strings.Replace(string(p), "ping", "pong", 1)
					err = conn.WriteMessage(messageType, []byte(pongMsg))
					if err != nil {
						t.Logf("Write error: %v", err)
						return
					}
				} else if strings.Contains(string(p), "SUBSCRIPTION") {
					// Handle subscription messages
					// Extract the channel from the message
					var subMsg map[string]interface{}
					if err := json.Unmarshal(p, &subMsg); err == nil {
						if params, ok := subMsg["params"].([]interface{}); ok && len(params) > 0 {
							channel := params[0].(string)
							// Send subscription confirmation
							confirmMsg := fmt.Sprintf(`{"e": "sub.success", "c": "%s"}`, channel)
							err = conn.WriteMessage(websocket.TextMessage, []byte(confirmMsg))
							if err != nil {
								t.Logf("Write error: %v", err)
								return
							}
						}
					}
				} else {
					// Echo other messages back
					err = conn.WriteMessage(messageType, p)
					if err != nil {
						t.Logf("Write error: %v", err)
						return
					}
				}
			}
		}()
	}))

	// Parse URL and replace scheme from http to ws
	u, err := url.Parse(mock.Server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %v", err)
	}
	u.Scheme = "ws"
	mock.URL = u.String()

	return mock
}

// Close closes the mock server
func (mock *MockWsServer) Close() {
	if mock.Conn != nil {
		mock.Conn.Close()
	}
	mock.Server.Close()
}

// WriteToClient sends a message to the connected client
func (mock *MockWsServer) WriteToClient(messageType int, data []byte) error {
	if mock.Conn == nil {
		mock.T.Logf("No client connection to write to")
		return nil
	}
	return mock.Conn.WriteMessage(messageType, data)
}
