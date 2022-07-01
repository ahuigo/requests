package examples

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
)

func TestSendCookie(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	resp, err := requests.Get(ts.URL+"/cookie/count",
		requests.Header{"Cookie": "id_token=1234"},
	)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	resp.Json(&data)
	fmt.Println(data)

}

// Test session Cookie
func TestSessionCookie(t *testing.T) {
	cookie1 := &http.Cookie{
		Name:  "name1",
		Value: "value1",
		Path:  "/",
	}
	cookie2 := &http.Cookie{
		Name:  "name2",
		Value: "value2",
	}
	session := requests.R().SetDebug()

	// 1. set cookie1
	session.SetCookie(cookie1).Get("https://www.httpbin.org/get")

	// 2. set cookie2 and get all cookies
	resp, err := session.Get("https://www.httpbin.org/get", cookie2)
	if err != nil {
		t.Fatal(err)
	}
	cookies := map[string]string{}
	// cookies's type is `[]*http.Cookies`
	for _, c := range resp.Cookies() {
		if _, exists := cookies[c.Name]; exists {
			t.Fatal("duplicated cookie:", c.Name, c.Value)
		}
		cookies[c.Name] = c.Value
	}
	if cookies["name1"] != "value1" || cookies["name2"] != "value2" {
		t.Fatalf("Failed to send valid cookie(%+v)", resp.Cookies())
	}

}

// Test session Cookie
func TestSessionCookieWithClone(t *testing.T) {
	url := "https://www.httpbin.org/get"
	cookie1 := &http.Cookie{
		Name:  "name1",
		Value: "value1",
		Path:  "/",
	}
	cookie2 := &http.Cookie{
		Name:  "name2",
		Value: "value2",
	}
	session := requests.R().SetDebug()

	// 1. set cookie1
	session.SetCookie(cookie1).Get(url)

	// 2. set cookie2 and get all cookies
	session = session.Clone()
	resp, err := session.Get(url, cookie2)
	if err != nil {
		t.Fatal(err)
	}
	cookies := map[string]string{}
	// cookies's type is `[]*http.Cookies`
	for _, c := range resp.Cookies() {
		if _, exists := cookies[c.Name]; exists {
			t.Fatal("duplicated cookie:", c.Name, c.Value)
		}
		cookies[c.Name] = c.Value
	}
	if resp.GetCookie("name1") != "value1" || resp.GetCookie("name2") != "value2" {
		t.Fatalf("Failed to send valid cookie(%+v)", resp.Cookies())
	}

}

// Test Set-Cookie
func TestResponseCookie(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	// resp, err := requests.Get("https://httpbin.org/json")
	resp, err := requests.Get(ts.URL + "/cookie/count")
	if err != nil {
		t.Fatal(err)
	}

	cs := resp.Cookies()
	if len(cs) == 0 {
		t.Fatalf("require cookies, body=%s", resp.Body())
	}
}

func TestResponseBuildCookie(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	// resp, err := requests.Get("https://httpbin.org/json")
	resp, err := requests.Get(ts.URL + "/cookie/count")
	if err != nil {
		t.Fatal(err)
	}

	// build new resposne
	resp.R.Body = ioutil.NopCloser(bytes.NewBuffer(resp.Body())) // important!!
	resp = requests.BuildResponse(resp.R)
	cs := resp.Cookies()
	if len(cs) == 0 {
		t.Fatalf("require cookies, headers=%#v, body=%s", resp.Header(), resp.Body())
	}
	findCount := false
	for _, c := range cs {
		if c.Name == "count" && c.Value == "1" {
			findCount = true
		}
	}
	if !findCount {
		t.Fatalf("could not find cookie, dumpResponse=%s", resp.GetDumpResponse())
	}
}
