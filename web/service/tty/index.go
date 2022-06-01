package tty

import (
	"app/library/log"
	"app/library/types/convert"
	"app/models"
	"app/setting"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"strings"
)

var upgrader = websocket.Upgrader{}

//gotty for kubectl logs 更换Arguments Json内容
//字符串模式：{"Arguments":"?pid=626&pod=demo-7b874bb986-k8sb4","AuthToken":""}
func argumentsExchange(assist string, b []byte) (res []byte) {
	res = b
	if strings.Index(string(b), `{"Arguments":`) == -1 {
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}
	u, err := url.Parse(convert.MustString(data["Arguments"]))
	if err != nil {
		return
	}
	var respStr string
	switch assist {
	case "local": //local bash
		return
	case "sshpass": //sshpass
		if node := models.GetNode(convert.MustUint32(u.Query().Get("id"))); node != nil {
			respStr = fmt.Sprintf(
				`{"Arguments":"?arg=-p&arg=%s&arg=ssh&arg=-o&arg=StrictHostKeyChecking=no&arg=-p&arg=%s&arg=%s@%s","AuthToken":"%s"}`,
				node.SshPassword, node.SshPort, node.SshUsername, node.IP,
				setting.SysGoTtyRandBasicAuth,
			)
			fmt.Println(respStr)
		}
	default: //kubectl
	}
	return []byte(respStr)
}

//gotty proxy
func Proxy(writer http.ResponseWriter, req *http.Request, goTtyAddress, assist string) {
	defer log.StdOut("gotty", "conn.closed", req.RequestURI)
	upgrader.Subprotocols = append(upgrader.Subprotocols, "gotty")
	ws, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		log.StdWarning("gotty", "ws.upgrader", err)
		return
	}
	defer ws.Close()

	wsUrl := url.URL{Scheme: "ws", Host: goTtyAddress, Path: "/ws"}
	websocket.DefaultDialer.Subprotocols = append(websocket.DefaultDialer.Subprotocols, "gotty")
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		log.StdWarning("gotty", "conn.dialer", wsUrl.String(), err)
		return
	}
	defer conn.Close()

	connError := false
	//conn->ws
	go func() {
		defer func() {
			connError = true
		}()
		for {
			_, bytes, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if err := ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
				break
			}
		}
	}()

	//ws->conn
	for {
		if connError {
			break
		}
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, argumentsExchange(assist, bytes)); err != nil {
			break
		}
	}
}
