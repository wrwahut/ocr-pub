package server

import (
	"fmt"
	"net"
	"bufio"
	"errors"
)

type socketServer struct{
	connSlice  []*net.TCPConn
	socketPort     int32
	data  chan string
}

func NewSocketServer(channeldata  chan string, socketPort int32) *socketServer{
	server := socketServer{
		socketPort : socketPort,
	}
	server.data = channeldata
	return &server
}

func (server *socketServer) listen(){
	tcpAdd,err:= net.ResolveTCPAddr("tcp",fmt.Sprintf(":%d", server.socketPort))  //解析tcp服务
	if err!=nil{
		fmt.Println("net.ResolveTCPAddr error:",err)
		return
	}
	tcpListener,err:=net.ListenTCP("tcp",tcpAdd)   //监听指定TCP服务
	if err!=nil{
		fmt.Println("net.ListenTCP error:",err)
		return
	}
	defer tcpListener.Close()
	for{
		tcpConn,err:=tcpListener.AcceptTCP() //阻塞，当有客户端连接时，才会运行下面
		if err!=nil{
			fmt.Println("tcpListener error :",err)
			continue
		}
		fmt.Println("A client connected:",tcpConn.RemoteAddr().String())
		// server.boradcastMessage(tcpConn.RemoteAddr().String()+"进入房间"+"\n")  //当有一个客户端进来之时，广播某某进入房间
		server.connSlice = append(server.connSlice,tcpConn)
		// 监听到被访问时，开一个协程处理
		go server.tcpPipe(tcpConn)
	}
}

func (server *socketServer) Run(){
	go server.listen()
	for {
		select {
		case i1 := <-server.data:
		   fmt.Println(i1)
		   server.boradcastMessage(i1)
		default:
		}
	 }    
  
}

func (server *socketServer) tcpPipe(conn *net.TCPConn){
	ipStr := conn.RemoteAddr().String()
	fmt.Println("ipStr:",ipStr)
	defer func(){
		fmt.Println("disconnected:",ipStr)
		conn.Close()
		server.deleteConn(conn)
		// server.boradcastMessage(ipStr+"离开了房间"+"\n")
	}()
	reader:=bufio.NewReader(conn)
	for{
		message,err:=reader.ReadString('\n')  //读取直到输入中第一次发生 ‘\n’
		//因为按强制退出的时候，他就先发送换行，然后在结束
		if message == "\n"{
			return
		}
		message = ipStr+"说："+message
		if err!=nil{
			fmt.Println("topPipe:",err)
			return
		}
		// 广播消息
		// fmt.Println(ipStr,"说：",message)
		// err = server.boradcastMessage(message)
		// if err!=nil{
		// 	fmt.Println(err)
		// 	return 
		// }
	}
}

// 移除已经关闭的客户端
func (server *socketServer) deleteConn(conn *net.TCPConn)error{
	if conn==nil{
		fmt.Println("conn is nil")
		return errors.New("conn is nil")
	}
	for i:= 0;i<len(server.connSlice);i++{
		if(server.connSlice[i]==conn){
			server.connSlice = append(server.connSlice[:i],server.connSlice[i+1:]...)
			break
		}
	}
	return nil
}

// 广播数据
func (server *socketServer) boradcastMessage(message string)error{
	b := []byte(message)
	for i:=0;i<len(server.connSlice);i++{
		// fmt.Println(server.connSlice[i])
		_,err := server.connSlice[i].Write(b)
		if err!=nil{
			fmt.Println("发送给",server.connSlice[i].RemoteAddr().String(),"数据失败"+err.Error())
			continue
		}
	}
	return nil
}
