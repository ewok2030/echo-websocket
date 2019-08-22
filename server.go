// websockets.go
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// function to nicely format the HTTP request as string
func formatRequestInfo(r *http.Request) string {
	// Dump HTTP request
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	}
	return string(dump)
}

// function to handle HTTP endpoint
func handleHttpEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(200)
	fmt.Fprint(w, "## Server Information\n\n")

	// Hostname
	host, err := os.Hostname()
	if err == nil {
		fmt.Fprintf(w, "Hostname: %s\n", host)
	} else {
		fmt.Fprintf(w, "Hostname: %s\n", err.Error())
	}

	// Environment Variables
	allowShowEnvStr := os.Getenv("ALLOW_SHOW_ENV")
	allowShowEnv, err := strconv.ParseBool(allowShowEnvStr)
	requestShowEnv := r.URL.Query().Get("show_env")
	requestShowK8s := r.URL.Query().Get("show_k8s")

	if len(requestShowK8s) != 0 {
		fmt.Fprint(w, "\n## Kubernetes Environment\n\n")
		// Loop through K8S Variables
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "K8S") {
				fmt.Fprintf(w, "%s\n", e)
			}
		}
	}

	if err == nil && allowShowEnv && len(requestShowEnv) != 0 {
		fmt.Fprint(w, "\n## Server Environment\n\n")
		// Loop through Environment Variables
		for _, e := range os.Environ() {
			fmt.Fprintf(w, "%s\n", e)
		}
	}

	// Client Info
	fmt.Fprintf(w, "\n## Client Information:\n\n")
	fmt.Fprintf(w, "Client: %s\n", r.RemoteAddr)

	// Request Info
	fmt.Fprintf(w, "\n## Request Information:\n\n")
	fmt.Fprintf(w, "%s\n", formatRequestInfo(r))

}

// function to handle WebSocket endpoint
func handleWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	for {
		// Read message from client
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to client
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

// Main function to setup HTTP server
func main() {
	http.HandleFunc("/ws", handleWebsocketEndpoint)
	http.HandleFunc("/http", handleHttpEndpoint)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public")
	})

	http.ListenAndServe(":8081", nil)
}
