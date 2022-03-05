package requests

import "net/http"

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
	newSession := Sessions()
	newSession.isdebug = session.isdebug

	// 1. clone temp cookies
	newSession.initCookies = session.initCookies

	// 2. clone client cookies
	newSession.clientCloneCookies(session.Client)
	return newSession
}
