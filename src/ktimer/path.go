package ktimer
import (
    "log"
    "os"
    "path/filepath"
    "strings"
)

//字符串截取
func Substr(s string, pos,length int) string {
    runes := []rune(s)
    l := pos + length
    if l > len(runes) {
        l = len(runes)
    }
    return string(runes[pos:l])
}

//获取父级目录
func GetParentDirectory(dirctory string) string {
    return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

//获取当前目录
func GetCurrentDirectory() string {
    dir,err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil{
        log.Fatal(err)
    }
    return strings.Replace(dir, "\\", "/", -1)
}

//检查文件是否存在
func FileExist(filename string) bool {
    _,err := os.Stat(filename)
    return err==nil || os.IsExist(err)
}
