package display

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"training.go/Grypt/crypt"
)

const sizeByteSlice = 16

func stringSliceToByteSlice(b []string) ([]byte, error) {
	var bytes []byte
	for _, v := range b {
		u, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, byte(u))
	}
	return bytes, nil
}

func intSliceToString(numbers []int) string {
	var strNumbers []string
	for _, number := range numbers {
		strNumbers = append(strNumbers, strconv.Itoa(number))
	}
	return strings.Join(strNumbers, ", ")
}
func byteSliceToIntSlice(b []byte) []int {
	var ints []int
	for _, v := range b {
		ints = append(ints, int(v))
	}
	return ints
}

func generateBytes() ([]byte, error) {
	b := make([]byte, sizeByteSlice)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	for i := range b {
		b[i] %= 100
	}
	return b, nil
}

func generateSecret() (secret string, err error) {
	validChars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
	randomBytes := make([]byte, 24)

	_, err = rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for i, b := range randomBytes {
		randomBytes[i] = validChars[b%byte(len(validChars))]
	}

	secret = string(randomBytes)
	return
}

func generateNewPwd(pwd string) (secret string, b string, encryptPwd string, err error) {
	byt, err := generateBytes()
	if err != nil {
		return
	}
	secret, err = generateSecret()
	if err != nil {
		return
	}
	encryptPwd, err = crypt.Encrypt(pwd, secret, byt)
	if err != nil {
		return
	}
	b = intSliceToString(byteSliceToIntSlice(byt))
	return
}
