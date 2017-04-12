package main

import (
	"ktimer"
    "runtime/pprof"
    "os"
)

func main() {
    f, _ := os.Create("profile_file")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

	//命令行处理
	ktimer.CatchCli()
}
