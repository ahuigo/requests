package examples

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
)

func createHttpbinServer() (ts *httptest.Server) {
	ts = createTestServer(func(w http.ResponseWriter, r *http.Request) {
		switch path := r.URL.Path; {
		case path == "/get":
			getHandler(w, r)
		case strings.HasPrefix(path, "/cookie/"):
			cookieHandler(w, r)
		default:
			_, _ = w.Write([]byte("404 " + path))
		}
	})

	return ts
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	m := map[string]interface{}{
		"headers": dumpRequestHeader(r),
		"args":    parseRequestArgs(r),
		"body":    string(body),
	}
	buf, _ := json.Marshal(m)
	_, _ = w.Write(buf)
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/cookie/count":
		count := "1"
		cookie, err := r.Cookie("count")
		if err == nil {
			i, _ := strconv.Atoi(cookie.Value)
			count = strconv.Itoa(i + 1)
		}
		http.SetCookie(w, &http.Cookie{Name: "count", Value: url.QueryEscape(count)})
		w.Header().Set("Content-Type", "application/json")

		body, _ := ioutil.ReadAll(r.Body)
		m := map[string]interface{}{
			"args":    parseRequestArgs(r),
			"body":    string(body),
			"count":   count,
			"headers": dumpRequestHeader(r),
		}
		buf, _ := json.Marshal(m)
		_, _ = w.Write(buf)
	default:
		_, _ = w.Write([]byte("404 " + r.URL.Path))
	}
}

func dumpRequestHeader(req *http.Request) string {
	var res strings.Builder
	headers := sortHeaders(req)
	for _, kv := range headers {
		res.WriteString(kv[0] + ": " + kv[1] + "\n")
	}
	return res.String()
}

// sortHeaders
func sortHeaders(request *http.Request) [][2]string {
	headers := [][2]string{}
	for k, vs := range request.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return headers
}

func createEchoServer() (ts *httptest.Server) {
	ts = createTestServer(func(w http.ResponseWriter, r *http.Request) {
		res := dumpRequest(r)
		_, _ = w.Write([]byte(res))
	})

	return ts
}
func parseRequestArgs(request *http.Request) map[string]string {
	query := request.URL.RawQuery
	params := map[string]string{}
	paramsList, _ := url.ParseQuery(query)
	for key, vals := range paramsList {
		params[key] = vals[len(vals)-1]
	}
	return params
}

func dumpRequest(request *http.Request) string {
	var r strings.Builder
	// dump header
	res := request.Method + " " + //request.URL.String() +" "+
		request.Host +
		request.URL.Path + "?" + request.URL.RawQuery + " " + request.Proto + " " +
		"\n"
	r.WriteString(res)
	r.WriteString(dumpRequestHeader(request))
	r.WriteString("\n")

	// dump body
	buf, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // important!!
	r.WriteString(string(buf))
	return r.String()
}

func createTestServer(fn func(w http.ResponseWriter, r *http.Request)) (ts *httptest.Server) {
	ts = httptest.NewServer(http.HandlerFunc(fn))
	return ts
}
