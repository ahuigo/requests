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
	R         *http.Response
	body      []byte
	text      string
	session   *Session
	startTime time.Time
	endTime   time.Time
}

func (resp *Response) ResponseDebug() {
	if !resp.session.isdebug {
		return
	}
	fmt.Println("===========ResponseDebug ============")

	message, err := httputil.DumpResponse(resp.R, false)
	if err != nil {
		return
	}

	fmt.Println(string(message))
	fmt.Println("-----------ResponseDebug(end) ------------")
}

func (resp *Response) Body() []byte {
	var err error
	if resp.body != nil {
		return resp.body
	}
	defer resp.R.Body.Close()

	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" && resp.session.Header.Get("Accept-Encoding") != "" {
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
	if resp.body == nil {
		resp.Body()
	}
	resp.text = string(resp.body)
	return resp.text
}

func (resp *Response) Size() int {
	if resp.body == nil {
		resp.Body()
	}
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
	httpreq := resp.session.httpreq
	client := resp.session.Client

	cookies = client.Jar.Cookies(httpreq.URL)

	return cookies

}
