package requests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/alessio/shellescape"
)

// BuildRequest -
func BuildRequest(method string, origurl string, args ...interface{}) (req *http.Request, err error) {
	// call request Get
	args = append(args, Method(method))
	req, err = NewSession().BuildRequest(origurl, args...)
	return
}

func BuildCurlRequest(req *http.Request, args ...interface{}) (curl string) {
	// 1. generate curl request
	curl = "curl -X " + req.Method + " "
	// req.Host + req.URL.Path + "?" + req.URL.RawQuery + " " + req.Proto + " "
	headers := getHeaders(req)
	for _, kv := range *headers {
		curl += `-H ` + shellescape.Quote(kv[0]+": "+kv[1]) + ` `
	}

	// 1.2 generate curl with cookies
	for _, arg := range args {
		switch arg := arg.(type) {
		case *cookiejar.Jar:
			cookies := arg.Cookies(req.URL)
			if len(cookies) > 0 {
				curl += ` -H ` + shellescape.Quote(dumpCookies(cookies)) + " "
			}
		}
	}

	// body
	if req.Body != nil {
		buf, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // important!!
		curl += `-d ` + shellescape.Quote(string(buf))
	}

	curl += " " + shellescape.Quote(req.URL.String())
	return curl
}

func dumpCookies(cookies []*http.Cookie) string {
	sb := strings.Builder{}
	sb.WriteString("Cookie: ")
	for _, cookie := range cookies {
		sb.WriteString(cookie.Name + "=" + url.QueryEscape(cookie.Value) + "&")
	}
	return strings.TrimRight(sb.String(), "&")
}

// getHeaders
func getHeaders(req *http.Request) *[][2]string {
	headers := [][2]string{}
	for k, vs := range req.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	// fmt.Printf("%#v\n", headers)
	// sort headers
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return &headers
}
