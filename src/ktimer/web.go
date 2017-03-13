package ktimer
import (
    //"net"
    "fmt"
    "os"
)

//WEB容器
func WebContainer() {
    var err error
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

        fmt.Println(open,port,bind_ip,passwd)
    }


    os.Exit(0)
}

func WebHandler() {
    
}
