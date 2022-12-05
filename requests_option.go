package requests

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// set timeout s = second
func (session *Session) SetTimeout(n time.Duration) *Session {
	session.Client.Timeout = n
	return session
}

func (session *Session) Close() {
	session.httpreq.Close = true
}

func (session *Session) Proxy(proxyurl string) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyurl)
	if err != nil {
		fmt.Println("Set proxy failed")
		return
	}
	session.Client.Transport = &http.Transport{
		Proxy:           http.ProxyURL(urlproxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// SkipSsl
func (session *Session) SkipSsl(v bool) *Session{
	transport := session.Client.Transport
	if transport == nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	switch tp := transport.(type){
		case *http.Transport:
			tlsConf := tp.TLSClientConfig
			if tlsConf ==nil{
				tlsConf = &tls.Config{}
			}
			tlsConf.InsecureSkipVerify= v

	}
	session.Client.Transport = transport
	return session
}

func (session *Session) SetRespHandler(fn func(*Response) error) *Session {
	session.respHandler = fn
	return session
}

// SetMethod
func (session *Session) SetMethod(method string) *Session {
	session.httpreq.Method = strings.ToUpper(method)
	return session
}

// SetHeader
func (session *Session) SetHeader(key, value string) *Session {
	session.httpreq.Header.Set(key, value)
	return session
}

// Set global header
func (session *Session) SetGlobalHeader(key, value string) *Session {
	session.gHeader[key] = value
	return session
}

// Get global header
func (session *Session) GetGlobalHeader() map[string]string {
	return session.gHeader
}

// Del global header
func (session *Session) DelGlobalHeader(key string) *Session {
	delete(session.gHeader, key)
	return session
}
