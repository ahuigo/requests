package requests

import (
	"fmt"
	"net/http/httputil"
)

func (session *Session) RequestDebug() string {
	if !session.isdebug {
		return ""
	}
	fmt.Println("-----------Go RequestDebug !------------")
	curl := BuildCurlRequest(session.httpreq, session.Client.Jar)
	fmt.Println(curl)
	message, err := httputil.DumpRequestOut(session.httpreq, false)
	if err != nil {
		return ""
	}
	fmt.Println(string(message))
	fmt.Println("-----Go RequestDebug(end)!------------")
	return curl
}

func (session *Session) SetDebug() *Session {
	session.isdebug = true
	return session
}
func (session *Session) SetDebugBody() *Session {
	session.isdebug = true
	session.isdebugBody = true
	return session
}
func (session *Session) SetDoNotCloseBody(flag bool) *Session {
	session.doNotCloseBody = flag
	return session
}
func (session *Session) DisableDebug() *Session {
	session.isdebug = false
	session.isdebugBody = false
	return session
}
