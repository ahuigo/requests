package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
)

// Delete Form Request
func TestDeleteForm(t *testing.T) {
	println("Test DELETE method: form data(x-wwww-form-urlencoded)")
	data := requests.Datas{
		"comments": "ew",
	}
	session := requests.R() //.SetDebug()
	resp, err := session.Delete("https://www.httpbin.org/delete", data)
	if err == nil {
		fmt.Println(resp.Text())
	}
}
