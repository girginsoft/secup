package main

import (
    "fmt"
    "net"
    "io/ioutil"
    "bufio"
    "encoding/gob"
    "os"
    "secupsrv/service"
    "secupsrv/dao"
    "github.com/howeyc/gopass"
    "os/exec"
    "secupsrv/cryptor"
)
var PASSWORD string = ""
func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &service.Request{}
    dec.Decode(p)
    service.Process(p, PASSWORD)
    encoder := gob.NewEncoder(conn)
    encoder.Encode(&p.Data)
    fmt.Printf("Received : %+v", p);
}

func main() {
   param := os.Args[1]
   options := []string{"run", "install"}
   switch param {
        case "run": run()
        case "install": install()
        case "sign": sign()
        default: fmt.Println("Allowed options", options)
   }
}
func run() {
   fmt.Println("start");
   fmt.Printf("Password: ")
   PASSWORD = string(gopass.GetPasswd())
   if len(PASSWORD) != 16 {
        fmt.Println("Password must be 16 character")
        return
   }
   c := exec.Command("clear")
   c.Stdout = os.Stdout
   c.Run()
   fmt.Println("Running...")
   ln, err := net.Listen("tcp", ":8080")
   if err != nil {
        // handle error
   }
   for {
       conn, err := ln.Accept() 
       if err != nil {
           // handle error
           continue
       }
       go handleConnection(conn) 
   }

}
func install() {
    dao.Create()
}
func verifySign() {
}
func sign() {
    sign := cryptor.RandString(128)
    signature := []byte(sign)
    err := ioutil.WriteFile("./.sign", signature, 0644)
    if err != nil {
        fmt.Println("File error could not write to file")
        return
    }
    fmt.Println("Your sign is:")
    fmt.Println(sign)
}
