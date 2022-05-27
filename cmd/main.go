package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go run main.go card=sayı team=sayı
func main() {

	http.ListenAndServe(":8008", Router().InitRouter())
	fmt.Println(getParamsFromArgs(os.Args))

}

func getParamsFromArgs(args []string) (int,int) {
	if strings.Contains(args[2],"card="){
		args[1], args[2]= args[2], args[1]
	}

	card := strings.ReplaceAll(args[1], "card=", "")
	team :=strings.ReplaceAll(args[2], "team=", "")
	c,_ := strconv.Atoi(card)
	t,_ := strconv.Atoi(team)

	return c,t
}