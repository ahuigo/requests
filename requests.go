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
	"time"

	"github.com/ahuigo/requests/rerrors"
)

var respHandler func(*Response) error

func SetRespHandler(fn func(*Response) error) {
	respHandler = fn
}

// Post -
func (session *Session) Run(origurl string, args ...interface{}) (resp *Response, err error) {
	if _, err = session.BuildRequest(session.httpreq.Method, origurl, args...); err != nil {
		return nil, err
	}
	dumpCurl := session.RequestDebug()
	startTime := time.Now()
	res, err := session.Client.Do(session.httpreq)

	if err != nil {
		err := rerrors.Wrapf(rerrors.NetworkError, err, "%s %s", session.httpreq.Method, origurl)
		return nil, err
	}

	resp = &Response{
		R:           res,
		startTime:   startTime,
		endTime:     time.Now(),
		httpreq:     session.httpreq,
		client:      session.Client,
		isdebug:     session.isdebug,
		isdebugBody: session.isdebugBody,
		dumpCurl:    dumpCurl,
	}
	resp.ResponseDebug()
	resp.SetStartEndTime(startTime, time.Now()).Body()
	session.reset()

	// global respnse hander & session response handler
	if session.respHandler != nil {
		err = session.respHandler(resp)
	} else if respHandler != nil {
		err = respHandler(resp)
	}
	return resp, err
}
