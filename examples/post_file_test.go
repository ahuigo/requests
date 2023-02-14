package examples

import (
	"os"
	"testing"

	"github.com/ahuigo/requests"
)

/*
An example about post both `file` and `form data`:
curl "https://www.httpbin.org/post" -F 'file1=@./go.mod' -F 'file2=@./version' -F 'name=alex'
*/
func TestPostFile(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	path, _ := os.Getwd()
	resp, err := requests.Post(
		ts.URL+"/file",
		requests.Files{
			"file1": path + "/go.mod",
			"file2": path + "/version",
		},
		requests.FormData{
			"name": "alex",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Files struct {
			File2 string
		}
		Form struct {
			Name string
		}
	}{}
	resp.Json(&data)
	if data.Files.File2 == "" {
		t.Error("invalid response files:", resp.Text())
	}
	if data.Form.Name == "" {
		t.Error("invalid response forms:", resp.Text())
	}

}
