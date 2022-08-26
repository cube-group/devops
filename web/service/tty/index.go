package tty

import (
	"app/library/log"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/url"
)

var upgrader = websocket.Upgrader{}

//gotty for kubectl logs 更换Arguments Json内容
//字符串模式：{"Arguments":"?pid=626&pod=demo-7b874bb986-k8sb4","AuthToken":""}
//func argumentsExchange(assist string, b []byte) (res []byte) {
//	res = b
//	if strings.Index(string(b), `{"Arguments":`) == -1 {
//		return
//	}
//	var data map[string]interface{}
//	if err := json.Unmarshal(b, &data); err != nil {
//		return
//	}
//	u, err := url.Parse(convert.MustString(data["Arguments"]))
//	if err != nil {
//		return
//	}
//	var query = u.Query()
//	var respStr string
//	switch assist {
//	case "local": //local bash
//		return
//	case "history": //history id
//		if i := models.GetHistory(convert.MustUint32(query.Get("id"))); i != nil {
//			var logFilePath = i.WorkspaceFollowLog()
//			if i.IsEnd() {
//				if endLogPath := i.WorkspaceEndLog(); endLogPath == "" {
//					return
//				} else {
//					logFilePath = endLogPath
//				}
//			}
//			respStr = fmt.Sprintf(
//				`{"Arguments":"?arg=-c&arg=%s","AuthToken":"%s"}`,
//				url.QueryEscape(fmt.Sprintf(`tail -n 1000 -f %s`, logFilePath)),
//				setting.SysGoTtyRandBasicAuth,
//			)
//		}
//	case "node": //node
//		if node := models.GetNode(convert.MustUint32(query.Get("id"))); node != nil {
//			respStr = fmt.Sprintf(
//				`{"Arguments":"?arg=-p&arg=%s&arg=ssh&arg=-o&arg=StrictHostKeyChecking=no&arg=-p&arg=%s&arg=%s@%s","AuthToken":"%s"}`,
//				node.SshPassword, node.SshPort, node.SshUsername, node.IP,
//				setting.SysGoTtyRandBasicAuth,
//			)
//		}
//	case "pod": //docker pod
//		podType := query.Get("type")
//		h := models.GetHistory(convert.MustUint32(query.Get("id")))
//		if h == nil {
//			return
//		}
//		if podType == "log" {
//			respStr = fmt.Sprintf(
//				`{"Arguments":"?arg=-p&arg=%s&arg=ssh&arg=-o&arg=StrictHostKeyChecking=no&arg=-t&arg=-p&arg=%s&arg=%s@%s&arg=docker&arg=logs&arg=-f&arg=-n&arg=1000&arg=%s","AuthToken":"%s"}`,
//				h.Node.SshPassword, h.Node.SshPort, h.Node.SshUsername, h.Node.IP,
//				h.Project.Name,
//				setting.SysGoTtyRandBasicAuth,
//			)
//		} else if podType == "tty" {
//			respStr = fmt.Sprintf(
//				`{"Arguments":"?arg=-p&arg=%s&arg=ssh&arg=-o&arg=StrictHostKeyChecking=no&arg=-t&arg=-p&arg=%s&arg=%s@%s&arg=docker&arg=exec&arg=-it&arg=%s&arg=sh","AuthToken":"%s"}`,
//				h.Node.SshPassword, h.Node.SshPort, h.Node.SshUsername, h.Node.IP,
//				h.Project.Name,
//				setting.SysGoTtyRandBasicAuth,
//			)
//		} else {
//			return
//		}
//	case "port":
//		return
//	default: //kubectl
//	}
//	fmt.Println("tty", respStr)
//	res = []byte(respStr)
//	return
//}

//gotty proxy
func Proxy(c *gin.Context, port, assist string) {
	defer func() {
		//close tty client & stream
		models.TTYCacheClear(convert.MustUint32(port))
	}()

	//create websocket
	upgrader.Subprotocols = append(upgrader.Subprotocols, "webtty")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.StdWarning("gotty", "ws.upgrader", err)
		return
	}
	defer ws.Close()

	wsUrl := url.URL{Scheme: "ws", Host: "127.0.0.1:" + port, Path: "/ws"}
	websocket.DefaultDialer.Subprotocols = append(websocket.DefaultDialer.Subprotocols, "webtty")
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
		//if err := conn.WriteMessage(websocket.TextMessage, argumentsExchange(assist, bytes)); err != nil {
		//	break
		//}
		if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
			break
		}
	}
}
