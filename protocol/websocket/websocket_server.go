package ws

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thinkonmay/signaling-server/protocol"
)

var wsserver = WebSocketServer{}

type WebSocketServer struct {
	fun protocol.OnTenantFunc
}

func (server *WebSocketServer) OnTenant(fun protocol.OnTenantFunc) {
	server.fun = fun
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

func handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	params := strings.Split(r.URL.RawQuery, "=")
	if len(params) == 1 {
		return
	}

	tenant := NewWsTenant(c)
	err = wsserver.fun(params[1], tenant)
	if err != nil {
		tenant.Exit()
	}

	for {
		if tenant.IsExited() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func InitSignallingWs(conf *protocol.SignalingConfig) *WebSocketServer {
	http.HandleFunc("/api/handshake", handle)
	if conf.KeyFile != "" && conf.CertFile != "" {
		go http.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", conf.WebsocketPort),conf.CertFile,conf.KeyFile, nil)
	} else {
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", conf.WebsocketPort), nil)
	}
	return &wsserver
}
