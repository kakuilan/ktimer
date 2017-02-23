package ktimer
import (
    "fmt"
    "os"
    "flag"
)

//打印帮助信息
func Help() {
    fmt.Println("Ktimer is a simple timer/ticker manager by golang.")
    fmt.Println("Version ", VERSION)
    fmt.Println("Author ", AUTHOR)
    fmt.Println("Usage:")
    fmt.Printf("%8s%s\n", " ", "ktimer command [arguments]") 
    fmt.Println("The commands are:")
    fmt.Printf("%8s%-10s%-s\n"," ", "help", "show help information and usage")
    fmt.Printf("%8s%-10s%-s\n"," ", "init", "initialization must be used for the first time")
    fmt.Printf("%8s%-10s%-s\n"," ", "start", "start service")
    fmt.Printf("%8s%-10s%-s\n"," ", "stop", "stop service")
    fmt.Printf("%8s%-10s%-s\n"," ", "restart", "restart service")
    fmt.Printf("%8s%-10s%-s\n"," ", "status", "show service status")
    fmt.Printf("%8s%-10s%-s\n"," ", "version", "show software version")
    fmt.Printf("%8s%-10s%-s\n"," ", "add", "add a timer,it has following parameters:")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "-type: specified type, timer or ticker ")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "-time: specify how many seconds to execute, or time stamp")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "-limit: limit execution times. 0 is not limited to ticker")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "-command: specific operations to be performed")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "e.g.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer add -type=timer -time=1 -limit=1 -command=\"echo -e Hello Ktimer\"")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer add -type=ticker -time=1 -limit=0 -command=\"date --rfc-3339=ns\"")
    fmt.Printf("%8s%-10s%-s\n"," ", "get", "get the timer information by a key.The key is MD5 when inserted timer return.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "e.g.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer get 912ec803b2ce49e4a541068d495ab570")
    fmt.Printf("%8s%-10s%-s\n"," ", "del", "delete the timer by a key.The key is MD5 when inserted timer return")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "e.g.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer del 912ec803b2ce49e4a541068d495ab570")
    fmt.Printf("%8s%-10s%-s\n"," ", "count", "show total number of current tasks")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "e.g.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer count")
    fmt.Printf("%8s%-10s%-s\n"," ", "clear", "clear all current tasks")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "e.g.")
    fmt.Printf("%8s%-10s%-s\n"," ", "", "ktimer clear")
    //fmt.Printf("%8s%-10s%-s\n"," ", "", "")

}

//捕获CLI命令参数
func CatchCli() {
    action := os.Args[1]
    tType := flag.String("type", "timer", "Timer type")
    tTime := flag.Int("time", 1, "seconds or timestamp")
    tLimit := flag.Int("limit", 0, "limit number")
    tCommand := flag.String("command", "", "specific operation")

    fmt.Println("acton=", action, *tType, *tTime, *tLimit, *tCommand)

}
