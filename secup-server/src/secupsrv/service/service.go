package service
import (
    "fmt"
    "secupsrv/cryptor"
    "strings"
    "github.com/nu7hatch/gouuid"
    "secupsrv/dao"
    "os"
    "bufio"
)
var PASSWORD string = ""
type Data struct {
    Value, Token, Pass1, Pass2 string
}

type Request struct {
    Action string
    Data Data
    Signature string
}
func _Verify(signature string) (bool) {
  file, err := os.Open("./.sign")
  if err != nil {
     fmt.Println("Signature file not found")
     return false
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
    break
  }
  if lines[0] == signature {
    return true
  }
  return false
}
func Process(request* Request, password string) {
     if !_Verify(request.Signature) {
        return
     }
     PASSWORD = password
     action := request.Action
     fmt.Println(action)
     switch action {
        case "encode": Encode(&request.Data)
        case "decode": Decode(&request.Data)
        case "delete": Delete(&request.Data)
     }
}
func Encode(data* Data){
    value := data.Value
    passphrase1 := cryptor.RandString(16)
    passphrase2 := cryptor.RandString(16)
    data.Pass1 = passphrase1
    data.Pass2 = passphrase2
    newKey := strings.Join([]string{data.Pass1, data.Pass2}, "")
    byteKey := []byte(newKey)
    byteValue := []byte(value)
    ciphertext := cryptor.Encrypt(byteKey, byteValue);
    token, err := uuid.NewV4()
    if err != nil {
        fmt.Println(err);
    }
    data.Token = token.String()
    data.Pass1 = cryptor.EncodeBase64(cryptor.Encrypt([]byte(PASSWORD), []byte(data.Pass1)))
    data.Pass2 = cryptor.EncodeBase64(cryptor.Encrypt([]byte(PASSWORD), []byte(data.Pass2)))
    model := dao.Model{cryptor.EncodeBase64(ciphertext), data.Token, data.Pass1}
    fmt.Println(model)
    res, err := dao.Insert(model)
    if err != nil || res == false {
        fmt.Println("Could not save token into db")
    }
    data.Value = ""
    data.Pass1 = ""
}
func Decode(data* Data) {
    fmt.Println("yahahahah");
    token := data.Token
    model := dao.Select(token)
    passpharse1 := cryptor.Decrypt([]byte(PASSWORD), cryptor.DecodeBase64(model.Pass1))
    passphrase2 := cryptor.Decrypt([]byte(PASSWORD), cryptor.DecodeBase64(data.Pass2))
    newKey := strings.Join([]string{string(passpharse1), string(passphrase2)}, "")
    byteKey := []byte(newKey)
    value := cryptor.Decrypt(byteKey, cryptor.DecodeBase64(model.Data))
    data.Value = value
}
func Delete(data *Data) {
    if dao.Delete(data.Token) {
        data = nil
    } 
}
