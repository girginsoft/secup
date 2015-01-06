package cryptor
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

func EncodeBase64(b []byte) string {
    return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) []byte {
    data, err := base64.StdEncoding.DecodeString(s)
    if err != nil {
        panic(err)
    }
    return data
}

func Encrypt(key, text []byte) []byte {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    b := EncodeBase64(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext
}

func Decrypt(key, text []byte) string {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }
    if len(text) < aes.BlockSize {
        panic("ciphertext is too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    return string(DecodeBase64(string(text)))
}
func RandString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}
