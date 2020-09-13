package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/tzusman/super-soft/tls"
)

func main() {

	hostnameRaw, err := exec.Command("hostname").Output()
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.ToLower(strings.TrimSpace(string(hostnameRaw)))
	fmt.Printf("Download ca cert to your phone at http://%s:7001\n", hostname)
	fmt.Printf("Proxy is at https://%s:7000\n", hostname)

	certDir := fmt.Sprintf("%s/.config/super-soft/", os.Getenv("HOME"))
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

	origin, _ := url.Parse("http://localhost:8080/")

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

		http.ListenAndServe(":7001", serveCA)
	}()

	err = http.ListenAndServeTLS(":7000", crtFilepath, keyFilepath, nil)
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
