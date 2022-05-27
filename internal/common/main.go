package common

import (
	"encoding/json"
	"fmt"
)

type StartGame struct {
	card int
	team int
}

func Response(r interface{}) []byte {
	a, _ := json.Marshal(r)
	fmt.Println(string(a)) // 20192
	return a
}