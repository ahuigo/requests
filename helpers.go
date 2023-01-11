package requests

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
)

var VERSION string = "v0.0.0"

func getVersion() string {
	_, filename, _, _ := runtime.Caller(0)
	versionFile := path.Dir(filename) + "/version"
	version, _ := ioutil.ReadFile(versionFile)
	VERSION = strings.TrimSpace(string(version))
	return VERSION

}

func init() {
	getVersion()
}

// open file for post upload files
func openFile(filename string) *os.File {
	r, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return r
}

// handle URL params
func buildURLParams(userURL string, params map[string]string, paramsArray map[string][]string) (*url.URL, error) {
	if strings.HasPrefix(userURL, "/") {
		userURL = "http://localhost" + userURL
	}
	parsedURL, err := url.Parse(userURL)

	if err != nil {
		return nil, err
	}

	values := parsedURL.Query()

	for key, value := range params {
		values.Set(key, value)
	}
	for key, vals := range paramsArray {
		for _, v := range vals {
			values.Add(key, v)
		}
	}
	parsedURL.RawQuery = values.Encode()
	return parsedURL, nil
}
