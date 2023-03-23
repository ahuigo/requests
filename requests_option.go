package requests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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

func (session *Session) SetContext(ctx context.Context) *Session {
	// build request: session.httpreq = session.httpreq.WithContext(ctx)
	session.initContext = ctx
	return session
}

func (session *Session) Close() {
	session.httpreq.Close = true
}

func (session *Session) SetProxy(proxyurl string) *Session {
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyurl)
	if err != nil {
		println("SetProxy, Bad proxy:" + proxyurl)
		return session
	}
	transport := session.GetTransport()
	transport.Proxy = http.ProxyURL(urlproxy)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return session
}

// In generally, you could SystemCertPool instead of NewCertPool to keep existing certs.
func (session *Session) SetCaCert(caCertPath string) *Session {
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConf := session.getTLSClientConfig()
	tlsConf.RootCAs = caCertPool
	return session
}

// SkipSsl
func (session *Session) SkipSsl(v bool) *Session {
	tlsConf := session.getTLSClientConfig()
	tlsConf.InsecureSkipVerify = v
	return session
}

func (session *Session) getTLSClientConfig() *tls.Config {
	tp := session.GetTransport()
	tlsConf := tp.TLSClientConfig
	if tlsConf == nil {
		tlsConf = &tls.Config{}
	}
	return tlsConf
}

func (session *Session) GetTransport() *http.Transport {
	transport := session.Client.Transport
	if transport == nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
		session.Client.Transport = transport
	}
	return transport.(*http.Transport)
}

func (session *Session) SetTransport(tsp *http.Transport) *Session {
	session.Client.Transport = tsp
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
