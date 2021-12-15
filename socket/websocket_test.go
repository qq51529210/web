package socket

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func Test_websocket(t *testing.T) {
	var ser http.Server
	ser.Addr = "127.0.0.1:3390"
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Test WebSocket</title>
	<script type="text/javascript">
		function append(d) {
			var output = document.getElementById('output')
			var p = document.createElement('p')
			p.innerHTML += d
			output.appendChild(p)
		}
		var ws = new WebSocket("ws://127.0.0.1:3390/ws")
		ws.onopen = function(e) {
			append('open')
		}
		ws.onmessage = function(e){
			append(e.data)
		}
		ws.onclose = function(e){
			append('close')
		}
		ws.onerror = function(e){
			append(e.data)
		}

		function sendMessage() {
			ws.send(document.getElementById('input').value)
		}

		function clearOuput() {
			var output = document.getElementById('output')
			var childs = output.childNodes
			for (var i = childs.length - 1; i >= 0; i--) { 
				output.removeChild(childs[i])
			}
		}

	</script>
</head>
<body>
	<div>
		<input type="text" id="input" />
		<input type="button" value="send" onclick="sendMessage()" />
		<button type="button" onclick="clearOuput()">clear</button>
	</div>
	<div id="output"></div>
</body>
</html>
`)
	})
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		payload := 1024
		conn, err := Accept(rw, r)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		err = conn.Write(CodeText, []byte("hello client."), payload)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		err = conn.ReadLoop(payload, func(c Code, b []byte) error {
			switch c {
			case CodePing:
				return conn.Write(CodePong, nil, payload)
			case CodeClose:
				buf.Reset()
				buf.WriteString("You close!")
				conn.Write(CodeClose, buf.Bytes(), payload)
				conn.Close()
				ser.Close()
				return nil
			case CodePong:
				return nil
			default:
				buf.Reset()
				if string(b) == "close" {
					err = conn.Write(CodeClose, nil, payload)
					ser.Close()
					return err
				}
				if len(b) == 0 {
					buf.WriteString("Empry message.")
					return conn.Write(c, buf.Bytes(), payload)
				}
				buf.WriteString("Your message is: ")
				buf.Write(b)
				return conn.Write(c, buf.Bytes(), payload)
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	ser.ListenAndServe()
}
