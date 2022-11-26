package utils

import (
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiZDQyYWMyMTUtY2U2ZC00MDkzLWFhYjUtNzM0ZGQyMmMzNjdiIiwiSUQiOjEsIlVzZXJuYW1lIjoiZ2VmZiIsIkJ1ZmZlclRpbWUiOjg2NDAwLCJpc3MiOiJnZWZmIiwiZXhwIjoxNjY5NTYyMTM0LCJuYmYiOjE2Njg5NTczMzR9.ur9tBZ5zxkpMVV5JT7cj2wkMa3HeWhr7L9-YarjkSaM"
	mj := &MyJwt{
		SigningKey: []byte("09b70634-1056-4891-8555-538e9d43573f"),
	}
	cl, err := mj.ParseToken(token)
	fmt.Println(cl, err)
}
