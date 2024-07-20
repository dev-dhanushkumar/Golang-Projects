package filecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)


func  Encrypt(source string, password []byte) {
  if _, err := os.Stat(source); os.IsNotExist(err) {
    panic(err.Error())
  }

  srcFile, err := os.Open(source)
  if err != nil {
    panic(err.Error())
  }

  defer srcFile.Close()

  plainText, err := io.ReadAll(srcFile)
  if err != nil {
    panic(err.Error())
  }

  key := password

  nonce := make([]byte, 12) //[0,0,0,0,0,0,0,0,0,0,0,0]
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    panic(err.Error())
  }

  dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New())
  
  block, err := aes.NewCipher(dk)
  if err != nil {
    panic(err.Error())
  }


  // GCM Function
  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    panic(err.Error())
  }
}

func Decrypt() {

}
