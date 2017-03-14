package ktimer
import (
    "net/http"
    "fmt"
    "os"
    "os/signal"
    "time"
    "context"
    "syscall"
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
        
        //注册http请求的处理方法
        http.HandleFunc("/", WebHandler)
        srv := &http.Server{
            Addr: portdesc,
            Handler: http.DefaultServeMux,
            ReadTimeout: 10 * time.Second,
            WriteTimeout: 10 * time.Second,
            MaxHeaderBytes: 1 << 20,
        }
        slg,_ := GetSerLoger()

        go func(){
            //启动http服务
            slg.Println("web server starting...")
            if err := srv.ListenAndServe(); err!=nil {
                ServiceError("web server start listen fail.", err)
            }
        }()

        //监听系统信号
        stopChan := make(chan os.Signal)
        signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
        <-stopChan
        msg = "shutting down web server..."
        slg.Println(msg)
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        srv.Shutdown(ctx)
        msg = "web server gracefully stopped."
        slg.Println(msg)
    }

    //os.Exit(0)
}


//定义http请求的处理方法
func WebHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World, %v\n", time.Now())

}
