package ktimer
import (
    "fmt"    
)

func Help() {
    fmt.Println("Ktimer is a simple timer/ticker manager by golang.")
    fmt.Println("Version ", VERSION)
    fmt.Println("Author ", AUTHOR)
    fmt.Println("Usage:")
    fmt.Printf("%8s%s\n", " ", "ktimer command [arguments]") 
    fmt.Println("The commands are:")
    fmt.Printf("%8s%-10s%-100s|\n"," ", "foo", "b")



}
