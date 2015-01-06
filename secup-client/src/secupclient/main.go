package main

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
    "os"
    "time"
    "crypto/md5"
)
var SIGN string = "0ZY3G2Kfc9aAvHTqGP96BrV6Mx3sosGTtVc9lqmrOO9NC0DYdg5grCYgpn37Q6ld0nPeml1cm9t2YiyoY0FcF3qUOGKSpaMhEhbbGwxqtLl7wvEXnEfat3UxG3rW7z0t"
type Data struct {
    Value, Token, Pass1, Pass2 string
}
type Request struct {
    Action string
    Data Data
    Signature string
}
func main() {
     if len(os.Args) < 3 {
        fmt.Println("Please specify action and parameter")
        fmt.Println("./secupclient encode value")
        fmt.Println("./secupclient decode token")
        fmt.Println("./secupclient delete token")
        return
     }
     action := os.Args[1]
     param := os.Args[2]
     switch action {
        case "encode": encode(param)
        case "decode": decode(param)
        case "delete": delete(param)
        default: fmt.Println("Please choose one of these operations encode, decode, delete")
    }
}
func _verify (token string, md5hash string) (bool) {
    t := time.Now()
    currentDate := t.Format("200601021504")
    hash := fmt.Sprint(token, currentDate)
    h := md5.New()
    calculatedHash := string(h.Sum([]byte(hash)))
    fmt.Println(md5hash)
    fmt.Println(calculatedHash)
    if md5hash == calculatedHash {
        return true
    }
    return false
}
func encode(value string) {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    //Request
    encoder := gob.NewEncoder(conn)
    data := &Data{value, "", "",""}
    request := &Request{"encode", *data, SIGN}
    encoder.Encode(request)
    //Response
    dec := gob.NewDecoder(conn)
    response := &Data{}
    dec.Decode(response)
    fmt.Println("Your token and password is:")
    fmt.Println(response.Token)
    fmt.Print(response.Pass2)
    conn.Close()
}
func decode(token string) {
    if len(os.Args) < 4 {
        fmt.Println("Your password is required!")
        return
    }
    md5 := os.Args[4]
    if !_verify(token, md5) {
        fmt.Println("Your request could not validated")
        return
    }
    password := os.Args[3]
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    //request
    encoder := gob.NewEncoder(conn)
    data := &Data{"", token, "", password}
    request := &Request{"decode", *data, SIGN}
    encoder.Encode(request)
    //Response
    dec := gob.NewDecoder(conn)
    data = &Data{}
    dec.Decode(data)
    fmt.Println("Your value is:");
    fmt.Println(data.Value)
    conn.Close()
}
func delete(token string) {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    md5 := os.Args[3]
    if !_verify(token, md5) {
        fmt.Println("Your request could not validated")
        return
    }
    //Request
    data := &Data{"", token, "", ""}
    encoder := gob.NewEncoder(conn)
    request := &Request{"delete", *data, SIGN}
    encoder.Encode(request)
    //Response
    dec := gob.NewDecoder(conn)
    data = &Data{}
    dec.Decode(data)
    fmt.Println("Deleted successfully");
    conn.Close()
}
