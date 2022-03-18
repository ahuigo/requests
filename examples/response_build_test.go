package examples

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ahuigo/requests"
)

func TestResponseBuilder(t *testing.T) {
	var data = 1
	responseBytes, _ := json.Marshal(data)

	respRecorder := httptest.NewRecorder()
	respRecorder.Write(responseBytes)

	// build response
	wrapResp := requests.BuildResponse(respRecorder.Result())

	var ndata int
	wrapResp.Json(&ndata)
	if ndata != data {
		t.Fatalf("expect response:%v", data)
	}

}
