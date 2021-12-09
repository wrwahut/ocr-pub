package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/gorilla/mux"
)

type WebSocketServer struct{
	data  chan string
	webSocketPort  int32
	pool map[*connection]bool
    b chan []byte
    r chan *connection
    u chan *connection
}

type connection struct {
    ws   *websocket.Conn
    sc   chan []byte
}

var wu = &websocket.Upgrader{ReadBufferSize: 512,
    WriteBufferSize: 512, CheckOrigin: func(r *http.Request) bool { return true }}

func NewWebSocketServer(channeldata  chan string,webSocketPort  int32) *WebSocketServer{
	server := WebSocketServer{
		pool: make(map[*connection]bool),
		u: make(chan *connection),
		b: make(chan []byte),
		r: make(chan *connection),
		data: channeldata,
		webSocketPort: webSocketPort,
	}
	// server.data = channeldata
	return &server
}

func (server *WebSocketServer) listen_channel(){
	for {
		select {
		case i1 := <-server.data:
			server.b <- []byte(i1)
		default:
		}
	 }    
}

func (server *WebSocketServer) RunListen(){
	go server.listen_channel()
	router := mux.NewRouter()
	go server.run()
	router.HandleFunc("/ws", server.myws)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", server.webSocketPort), router); err != nil {
        fmt.Println("err:", err)
    }
}

func (server *WebSocketServer) run(){
	for {
		select {
		case c := <- server.r:
			server.pool[c] = true
		case c := <- server.u:
			if _, ok := server.pool[c]; ok {
                delete(server.pool, c)
                close(c.sc)
            }
		case data := <-server.b:
			fmt.Println("wesocket send data->",string(data))
            for c := range server.pool {
                select {
                case c.sc <- data:
					
                default:
                    delete(server.pool, c)
                    close(c.sc)
                }
            } 
		}
	}
}

func (server *WebSocketServer) myws(w http.ResponseWriter, r *http.Request) {
    ws, err := wu.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    c := &connection{sc: make(chan []byte, 256), ws: ws}
    server.r <- c
    go c.writer()
    c.reader(server)
    defer func() {
        server.u <- c
    }()
}


func (conn *connection) writer() {
	for message := range conn.sc{
		conn.ws.WriteMessage(websocket.TextMessage, message)
	}
	conn.ws.Close()
}

func (conn *connection) reader(server *WebSocketServer) {
	for {
        _, _, err := conn.ws.ReadMessage()
		if err != nil {
            server.r <- conn
            break
        }
	}
}