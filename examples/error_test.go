package examples

import (
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
	"github.com/ahuigo/requests/rerrors"
)

func TestError(t *testing.T) {
	_, err := requests.Get("http://127.0.0.1:12346/connect-refused")
	if err2, ok := err.(*rerrors.Error); !ok {
		t.Fatalf("unexpected error:%+v", err)
	} else {
		switch err2.ErrType {
		case rerrors.NetworkTimeout:
			t.Log(err2.ErrType)
		case rerrors.NetworkError:
			t.Log(err2.ErrType)
		case rerrors.URLError:
			t.Log(err2.ErrType)
		default:
			t.Log(err2.ErrType)
		}
	}
	t.Log(err)
}
func TestErrorTimeout(t *testing.T) {
	ts := createHttpbinServer()
	defer ts.Close()

	// resp, err := requests.Get("https://httpbin.org/json")
	_, err := requests.R().SetTimeout(1).Get(ts.URL + "/sleep/10")
	if err2, ok := err.(*rerrors.Error); !ok {
		t.Fatalf("unexpected error:%+v", err)
	} else {
		switch err2.ErrType {
		case rerrors.NetworkTimeout:
			t.Log(err2.ErrType)
		default:
			t.Fatalf("unexpected error type:%+v", err2.ErrType)
		}
	}
}

func TestErrorNetwork(t *testing.T) {
	_, err := requests.Get("http://127.0.0.1:12346/connect-refused")
	if err2, ok := err.(*rerrors.Error); !ok {
		t.Fatalf("unexpected error:%+v", err)
	} else {
		switch err2.ErrType {
		case rerrors.NetworkError:
			t.Log(err2.ErrType)
		default:
			t.Fatalf("expect type:NetworkError, unexpected error type:%+v", err2.ErrType)
		}
	}
}

func TestErrorURL(t *testing.T) {
	_, err := requests.Get("xxxx")
	if err2, ok := err.(*rerrors.Error); !ok {
		t.Fatalf("unexpected error:%+v", err)
	} else {
		switch err2.ErrType {
		case rerrors.URLError:
			t.Log(err2.ErrType)
		default:
			t.Fatalf("expect type:NetworkError, unexpected error type:%+v", err2.ErrType)
			t.Fatalf("unexpected error type:%+v", err2.ErrType)
		}
	}
}
