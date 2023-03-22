/**
 * refer to: git@github.com:go-resty/resty.git
 */
package examples

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"testing"
	"time"

	"github.com/ahuigo/requests"
)

// test context: cancel multi
func TestSetContextCancelMulti(t *testing.T) {
	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Microsecond)
		n, err := w.Write([]byte("TestSetContextCancel: response"))
		t.Logf("%s Server: wrote %d bytes", time.Now(), n)
		t.Logf("%s Server: err is %v ", time.Now(), err)
	}, 0)
	defer ts.Close()

	// client
	ctx, cancel := context.WithCancel(context.Background())
	client := requests.R().SetContext(ctx)
	go func() {
		time.Sleep(1 * time.Microsecond)
		cancel()
	}()

	// first
	_, err := client.Get(ts.URL + "/get")
	if !errIsContextCancel(err) {
		t.Fatalf("Got unexpected error: %v", err)
	}

	// second
	_, err = client.Get(ts.URL + "/get")
	if !errIsContextCancel(err) {
		t.Fatalf("Got unexpected error: %v", err)
	}
}

// test context: cancel with chan
func TestSetContextCancelWithChan(t *testing.T) {
	ch := make(chan struct{})
	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			ch <- struct{}{} // tell test request is finished
		}()
		t.Logf("%s Server: %v %v", time.Now(), r.Method, r.URL.Path)
		ch <- struct{}{} // tell test request is canceld
		t.Logf("%s Server: call canceld", time.Now())

		<-ch // wait for client to finish request
		n, err := w.Write([]byte("TestSetContextCancel: response"))
		// FIXME? test server doesn't handle request cancellation
		t.Logf("%s Server: wrote %d bytes", time.Now(), n)
		t.Logf("%s Server: err is %v ", time.Now(), err)

	}, 0)
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-ch // wait for server to start request handling
		cancel()
	}()

	_, err := requests.R().SetContext(ctx).Get(ts.URL + "/get")
	t.Logf("%s:client:is canceled", time.Now())

	ch <- struct{}{} // tell server to continue request handling
	t.Logf("%s:client:tell server to continue", time.Now())

	<-ch // wait for server to finish request handling

	if !errIsContextCancel(err) {
		t.Fatalf("Got unexpected error: %v", err)
	}
}

// test with trace context
func TestContextWithTrace(t *testing.T) {
	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TestSetContextWithTrace: response"))
	}, 0)
	defer ts.Close()

	//1. Create Trace context
	traceInfo := struct {
		dnsDone     time.Time
		connectDone time.Time
	}{}

	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) {
			traceInfo.dnsDone = time.Now()
			t.Log(time.Now(), "ConnectStart:", "network=", network, ",addr=", addr)
		},
		ConnectDone: func(network, addr string, err error) {
			traceInfo.connectDone = time.Now()
			t.Log(time.Now(), "ConnectDone:", "network=", network, ",addr=", addr)
		},
	}
	ctx := httptrace.WithClientTrace(context.Background(), trace)

	//2. Send request with Trace context
	session := requests.R().SetContext(ctx)
	params := requests.Params{"name": "ahuigo", "page": "1"}
	_, err := session.Get(ts.URL+"/get", params)
	if err != nil {
		t.Fatal(err)
	}
	if traceInfo.connectDone.Sub(traceInfo.dnsDone) <= 0 {
		t.Fatal("Bad trace info")
	}

}
