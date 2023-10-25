package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {

	var url = "https://jsonplaceholder.typicode.com/todos/1"
	var resp, err = http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	var bodyString = new(strings.Builder)
	io.Copy(bodyString, resp.Body)

	fmt.Println(bodyString.String())

}
