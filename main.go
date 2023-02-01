package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("fetch url error:%v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println("error status code:%v", res.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	fmt.Println("body:", string(body))
}
