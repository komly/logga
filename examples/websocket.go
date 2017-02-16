package main

import (
	"fmt"
	"github.com/Komly/logga"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := logga.NewLogger(
		logga.WithFormatter(logga.JSONFormatter{}),
		logga.WithOutput(os.Stderr),
	)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if websocket.IsWebSocketUpgrade(r) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				logger.Fatalf("WebSocket upgrade error: %s", err)
			}

			defer conn.Close()
			for {
				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					logger.SetOption(logga.WithOutput(os.Stderr))
					logger.Errorf("NextWriter() error: %s", err)
					return
				}
				logger.SetOption(logga.WithOutput(w))
				logger.Debugf("Debug message: %d", 1)
				time.Sleep(1 * time.Second)
			}
		}
		fmt.Fprintf(w, `
		<h1>Log</h1>	
		<ul id="log"></ul>
		<script>
		var log = document.getElementById('log');
		var ws = new WebSocket('ws://' + location.host + '/ws');
		ws.onopen = function() {
			log.appendChild(document.createTextNode('Connected'));
		}
		ws.onmessage = function(e) {
			var row = JSON.parse(e.data);
			var el = document.createElement('LI');
			el.innerHTML = row.level + ' - ' + row.time + ' - <b>' + row.message + '</b>';
			log.appendChild(el);
		}
		</script>
		`)
	})
	logger.Debugf("Debug message: %d", 1)
	logger.Debugf("Info message: %d", 2)
	logger.Warningf("Warning message: %d", 3)
	logger.Errorf("Error message: %d", 4)

	http.ListenAndServe(":3000", nil)
}
