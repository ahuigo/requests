/* Copyright（2） 2018 by  asmcos and ahuigo .
Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
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

var respHandler func(*Response) error
var gHeader = map[string]string{
	"User-Agent": "Go-requests-" + getVersion(),
}

func SetRespHandler(fn func(*Response) error) {
	respHandler = fn
}

type Session struct {
	httpreq     *http.Request
	Client      *http.Client
	isdebug     bool
	respHandler func(*Response) error
	// global header
	Header      *http.Header
	initCookies []*http.Cookie
}

type Header map[string]string
type Params map[string]string
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
type Method string

// New request session
func R() *Session {
	return NewSession()
}

// New request session
// @params method  GET|POST|PUT|DELETE|PATCH
func NewSession() *Session {
	session := new(Session)
	session.reset()

	session.Client = &http.Client{}

	// cookiejar.New source code return jar, nil
	jar, _ := cookiejar.New(nil)

	session.Client.Jar = jar

	return session
}

// Set global header
func SetHeader(key, value string) {
	if value == "" {
		delete(gHeader, key)
		return
	}
	gHeader[key] = value
}

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
	session.Header.Set(key, value)
	return session
}

// BuildRequest
func (session *Session) BuildRequest(origurl string, args ...interface{}) (*http.Request, error) {
	var params []map[string]string
	var datas []map[string]string // form data
	var files []map[string]string //file data
	dataType := ContentTypeNone
	bodyBytes := []byte{}

	for _, arg := range args {
		switch arg := arg.(type) {
		// arg is Header , set to request header
		case Method:
			session.httpreq.Method = strings.ToUpper(string(arg))
		case Header:
			for k, v := range arg {
				session.httpreq.Header.Set(k, v)
			}
		case Auth:
			session.httpreq.SetBasicAuth(arg[0], arg[1])
		case *http.Cookie:
			session.SetCookie(arg)
		case ContentType:
			session.setContentType(string(arg))
			dataType = arg
		case Params:
			params = append(params, arg)
		case Datas:
			datas = append(datas, arg)
		case FormData:
			dataType = ContentTypeFormData
			datas = append(datas, arg)
		case Files:
			dataType = ContentTypeFormData
			files = append(files, arg)
		case string:
			dataType = ContentTypeFormEncode
			bodyBytes = []byte(arg)
		case []byte:
			dataType = ContentTypePlain
			bodyBytes = arg
		case Json, Jsoni:
			dataType = ContentTypeJson
			bodyBytes = session.buildJSON(arg)
		default:
			dataType = ContentTypeJson
			bodyBytes = session.buildJSON(arg)
		}
	}

	disturl, _ := buildURLParams(origurl, params...)

	switch dataType {
	case ContentTypeFormEncode:
		session.setContentType("application/x-www-form-urlencoded")
		session.setBodyBytes(bodyBytes)
	case ContentTypeFormData:
		// multipart/form-data
		session.buildFilesAndForms(files, datas)
	case ContentTypeJson:
		session.setContentType("application/json")
		session.setBodyBytes(bodyBytes)
	case ContentTypePlain:
		session.setContentType("text/plain")
		session.setBodyBytes(bodyBytes)
	default:
		if len(datas) > 0 {
			session.setContentType("application/x-www-form-urlencoded")
			formEncodeValues := session.buildFormEncode(datas...)
			session.setBodyFormEncode(formEncodeValues)
		}
	}
	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	session.httpreq.URL = URL

	session.clientLoadCookies()
	// fmt.Printf("session:%#v\n", session.httpreq)
	// fmt.Printf("session-url:%#v\n", session.httpreq.URL.String())
	return session.httpreq, nil

}
func (session *Session) setContentType(ct string) {
	if session.httpreq.Header.Get("Content-Type") == "" && ct != "" {
		session.httpreq.Header.Set("Content-Type", ct)
	}
}

// Post -
func (session *Session) Run(origurl string, args ...interface{}) (resp *Response, err error) {
	session.BuildRequest(origurl, args...)
	session.RequestDebug()
	startTime := time.Now()
	res, err := session.Client.Do(session.httpreq)

	if err != nil {
		return nil, errors.New(session.httpreq.Method + " " + origurl + " " + err.Error())
	}

	resp = &Response{
		R:         res,
		startTime: startTime,
		endTime:   time.Now(),
	}
	req_dup := *session
	resp.session = &req_dup
	resp.ResponseDebug()
	session.reset()

	// global respnse hander & session response handler
	if session.respHandler != nil {
		err = session.respHandler(resp)
	} else if respHandler != nil {
		err = respHandler(resp)
	}
	return resp, err
}

// only set forms
func (session *Session) setBodyFormEncode(Forms url.Values) {
	data := Forms.Encode()
	session.httpreq.Body = ioutil.NopCloser(strings.NewReader(data))
	session.httpreq.ContentLength = int64(len(data))
}

// only set forms
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
	session.Header.Set("Content-Type", w.FormDataContentType())
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
