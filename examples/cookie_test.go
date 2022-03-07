package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
)

func TestSendCookie(t *testing.T) {
	resp, err := requests.Get("https://www.httpbin.org/cookies",
		requests.Header{"Cookie": "id_token=1234"},
		requests.Json{"workflow_id": "wfid1234"},
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
	session := requests.Sessions().SetDebug(true)

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
	url = "http://m:4500/echo/get"
	cookie1 := &http.Cookie{
		Name:  "name1",
		Value: "value1",
		Path:  "/",
	}
	cookie2 := &http.Cookie{
		Name:  "name2",
		Value: "value2",
	}
	session := requests.Sessions().SetDebug(true)

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
	if cookies["name1"] != "value1" || cookies["name2"] != "value2" {
		t.Fatalf("Failed to send valid cookie(%+v)", resp.Cookies())
	}

}
