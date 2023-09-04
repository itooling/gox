package aes

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
)

var (
	origin = "hello world"
	aesKey = "0123456789ABCDEF"
	gobKey = []byte("0123456789ABCDEF")
)

type User struct {
	Name string
	Data map[string]interface{}
}

func TestAesKeyGenerate(t *testing.T) {
	key := GenAesKey([]byte(aesKey))
	fmt.Println(string(key))
}

func TestAesKeyGenerateRand(t *testing.T) {
	key := GenAesKeyRand()
	fmt.Println(string(key))
}

func TestAes(t *testing.T) {
	key := GenAesKey([]byte(aesKey))
	en, _ := EncryptCBC([]byte(origin), key)
	fmt.Println(base64.StdEncoding.EncodeToString(en))
	de, _ := DecryptCBC(en, key)
	fmt.Println(string(de))
}

func TestWriteObj(t *testing.T) {
	u := User{Name: "test", Data: map[string]interface{}{
		"code":    200,
		"result":  "success",
		"content": "hello world",
	}}
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(&u)
	enc, _ := EncryptCBC(buf.Bytes(), gobKey)
	fo, _ := os.Create("out")
	gob.NewEncoder(fo).Encode(enc)
}

func TestReadObj(t *testing.T) {
	fi, _ := os.Open("out")
	buf := make([]byte, 0)
	gob.NewDecoder(fi).Decode(&buf)
	dec, _ := DecryptCBC(buf, gobKey)
	data := bytes.NewBuffer(dec)
	u := new(User)
	gob.NewDecoder(data).Decode(u)
	fmt.Println(u)
}
