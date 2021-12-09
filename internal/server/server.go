package server

type channelData struct{
	billet_num string
	billet_type int8
	warpage_value int
	warpage_direction int
	billet_pic_path string
	warpage_pic_path string
}



type Server struct{
	redis_server    *redisServer
	socket_server   *socketServer
	web_server      *WebSocketServer
	stopch          chan struct{}
	channeldata     chan string
	socketdata      chan string
	webdata         chan string
}

func NewServer(redisHost string, redisDB int8,socketPort,webSocketPort int32) (*Server, error){
	server := Server{
		channeldata : make(chan string),
		socketdata : make(chan string),
		webdata : make(chan string),
	}
	redis_server, err:= NewRedisServer(redisHost, 1, server.channeldata)
	if err != nil{
		return nil, err
	}
	server.redis_server = redis_server
	socket_server := NewSocketServer(server.socketdata, socketPort)
	server.socket_server = socket_server

	web_server := NewWebSocketServer(server.webdata, webSocketPort)
	server.web_server = web_server
	return &server, nil
}

func (server *Server) scribeRedis(){
	for{
		select{
		case data := <- server.channeldata:
			server.socketdata <- data
			server.webdata <- data
		}
	}
}

func (server *Server) Start(){
	go server.scribeRedis()
	go server.redis_server.StartScribe()
	go server.socket_server.Run()
	go server.web_server.RunListen()
	<-server.stopch
}



