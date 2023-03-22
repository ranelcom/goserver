package rest

import (
	"fmt"
	"net/http"
	"strings"
)

func Post(body string) {

	url := "http://localhost:59986/api/v2/resource/GPS-POC/json"
	method := "POST"

	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("response Status:", res.Status)
	fmt.Println(string(body))
}
