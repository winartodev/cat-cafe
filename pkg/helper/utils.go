package helper

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/mail"
	"time"
)

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateRandNumber(base string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(9000) + 1000
	return fmt.Sprintf("%s%d", base, num)
}

func PrettyPrint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}
