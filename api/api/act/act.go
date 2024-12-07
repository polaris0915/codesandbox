package act

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/polaris/codesandbox/api/request"
	"github.com/polaris/codesandbox/api/response"
	"github.com/polaris/codesandbox/logger"
	"net/http"
	"sync"
	"time"
)

var (
	CPP    = "c++"
	Golang = "go"
)

// 执行代码沙箱的锁
var mu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有域进行连接
		return true
	},
}

// 客户端连接
type WsConnection struct {
	Conn *websocket.Conn // 底层websocket
	//inChan  chan *request.WebSocketRequest  // 读队列
	OutChan chan response.WebSocketResponse // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	CloseChan chan byte // 关闭通知
}

// 关闭连接
func (wsConn *WsConnection) close() {
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed {
		return
	}
	wsConn.Conn.Close()
	wsConn.isClosed = true
	//close(wsConn.inChan)
	close(wsConn.OutChan)
	close(wsConn.CloseChan)
}

// 接收循环
func (wsConn *WsConnection) readLoop() {
	for {
		wsConn.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, message, err := wsConn.Conn.ReadMessage()
		if err != nil {
			fmt.Println(" wsConn.conn.ReadMessage错误: ", err)
			//wsConn.outChan <- response.NewSystemErrorResponse(wsConn.conn, response.SystemError)
			//wsConn.close()
			return
		}
		_request := new(request.WebSocketRequest)
		if err := json.Unmarshal(message, _request); err != nil {
			wsConn.OutChan <- response.NewSystemErrorResponse(response.ParamsError)
			continue
		}
		switch _request.Activity {
		case request.HeartBeatAct: // 在这里编写关闭连接逻辑
			fmt.Printf("%s检查心跳包....\n", wsConn.Conn.RemoteAddr().String())
			break
		case request.SubmitAct:
			fmt.Printf("%s用户提交问题....\n", wsConn.Conn.RemoteAddr().String())
			HandleSubmitAct(wsConn, message)
			break
		case request.RunAct:
			fmt.Printf("%s用户调试问题....\n", wsConn.Conn.RemoteAddr().String())
			HandleRunAct(wsConn, message)
			break
		default:
			wsConn.OutChan <- response.NewSystemErrorResponse(response.ParamsError)
		}
	}
}

// 响应循环
func (wsConn *WsConnection) writeLoop() {
	for {
		select {
		//case _response := <-wsConn.OutChan:
		//	_response.Response()
		case _response := <-wsConn.OutChan:
			if err := wsConn.Conn.WriteJSON(_response.Response()); err != nil {
				logger.GetLogger().Errorln(err)
			}
		// 处理主动关闭连接
		case <-wsConn.CloseChan:
			wsConn.close()
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	// 应答客户端告知升级连接为websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	wsConn := &WsConnection{
		Conn: conn,
		//inChan:    make(chan *request.WebSocketRequest, 10),
		OutChan:   make(chan response.WebSocketResponse, 10),
		CloseChan: make(chan byte),
		isClosed:  false,
	}

	// 读协程
	wg.Add(1)
	go wsConn.readLoop()
	// 写协程
	wg.Add(1)
	go wsConn.writeLoop()

	fmt.Printf("%s用户连接已经准备就绪了\n", wsConn.Conn.RemoteAddr().String())
	wg.Wait()
}

func GinWsHandler(c *gin.Context) {
	WsHandler(c.Writer, c.Request)
}
