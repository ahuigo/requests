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
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type Response struct {
	R              *http.Response
	Attempt        int
	body           []byte
	doNotCloseBody bool
	httpreq        *http.Request
	client         *http.Client
	TraceInfo      *TraceInfo
	isdebug        bool
	isdebugBody    bool
	startTime      time.Time
	endTime        time.Time
	dumpCurl       string
	dumpResponse   string
}

func BuildResponse(response *http.Response) *Response {
	r := &Response{
		R: response,
	}
	// resp.R.Body = ioutil.NopCloser(bytes.NewBuffer(resp.Body())) // important!!
	r._DumpResponse(true)
	r.Body()
	return r
}

func (resp *Response) GetReq() (req *http.Request) {
	return resp.httpreq
}

func (resp *Response) SetDoNotCloseBody() *Response {
	resp.doNotCloseBody = true
	return resp
}

func (resp *Response) SetClientReq(req *http.Request, client *http.Client) *Response {
	resp.client = client
	resp.httpreq = req
	return resp
}

func (resp *Response) SetStartEndTime(start, end time.Time) *Response {
	resp.startTime = start
	resp.endTime = end
	return resp
}

func (resp *Response) ResponseDebug() {
	if !resp.isdebug {
		return
	}
	fmt.Println("===========ResponseDebug ============")
	err := resp._DumpResponse(resp.isdebugBody)
	if err != nil {
		return
	}
	fmt.Println(resp.dumpResponse)
	fmt.Println("========== ResponseDebug(end) ============")
}

func (resp *Response) _DumpResponse(isdebugBody bool) error {
	//resp.isdebug
	message, err := httputil.DumpResponse(resp.R, isdebugBody)
	resp.dumpResponse = string(message)
	return err
}

func (resp *Response) GetDumpCurl() string {
	return resp.dumpCurl
}
func (resp *Response) GetDumpResponse() string {
	if resp.dumpResponse == "" {
		if resp.R.Body == nil {
			resp.R.Body = ioutil.NopCloser(bytes.NewBuffer(resp.Body())) // important!!
		}
		resp._DumpResponse(true)
	}
	return resp.dumpResponse
}

func (resp *Response) Body() []byte {
	var err error
	if resp.body != nil {
		return resp.body
	}
	resp.body = []byte{}
	if !resp.doNotCloseBody {
		defer resp.R.Body.Close()
	}

	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}

	resp.body, err = ioutil.ReadAll(Body)
	if err != nil {
		return nil
	}

	return resp.body
}

func (resp *Response) Text() string {
	return string(resp.body)
}

func (resp *Response) Size() int {
	return len(resp.body)
}

func (resp *Response) RaiseForStatus() (code int, err error) {
	code = resp.R.StatusCode
	if resp.R.StatusCode >= 400 && resp.R.StatusCode != 401 {
		err = errors.New(resp.Text())
	}
	return
}

func (resp *Response) StatusCode() (code int) {
	return resp.R.StatusCode
}

func (resp *Response) Time() time.Duration {
	return resp.endTime.Sub(resp.startTime)
}

func (resp *Response) Header() http.Header {
	return resp.R.Header
}

func (resp *Response) SaveFile(filename string) error {
	if resp.body == nil {
		resp.Body()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.body)
	f.Sync()

	return err
}

func (resp *Response) Json(v interface{}) error {
	if resp.body == nil {
		resp.Body()
	}
	return json.Unmarshal(resp.body, v)
}

func (resp *Response) Cookies() (cookies []*http.Cookie) {
	httpreq := resp.httpreq
	client := resp.client

	if httpreq == nil || client == nil {
		return resp.R.Cookies()
	}
	// cookies's type is `[]*http.Cookies`
	cookies = client.Jar.Cookies(httpreq.URL)
	return cookies
}

func (resp *Response) GetCookie(key string) (val string) {
	cookies := map[string]string{}
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	val = cookies[key]
	return val
}

func (resp *Response) HasCookie(key string) (exists bool) {
	cookies := map[string]string{}
	for _, c := range resp.Cookies() {
		cookies[c.Name] = c.Value
	}
	_, exists = cookies[key]
	return exists
}
