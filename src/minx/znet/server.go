package znet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"../utils"
	"../ziface"
)

// Server 服务器基础机构
type Server struct {
	Name      string         // 服务器名
	IPVersion string         // tcp4 or other
	IP        string         // IP地址
	Port      int            // 端口号
	Router    ziface.IRouter // 路由
}

// CallbackToClient 客户端的HandleAPI
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] CallbackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("Write Error")
		return errors.New("CallBackToClient Error")
	}
	return nil
}

// Start 开启服务器
func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP:%s, Port:%d, is starting.\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version:%s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	// 开一个服务去监听
	go func() {
		// 1.获取一个监听套接字
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		// 监听失败直接return
		if err != nil {
			fmt.Println("resolve tcp error:", err)
			return
		}

		// 2.开始监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err:", err)
		}

		// 监听成功
		fmt.Println("start Server", s.Name, " success listening")

		// 客户端Id
		var cid uint32

		// 3.循环监听新链接
		for {
			// 3.1接受一个套接字
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// 3.2 TODO 最大连接数
			// 3.3 TODO 新链接的业务请求

			// 生成一个新的客户端连接
			dealConn := NewConnection(conn, cid, s.Router)
			// id ++, TODO 自定义的连接生成方法
			cid++
			// 3.4 启动客户端业务
			go dealConn.Start()

		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server :", s.Name, " stop")

	// TODO 关闭后的清理
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	// TODO 启动服务时的各种初始化

	// 阻塞主线程, 不让 listener 退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 添加路由
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router success!")
}

// NewServer 返回新的服务器
func NewServer(name string) ziface.IServer {

	conf := utils.GlobalObject
	conf.Reload()

	s := &Server{
		Name:      conf.Name,
		IPVersion: "tcp4",
		IP:        conf.Host,
		Port:      conf.TCPPort,
		Router:    nil,
	}

	return s
}
