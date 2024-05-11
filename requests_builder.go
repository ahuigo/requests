package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Session struct {
	httpreq     *http.Request
	Client      *http.Client
	isdebug     bool
	isdebugBody bool
	respHandler func(*Response) error
	// session header
	gHeader     map[string]string
	initCookies []*http.Cookie
	initContext context.Context
	// connection
	doNotCloseBody bool
	// retry
	retryCount         int
	retryWaitTime      time.Duration
	retryConditionFunc func(*Response, error) bool
	clientTrace        *clientTrace
}

type Header map[string]string
type Params map[string]string
type ParamsArray map[string][]string
type Datas map[string]string     // for post form urlencode
type FormData map[string]string  // for post multipart/form-data
type Json map[string]interface{} // for Json map
type Jsoni interface{}           // for Json interface
type Files map[string]string     // name ,filename
// type AnyData interface{}         // for AnyData
type ContentType string

const (
	ContentTypeNone       ContentType = ""
	ContentTypeFormEncode ContentType = "application/x-www-form-urlencoded"
	ContentTypeFormData   ContentType = "multipart/form-data"
	ContentTypeJson       ContentType = "application/json"
	ContentTypePlain      ContentType = "text/plain"
)

// Auth - {username,password}
type Auth []string

// New request session
func NewSession() *Session {
	return R()
}

// New request session
// @params method  GET|POST|PUT|DELETE|PATCH
func R() *Session {
	var gHeader = map[string]string{
		"User-Agent": "Go-requests-" + getVersion(),
	}
	session := &Session{
		gHeader: gHeader,
	}
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

// BuildRequest
func (session *Session) BuildRequest(method, origurl string, args ...interface{}) (*http.Request, error) {
	var params map[string]string
	var paramsArray map[string][]string
	var datas []map[string]string // form data
	var files []map[string]string //file data
	dataType := ContentTypeNone
	bodyBytes := []byte{}

	session.httpreq.Method = strings.ToUpper(method)
	for _, arg := range args {
		switch arg := arg.(type) {
		case context.Context:
			session.SetContext(arg)
		// arg is Header , set to request header
		case Header:
			for k, v := range arg {
				session.httpreq.Header.Set(k, v)
			}
		case Auth:
			session.httpreq.SetBasicAuth(arg[0], arg[1])
		case *http.Cookie:
			session.SetCookie(arg)
		case ContentType:
			dataType = arg
		case Params:
			params = arg
		case ParamsArray:
			paramsArray = arg
		case Datas:
			dataType = ContentTypeFormEncode
			datas = append(datas, arg)
		case FormData:
			dataType = ContentTypeFormData
			datas = append(datas, arg)
		case Files:
			dataType = ContentTypeFormData
			files = append(files, arg)
		case string:
			if dataType == "" {
				dataType = ContentTypeFormEncode
			}
			bodyBytes = []byte(arg)
		case []byte:
			if dataType == "" {
				dataType = ContentTypePlain
			}
			bodyBytes = arg
		case Json, Jsoni:
			dataType = ContentTypeJson
			bodyBytes = session.buildJSON(arg)
		default:
			dataType = ContentTypeJson
			bodyBytes = session.buildJSON(arg)
		}
	}

	URL, err := buildURLParams(origurl, params, paramsArray)
	if err != nil {
		return nil, err
	}
	if URL.Scheme == "" || URL.Host == "" {
		err = &url.Error{Op: "parse", URL: origurl, Err: fmt.Errorf("failed")}
		return nil, err
	}

	switch dataType {
	case ContentTypeFormEncode:
		session.setContentType("application/x-www-form-urlencoded")
		if len(datas) > 0 {
			formEncodeValues := session.buildFormEncode(datas...)
			session.setBodyFormEncode(formEncodeValues)
		} else {
			session.setBodyBytes(bodyBytes)
		}
	case ContentTypeFormData:
		// multipart/form-data
		session.buildFilesAndForms(files, datas)
	case ContentTypeJson:
		session.setContentType("application/json")
		session.setBodyBytes(bodyBytes)
	case ContentTypePlain:
		session.setContentType("text/plain")
		session.setBodyBytes(bodyBytes)
	}
	if session.httpreq.Body == nil && session.httpreq.Method != "GET" {
		session.httpreq.Body = http.NoBody
	}

	// set header
	for key, value := range session.gHeader {
		session.httpreq.Header.Set(key, value)
	}

	session.httpreq.URL = URL

	// set context
	if ctx := session.initContext; ctx != nil {
		session.httpreq = session.httpreq.WithContext(ctx)
	}

	// set trace context
	if session.isdebug {
		trace := clientTraceNew(session.httpreq.Context())
		session.clientTrace = trace
		session.httpreq = session.httpreq.WithContext(trace.ctx)
	}

	// set host
	host := session.httpreq.Header.Get("Host")
	if host != "" {
		session.httpreq.Host = host
	} else {
		session.httpreq.Host = session.httpreq.URL.Host
	}

	session.clientLoadCookies()
	// fmt.Printf("session:%#v\n", session.httpreq)
	// fmt.Printf("session-url:%#v\n", session.httpreq.URL.String())
	return session.httpreq, nil

}
func (session *Session) reset() {
	session.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	session.clientTrace = nil

}

func (session *Session) Clone() *Session {
	newSession := R()
	newSession.isdebug = session.isdebug
	newSession.isdebugBody = session.isdebugBody

	// 1. clone temp cookies
	newSession.initCookies = session.initCookies

	// 2. clone client cookies
	newSession.clientCloneCookies(session.Client)
	return newSession
}

func (session *Session) setContentType(ct string) {
	// session.httpreq.Header.Get("Content-Type") == "" &&
	if ct != "" {
		session.httpreq.Header.Set("Content-Type", ct)
	}
}

// set form urlencode
func (session *Session) setBodyFormEncode(Forms url.Values) {
	data := Forms.Encode()
	session.httpreq.Body = ioutil.NopCloser(strings.NewReader(data))
	session.httpreq.ContentLength = int64(len(data))
}

// set body
func (session *Session) setBodyBytes(data []byte) {
	session.httpreq.Body = ioutil.NopCloser(bytes.NewReader(data))
	session.httpreq.ContentLength = int64(len(data))
}

// upload file and form
// build to body format
func (session *Session) buildFilesAndForms(files []map[string]string, datas []map[string]string) {
	//handle file multipart
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for _, data := range datas {
		for k, v := range data {
			w.WriteField(k, v)
		}
	}

	for _, file := range files {
		for k, v := range file {
			part, err := w.CreateFormFile(k, v)
			if err != nil {
				fmt.Printf("Upload %s failed!", v)
				panic(err)
			}
			file := openFile(v)
			_, err = io.Copy(part, file)
			if err != nil {
				panic(err)
			}
		}
	}

	w.Close()
	// set file header example:
	// "Content-Type": "multipart/form-data; boundary=------------------------7d87eceb5520850c",
	session.httpreq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	session.httpreq.ContentLength = int64(b.Len())
	session.httpreq.Header.Set("Content-Type", w.FormDataContentType())
}

// build post Form encode
func (session *Session) buildFormEncode(datas ...map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for _, data := range datas {
		for key, value := range data {
			Forms.Add(key, value)
		}
	}
	return Forms
}

func (session *Session) buildJSON(data interface{}) []byte {
	jsonBytes, _ := json.Marshal(data)

	// fmt.Printf("a1=%#v,jsons=%#v\nahui\n", data, string(jsonBytes))
	return jsonBytes
}
