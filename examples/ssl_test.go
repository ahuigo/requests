package examples

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/ahuigo/requests"
)

func TestSkipSsl(t *testing.T) {
	// 1. create tls test server
	ts := createHttpbinServer(2)
	defer ts.Close()

	session := requests.R()

	// 2. fake CA certificate
	// session.SetCaCert("conf/rootCA.crt")

	// 3. skip ssl
	session = session.SkipSsl(true)

	// 4. send get request
	resp, err := session.Get(ts.URL + "/get?a=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}

func TestSetTransport(t *testing.T) {
	// 1. create tls test server
	ts := createHttpbinServer(2)
	defer ts.Close()

	session := requests.R()

	// 3. skip ssl & proxy connect
	tsp := session.GetTransport()
	_ = tsp
	tsp.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// not connect to a proxy server,, keep pathname only
		return net.Dial("tcp", ts.URL[strings.LastIndex(ts.URL, "/")+1:])
	}
	tsp.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// 4. send get request
	resp, err := session.Get(ts.URL + "/get?a=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}

func TestSslCertSelf(t *testing.T) {
	// 1. create tls test server
	ts := createHttpbinServer(1)
	defer ts.Close()

	session := requests.R()
	// 2. certs
	certs := x509.NewCertPool()
	for _, c := range ts.TLS.Certificates {
		roots, err := x509.ParseCertificates(c.Certificate[len(c.Certificate)-1])
		if err != nil {
			log.Fatalf("error parsing server's root cert: %v", err)
		}
		for _, root := range roots {
			certs.AddCert(root)
		}
	}

	// 3. skip ssl & proxy connect
	tsp := session.GetTransport()
	_ = tsp
	tsp.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		// not connect to a proxy server,, keep pathname only
		return net.Dial("tcp", ts.URL[strings.LastIndex(ts.URL, "/")+1:])
	}
	tsp.TLSClientConfig = &tls.Config{
		// InsecureSkipVerify: true,
		RootCAs: certs,
	}

	// 4. send get request
	resp, err := session.Get(ts.URL + "/get?a=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}

// go test -timeout 6000s -run '^TestSslCert$'   github.com/ahuigo/requests/examples -v -httptest.serve=127.0.0.1:443
func TesSslCertCa(t *testing.T) {
	// 1. create tls test server
	ts := createHttpbinServer(2)
	defer ts.Close()

	session := requests.R()

	// 2. fake CA certificate
	session.SetCaCert("conf/rootCA.crt")

	url := strings.Replace(ts.URL, "127.0.0.1", "local.com", 1) + "/get?a=1"
	t.Log(url)
	// time.Sleep(10 * time.Minute)
	// 4. send get request
	resp, err := session.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() == "" {
		t.Fatal(resp.Text())
	}
}
