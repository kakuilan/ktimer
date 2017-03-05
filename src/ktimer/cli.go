package ktimer

import (
	//"flag"
	"fmt"
	"os"
	"strings"
)

var Commands = []string{
	"init",
	"start",
	"stop",
	"restart",
	"status",
	"version",
	"help",
	"count",
	"clear",
	"get",
	"del",
	"add",
}

//打印帮助信息
func Help() {
	fmt.Println(PRODESC)
	fmt.Printf("Version %s [%s]\n", VERSION, PUBDATE)
	fmt.Println("Author ", AUTHOR)
	fmt.Println("Usage:")
	fmt.Printf("%8s%s\n", " ", "ktimer command [arguments]")
	fmt.Println("The commands are:")
	fmt.Printf("%8s%-10s%-s\n", " ", "help", "show help information and usage")
	//fmt.Printf("%8s%-10s%-s\n", " ", "init", "initialization must be used for the first time")
	fmt.Printf("%8s%-10s%-s\n", " ", "start", "start service")
	fmt.Printf("%8s%-10s%-s\n", " ", "stop", "stop service")
	fmt.Printf("%8s%-10s%-s\n", " ", "restart", "restart service")
	fmt.Printf("%8s%-10s%-s\n", " ", "status", "show service status")
	fmt.Printf("%8s%-10s%-s\n", " ", "version", "show software version")
	fmt.Printf("%8s%-10s%-s\n", " ", "add", "add a timer,it has following parameters:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "-type: specified type, timer or ticker ")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "-time: specify how many seconds to execute, or time stamp")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "-limit: limit execution times. 0 is not limited to ticker")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "-command: specific operations to be performed")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "example:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer add -type=timer -time=1 -limit=1 -command=\"echo -e Hello Ktimer\"")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer add -type=ticker -time=1 -limit=0 -command=\"date --rfc-3339=ns\"")
	fmt.Printf("%8s%-10s%-s\n", " ", "get", "get the timer information by a key.The key is MD5 when inserted timer return.")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "example:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer get 912ec803b2ce49e4a541068d495ab570")
	fmt.Printf("%8s%-10s%-s\n", " ", "del", "delete the timer by a key.The key is MD5 when inserted timer return")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "example:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer del 912ec803b2ce49e4a541068d495ab570")
	fmt.Printf("%8s%-10s%-s\n", " ", "count", "show total number of current tasks")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "example:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer count")
	fmt.Printf("%8s%-10s%-s\n", " ", "clear", "clear current all tasks")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "example:")
	fmt.Printf("%8s%-10s%-s\n", " ", "", "ktimer clear")
	//fmt.Printf("%8s%-10s%-s\n"," ", "", "")
	os.Exit(0)
}

func commandErr(command string) {
	fmt.Printf("The command error,please see help: [ktimer -help]\n")
	os.Exit(0)
}

//捕获CLI命令参数
func CatchCli() {
	//获取命令行参数
	argNum := len(os.Args)

	//无参数直接显示使用方式
	if argNum == 1 {
		Help()
	} else {
		action := os.Args[1]
		action = strings.ToLower(action)
		if action == "help" || action == "-h" || action == "--h" || action == "-help" || action == "--help" {
			Help()
		}

		//检查是否存在该命令
		var isCommand bool = false
		for _, ac := range Commands {
			if ac == action {
				isCommand = true
				break
			}
		}
		if !isCommand {
			commandErr(action)
		}

		switch action {
		case "init":
			ServiceInit()
		case "start":
			ServiceStart()
		case "stop":
			ServiceStop()
		case "restart":
			ServiceRestart()
		case "status":
			ServiceStatus()
		case "version":
			ServiceVersion()
		}

		for j, arg := range os.Args {
			fmt.Printf("arg[%d] = %s \n", j, arg)
		}

	}

}
