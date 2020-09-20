package main

import (
    "fmt"
    "inhand/clawer"
    "inhand/parser"
)

func main() {
    
    uri := "http://www.cnblogs.com/"
    res := clawer.Claw(uri)
    newUris := parser.Parse(res.Body)
    fmt.Println(res.Body)
    fmt.Println(newUris)
}

