package rest

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
)

func Post(resource string, value string) {

  url := "http://localhost:59986/api/v2/resource/GPS-POC/" + resource
  method := "POST"

  payload := strings.NewReader(value)

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Content-Type", "text/plain")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}