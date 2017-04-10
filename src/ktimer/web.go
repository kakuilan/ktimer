package ktimer
import (
    "net/http"
    "fmt"
    "os"
    "os/signal"
    "time"
    "context"
    "syscall"
    "strings"
    "encoding/json"
)

//请求日志结构体
type ReqLog struct {
    Addr string `json:"addr"`
    Method string `json:"method"`
    Url string `json:"url"`
    Params interface{} `json:"params"`
    Header interface{} `json:"header"`
}

//输出结构体
type OutPut struct {
    Status bool `json:"status"`
    Code int `json:"code"`
    Data interface{} `json:"data"`
    Msg string `json:"msg"`
}

//timer接受参数结构体
type TimerParm struct {
    Action string `json:"action"`
    Type string `json:"type"`
    Time string `json:"time"`
    Limit string `json:"limit"`
    Command string `json:"command"`
    Kid string `json:"kid"`
    Passwd string `json:"passwd"`
    Starttime string `json:"starttime"`
    Endtime string `json:"endtime"`
}

//排除的参数字符
var ExcParChat = [...]string{
    "\\",
    "/",
    "\"",
    "'",
    " ",
    ":",
}

//允许的动作
var Actions = []string {
    "add",
    "get",
    "del",
    "update",
    "count",
    "list",
}


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
    }else if open<1{
        LogService("web server is setting close.", open)
    }else{
        port,err := CnfObj.Int("web::web.port")
        if err!= nil {
            ServiceError("web server get conf fail.", err)
        }
        bind_ip := CnfObj.String("web::web.bind_ip")
        passwd := CnfObj.String("web::web.passwd")
        portdesc := ":"+fmt.Sprint(port)
        LogService("web server listen to:", bind_ip, portdesc, passwd)

        //注册http请求的处理方法
        http.HandleFunc("/", WebHandler)
        srv := &http.Server{
            Addr: portdesc,
            Handler: http.DefaultServeMux,
            ReadTimeout: 10 * time.Second,
            WriteTimeout: 10 * time.Second,
            MaxHeaderBytes: 1 << 20,
        }

        go func(){
            //启动http服务
            LogWebes("web server starting...")
            if err := srv.ListenAndServe(); err!=nil {
                ServiceError("web server start listen fail.", err)
            }
        }()

        //监听系统信号
        stopChan := make(chan os.Signal)
        signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
        <-stopChan
        msg = "shutting down web server..."
        LogWebes(msg)
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        srv.Shutdown(ctx)
        msg = "web server gracefully stopped."
        LogWebes(msg)
    }

    //os.Exit(0)
}

//定义http请求的处理方法
func WebHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    var isAct bool
    LogWebes("accept request:", getRequestLog(r)) 

    //检查密码是否正确
    timPar := getTimerParams(r)
    CnfObj, err = GetConfObj()
    if err!=nil {
        LogWebes("web server accept request has err:", err)
        outputJson(w, false, 500, "web server has error.", "")
    }
    pwd := CnfObj.String("web::web.passwd")
    if(timPar.Passwd != pwd && false) {
        outputJson(w, false, 401, "You are not authorized to access", "")
    }else{
        for _,ac := range Actions {
            if ac == timPar.Action {
                isAct = true
                break
            }
        }
 
        if timPar.Action=="" || !isAct {
            outputJson(w, false, 200, "parameter act missing or error", "")
        }else{
            
        }
    }
}

//输出json
func outputJson(w http.ResponseWriter,status bool, code int, msg string, data interface{}) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if data==nil {
        data = time.Now()
    }
    res := OutPut{
        Status : status,
        Code : code,
        Data : data,
        Msg : msg,
    }
    jsonRes,err := json.Marshal(res)
    if err!=nil {
        fmt.Fprintln(w, "json err:", err)
    }else{
        fmt.Fprint(w, string(jsonRes))
    }
}

//获取完整url
func getFullUrl(r *http.Request) string {
    scheme := "http://"
    if r.TLS != nil {
        scheme = "https://"
    }

    return strings.Join([]string{scheme, r.Host, r.RequestURI}, "")
}

//获取头信息
func getHeader(r *http.Request) interface{} {
    m := make(map[string] interface{}) 
    for k,v := range r.Header {
        key := strings.ToLower(k)
        if key=="referer" || key== "user-agent" {
            m[k] = v
        }
    }

    return m
}

//获取请求参数
func getRequestParams(r *http.Request) interface{} {
    var pos int
    m := make(map[string] interface{}) 
    
    r.ParseForm()
    for k,v := range r.Form {
        for _,chr := range ExcParChat {
            pos = strings.Index(k, chr)
            if pos!=-1 {
                break
            }
        } 

        if pos!=-1 {
            continue
        }

        m[k] = v
    }

    return m
}

//获取定时器的参数
func getTimerParams(r *http.Request) TimerParm {
    var tp TimerParm
    r.ParseForm()

    if len(r.Form["action"])>0 {
        tp.Action  = r.Form["action"][0]
    }
    if len(r.Form["type"])>0 {
        tp.Type = r.Form["type"][0]
    }
    if len(r.Form["time"])>0 {
        tp.Time = r.Form["time"][0]
    }
    if len(r.Form["limit"])>0 {
        tp.Limit = r.Form["limit"][0]
    }
    if len(r.Form["command"])>0 {
        tp.Command = r.Form["command"][0]
    }
    if len(r.Form["kid"])>0 {
        tp.Kid = r.Form["kid"][0]
    }
    if len(r.Form["passwd"])>0 {
        tp.Passwd = r.Form["passwd"][0]
    }
    if len(r.Form["starttime"])>0 {
        tp.Starttime = r.Form["starttime"][0]
    }
    if len(r.Form["endtime"])>0 {
        tp.Endtime = r.Form["endtime"][0]
    }

    return tp
}


//获取访问日志
func getRequestLog(r *http.Request) ReqLog {
    url := getFullUrl(r)
    hea := getHeader(r)
    params := getRequestParams(r)
    jsonRes,err := json.Marshal(params)
    par := string(jsonRes)
    if err!= nil{
        par = "params json err"
    }

    log := ReqLog{
        r.RemoteAddr ,
        r.Method ,
        url,
        par,
        hea,
    }

    return log
}
