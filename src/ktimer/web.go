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

    pwd := RandPwd(12,3)
    pwdstr := string(pwd)
    fmt.Println(pwd, pwdstr)
    os.Exit(0)
    //open,err := CnfObj.Int("web::web.enable")


    


}

func WebHandler() {
    
}
