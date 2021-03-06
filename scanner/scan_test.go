package scanner

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"testing"

	"github.com/empijei/go-sslscan/scanner/tls_dirty"
)

var gosuppciph = []uint16{
	tls_dirty.TLS_RSA_WITH_RC4_128_SHA,
	tls_dirty.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	tls_dirty.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls_dirty.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls_dirty.TLS_RSA_WITH_AES_128_CBC_SHA256,
	tls_dirty.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls_dirty.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls_dirty.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	tls_dirty.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	tls_dirty.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls_dirty.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	tls_dirty.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	tls_dirty.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls_dirty.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls_dirty.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls_dirty.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
}

func TestScanHostWithProto(t *testing.T) {
	l, err := fooserver()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = l.Close() }()

	//Default golang tls listener configuration
	expected := gosuppciph
	//cs, err := ScanHost("localhost:8043")
	cs, err := ScanHostWithProto("localhost:8043", TLS12)
	if err != nil {
		t.Error(err)
		return
	}
	if len(expected) != len(cs) {
		t.Errorf("ScanHost: expected %d ciphers but got %d", len(expected), len(cs))
		buf, _ := json.MarshalIndent(cs, "", " ")
		t.Errorf(string(buf))
		return
	}
	for i, c := range expected {
		if cs[i].ID != c {
			t.Errorf("ScanHost: expected <%X> in position %d but got <%X>", c, i, cs[i].Name)
		}
	}
}

func fooserver() (net.Listener, error) {
	// Simple static webserver:
	cer, err := tls.LoadX509KeyPair("data/test_server.crt", "data/test_server.key")
	if err != nil {
		log.Fatal(err)
	}
	l, err := tls.Listen("tcp", ":8043", &tls.Config{
		Certificates:             []tls.Certificate{cer},
		MinVersion:               tls.VersionSSL30,
		PreferServerCipherSuites: true,
		CipherSuites:             gosuppciph})
	go func() {
		for {
			c, errr := l.Accept()
			if errr != nil {
				return
			}
			_, _ = c.Write([]byte("hello"))
			_ = c.Close()
		}
	}()
	return l, err
}

/*
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
	"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		"TLS_RSA_WITH_AES_256_CBC_SHA",
		"TLS_RSA_WITH_AES_256_GCM_SHA384",
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
		"TLS_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_RSA_WITH_AES_128_CBC_SHA",
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA",
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
		/*TODO
		Preferred TLSv1.2  256 bits  ECDHE-RSA-AES256-GCM-SHA384   Curve P-256 DHE 256
		Accepted  TLSv1.2  256 bits  ECDHE-RSA-AES256-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.2  256 bits  AES256-GCM-SHA384
		Accepted  TLSv1.2  256 bits  AES256-SHA
		Accepted  TLSv1.2  128 bits  ECDHE-RSA-AES128-GCM-SHA256   Curve P-256 DHE 256
		Accepted  TLSv1.2  128 bits  ECDHE-RSA-AES128-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.2  128 bits  AES128-GCM-SHA256
		Accepted  TLSv1.2  128 bits  AES128-SHA
		Accepted  TLSv1.2  112 bits  ECDHE-RSA-DES-CBC3-SHA        Curve P-256 DHE 256
		Accepted  TLSv1.2  112 bits  DES-CBC3-SHA
		Preferred TLSv1.1  256 bits  ECDHE-RSA-AES256-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.1  256 bits  AES256-SHA
		Accepted  TLSv1.1  128 bits  ECDHE-RSA-AES128-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.1  128 bits  AES128-SHA
		Accepted  TLSv1.1  112 bits  ECDHE-RSA-DES-CBC3-SHA        Curve P-256 DHE 256
		Accepted  TLSv1.1  112 bits  DES-CBC3-SHA
		Preferred TLSv1.0  256 bits  ECDHE-RSA-AES256-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.0  256 bits  AES256-SHA
		Accepted  TLSv1.0  128 bits  ECDHE-RSA-AES128-SHA          Curve P-256 DHE 256
		Accepted  TLSv1.0  128 bits  AES128-SHA
		Accepted  TLSv1.0  112 bits  ECDHE-RSA-DES-CBC3-SHA        Curve P-256 DHE 256
		Accepted  TLSv1.0  112 bits  DES-CBC3-SHA
		Preferred SSLv3    256 bits  AES256-SHA
		Accepted  SSLv3    128 bits  AES128-SHA
		Accepted  SSLv3    112 bits  DES-CBC3-SHA
*/
