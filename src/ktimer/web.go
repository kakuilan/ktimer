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
    Key string `json:"key"`
    Passwd string `json:"passwd"`
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
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    wlg,_ := GetWebLoger()
    wlg.Println("accept a new request:", getRequestLog(r))
    wlg.Println("full log:", r)

    //json
    webRes := OutPut{
        Status : true,
        Data : time.Now(),
        Msg : "success",
    }
    
    jsonRes,err := json.Marshal(webRes)
    if err != nil {
        fmt.Fprintln(w, "json err:", err)
    }else{
        fmt.Fprint(w, string(jsonRes))
    }

    p := getRequestParams(r)
    p2,_ := json.Marshal(p)
    fmt.Fprint(w, string(p2))
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
func getTimerParams(r *http.Request) interface{} {
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
    if len(r.Form["key"])>0 {
        tp.Key = r.Form["key"][0]
    }
    if len(r.Form["passwd"])>0 {
        tp.Passwd = r.Form["passwd"][0]
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
