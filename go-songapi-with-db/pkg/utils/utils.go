package utils

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func ParseBody(b *http.Request, i interface{}){
	if body, err := ioutil.ReadAll(b.Body); err == nil{
		if err := json.Unmarshal([]byte(body), i); err != nil{
			return
		}
	}
}