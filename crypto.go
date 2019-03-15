// Playbook - http://play.golang.org/p/3wFl4lacjX

package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}
	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

// Encrypt a string
func Encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := Pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("error", err, string(iv))
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return finalMsg, nil
}

// Decrypt a string
func Decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}

// // EncryptKey using key the SmsDbIO object
// func (sms *SmsDbIO) EncryptKey(key []byte) (err error) {

// 	if len(key) != 16 {
// 		return fmt.Errorf("error illegal key length required length 16")
// 	}

// 	if sms.Msg.Address, err = Encrypt(key, sms.Msg.Address); err != nil {
// 		return
// 	}
// 	if sms.Msg.Timestamp, err = Encrypt(key, sms.Msg.Timestamp); err != nil {
// 		return
// 	}

// 	return
// }

// // DecryptKey using key to an SmsDbIO object
// func (sms *SmsDbIO) DecryptKey(key []byte) (err error) {

// 	if len(key) != 16 {
// 		return fmt.Errorf("error illegal key length required length 16")
// 	}

// 	if sms.Msg.Address, err = Decrypt(key, sms.Msg.Address); err != nil {
// 		return
// 	}
// 	if sms.Msg.Timestamp, err = Decrypt(key, sms.Msg.Timestamp); err != nil {
// 		return
// 	}
// 	return err
// }

// Encrypt using key the SmsDbIO object
func (sms *SmsDbIO) Encrypt(key []byte) (err error) {

	if len(key) != 16 {
		return fmt.Errorf("error illegal key length required length 16")
	}

	// if sms.Msg.Address, err = Encrypt(key, sms.Msg.Address); err != nil {
	// 	return
	// }
	// if sms.Msg.Timestamp, err = Encrypt(key, sms.Msg.Timestamp); err != nil {
	// 	return
	// }

	if sms.Msg.ContactName, err = Encrypt(key, sms.Msg.ContactName); err != nil {
		return
	}
	if sms.Msg.Date, err = Encrypt(key, sms.Msg.Date); err != nil {
		return
	}
	if sms.Msg.Subject, err = Encrypt(key, sms.Msg.Subject); err != nil {
		return
	}
  var text string
	if text, err = Encrypt(key, string(sms.Msg.Body)); err != nil {
		return
	}
  sms.Msg.Body = Body(text)
	// if sms.Msg.Type, err = Encrypt(key, sms.Msg.Type); err != nil {
	// 	return err
	// }

	return
}

// Decrypt using key to an SmsDbIO object
func (sms *SmsDbIO) Decrypt(key []byte) (err error) {

	if len(key) != 16 {
		return fmt.Errorf("error illegal key length required length 16")
	}

	// if sms.Msg.Address, err = Decrypt(key, sms.Msg.Address); err != nil {
	// 	return
	// }
	// if sms.Msg.Timestamp, err = Decrypt(key, sms.Msg.Timestamp); err != nil {
	// 	return
	// }

	if sms.Msg.ContactName, err = Decrypt(key, sms.Msg.ContactName); err != nil {
		return
	}
	if sms.Msg.Date, err = Decrypt(key, sms.Msg.Date); err != nil {
		return
	}
	if sms.Msg.Subject, err = Decrypt(key, sms.Msg.Subject); err != nil {
		return
	}
  var text string
	if text, err = Decrypt(key, string(sms.Msg.Body)); err != nil {
		return
	}
  sms.Msg.Body = Body(text)
	// if sms.Msg.Type, err = Decrypt(key, sms.Msg.Type); err != nil {
	// 	return err
	// }

	return
}
