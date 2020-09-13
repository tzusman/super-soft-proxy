package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/tzusman/super-soft/tls"
	"github.com/tzusman/super-soft/util"
)

func main() {

	proxyPort := flag.Int("proxy-port", 7000, "proxy port")
	proxyCAPort := flag.Int("ca-cert-server", 7001, "port where the CA certificate will be served")
	upstreamPort := flag.Int("upstream", 8080, "upstream service port")

	flag.Parse()

	hostname, err := util.GetHostname()
	if err != nil {
		panic(err)
	}

	certDir := fmt.Sprintf("%s/.config/super-soft-proxy/", os.Getenv("HOME"))
	_ = os.Mkdir(certDir, 0755)

	caCrtFilepath := fmt.Sprintf("%s/ca.crt", certDir)
	crtFilepath := fmt.Sprintf("%s/%s.crt", certDir, hostname)
	keyFilepath := fmt.Sprintf("%s/%s.key", certDir, hostname)

	if !fileExists(caCrtFilepath) {
		caCert, err := tls.CreateCertificateAuthority()
		if err != nil {
			panic(err)
		}

		sslCert, err := tls.CreateTLSCertificate(hostname, *caCert)
		if err != nil {
			panic(err)
		}

		err = writeFile(caCrtFilepath, caCert.CertPEM)
		if err != nil {
			panic(err)
		}

		err = writeFile(crtFilepath, sslCert.CertPEM)
		if err != nil {
			panic(err)
		}

		err = writeFile(keyFilepath, sslCert.KeyPEM)
		if err != nil {
			panic(err)
		}
	}

	upstreamURL := fmt.Sprintf("http://localhost:%d/", *upstreamPort)
	origin, _ := url.Parse(upstreamURL)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	go func() {
		// Serve CA over non-tls
		serveCA := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, caCrtFilepath)
		})

		listenPort := fmt.Sprintf(":%d", *proxyCAPort)
		http.ListenAndServe(listenPort, serveCA)
	}()

	fmt.Printf("https://%s:%d is the proxy address\n", hostname, *proxyPort)
	fmt.Printf(" http://%s:%d serves the CA cert for easy device installation\n", hostname, *proxyCAPort)

	listenPort := fmt.Sprintf(":%d", *proxyPort)
	err = http.ListenAndServeTLS(listenPort, crtFilepath, keyFilepath, nil)
	if err != nil {
		panic(err)
	}

}

func writeFile(filepath string, contents []byte) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(contents)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
