package tcpserver

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "9999", "port")
var buffer [1024]byte

type Msg struct {
	Data string `json:"data"`
	Type int    `json:"type"`
}

type Resp struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
}

func Tcpstart(ports string) {
	// 解析参数
	flag.Parse()
	var l net.Listener
	var err error
	// 监听
	l, err = net.Listen("tcp", *host+":"+ports)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on " + *host + ":" + *port)

	for {
		// 接收一个client
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}

		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		// 执行
		go handleRequest(conn)
	}
}

// 处理接收到的connection
//
func handleRequest(conn net.Conn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("Disconnected :" + ipStr)
		conn.Close()
	}()

	// 构建reader和writer
	//reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)

	for {
		// 读取一行数据, 以"\n"结尾
		n, err := conn.Read(buffer[0:])
		if err != nil {
			return
		}
		if n > 0 {
			buffer[n] = 0
			fmt.Println(buffer[0:n])
			conn.Write(buffer[0:n])
		}
		//conn.Write(r)
		//conn.Write([]byte("\n"))
	}

	fmt.Println("Done!")
}
