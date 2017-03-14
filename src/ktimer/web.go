package ktimer
import (
    "net"
    "fmt"
//    "os"
    "time"
)

//WEB容器
func WebContainer() {
    var err error
    var msg string
    fmt.Println("start web server...")
    CnfObj, err = GetConfObj()
    if err!= nil {
        ServiceError("web server get conf fail.", err)
    }

    open,err := CnfObj.Int("web::web.enable")
    if err!=nil {
        ServiceError("web server get conf fail.", err)
    }else{
        port,err := CnfObj.Int("web::web.port")
        if err!= nil {
            ServiceError("web server get conf fail.", err)
        }
        bind_ip := CnfObj.String("web::web.bind_ip")
        passwd := CnfObj.String("web::web.passwd")
        portdesc := ":"+fmt.Sprint(port)
        fmt.Println(333, open,port,bind_ip,passwd, portdesc)
        
        //开启监听
        listener,err := net.Listen("tcp", portdesc)
        if err!=nil {
            ServiceError("web server start listener fail.", err)
        }
        defer listener.Close()

        wlg,_ := GetWebLoger()
        for{
            //循环接收客户端的连接,没有连接时会阻塞,出错则跳出循环
            conn,err := listener.Accept()
            if err != nil {
                msg = "client accept has error."
                fmt.Println(msg, err)
                wlg.Println(msg, err)
                break
            }

            msg = "web server accept new connection."
            fmt.Println(msg)
            wlg.Println(msg)

            go WebHandler(conn)
        }
    }

    //os.Exit(0)
}

func WebHandler(conn net.Conn) {
   defer conn.Close() 
   for {
       //循环从连接中读取请求内容,没有请求时会阻塞,出错则跳出循环
        request := make([]byte, 128)
        readLength,err := conn.Read(request)

        if err != nil {
            fmt.Println(err)
            break
        }

        if readLength == 0{
            fmt.Println(err)
            break
        }

        //控制台输出读取到的请求内容，并在请求内容前加上hello和时间后向客户端输出
        fmt.Println("[server] request from ", string(request))
        conn.Write([]byte("hello " + string(request) + ", time: " + time.Now().Format("2006-01-02 15:04:05")))


   }


}
