package requests

import (
	"net/http"
	"net/http/cookiejar"
)

// cookies only save to Client.Jar
// session.initCookies is temporary
func (session *Session) SetCookie(cookie *http.Cookie) *Session {
	session.initCookies = append(session.initCookies, cookie)
	return session
}

func (session *Session) ClearCookies() {
	session.initCookies = session.initCookies[0:0]
}

func (session *Session) clientLoadCookies() {
	if len(session.initCookies) > 0 {
		// session.httpreq.AddCookie(cookie) // Do not use httpreq here, this is a temporary cookie, client will drop it in next request
		// 1. set session cookies
		session.Client.Jar.SetCookies(session.httpreq.URL, session.initCookies)
		// 2. clear init cookies
		session.ClearCookies()
	}
}

/**
1. The new Jar is a reference to old jar
2. We could use jar in multiple routines safely.
**/
func (session *Session) clientCloneCookies(client *http.Client) {
	if client != nil && client.Jar != nil {
		switch jar := client.Jar.(type) {
		case *cookiejar.Jar:
			session.Client.Jar = jar
		default:
			session.Client.Jar = client.Jar
		}
	}
}
