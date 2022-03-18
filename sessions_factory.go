package requests

import (
	"net/http"
	"net/http/cookiejar"
)

var gHeader = map[string]string{
	"User-Agent": "Go-requests-" + getVersion(),
}

// Set global header
func SetHeader(key, value string) {
	if value == "" {
		delete(gHeader, key)
		return
	}
	gHeader[key] = value
}

// New request session
func R() *Session {
	return NewSession()
}

// New request session
// @params method  GET|POST|PUT|DELETE|PATCH
func NewSession() *Session {
	session := new(Session)
	session.reset()

	session.Client = NewHttpClient()

	return session
}

func NewHttpClient() *http.Client {
	// cookiejar.New source code return jar, nil
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return client
}

func (session *Session) reset() {
	session.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	session.Header = &session.httpreq.Header
	for key, value := range gHeader {
		session.httpreq.Header.Set(key, value)
	}
}

func (session *Session) Clone() *Session {
	newSession := R()
	newSession.isdebug = session.isdebug

	// 1. clone temp cookies
	newSession.initCookies = session.initCookies

	// 2. clone client cookies
	newSession.clientCloneCookies(session.Client)
	return newSession
}
