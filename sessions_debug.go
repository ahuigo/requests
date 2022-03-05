package requests

import (
	"fmt"
	"net/http/httputil"
)

func (session *Session) RequestDebug() {
	if !session.isdebug {
		return
	}
	fmt.Println("===========Go RequestDebug !============")
	curl := BuildCurlRequest(session.httpreq, session.Client.Jar)
	fmt.Println(curl)
	message, err := httputil.DumpRequestOut(session.httpreq, false)
	if err != nil {
		return
	}
	fmt.Println(string(message))
	fmt.Println("-----Go RequestDebug(end)!------------")
}

func (session *Session) SetDebug(debug bool) *Session {
	session.isdebug = debug
	return session
}
