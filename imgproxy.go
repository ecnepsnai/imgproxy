package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
)

var Version = "dev"
var BuildDate = "N/A"

var log = logtic.Log.Connect("imgproxy")

type httpServer struct{}

func (h *httpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		log.Warn("Invalid request: ignore index")
		rw.WriteHeader(404)
		return
	}

	if r.Method != "GET" && r.Method != "HEAD" {
		log.Warn("Invalid request: unsupported method %s", r.Method)
		rw.WriteHeader(404)
		return
	}

	path := strings.Split(r.URL.Path[1:], ".")[0]
	urlBytes, err := base64.RawURLEncoding.DecodeString(path)
	if err != nil {
		log.Warn("Invalid request: invalid base64 encoded data")
		rw.WriteHeader(404)
		return
	}

	uri, err := url.Parse(string(urlBytes))
	if err != nil {
		log.Warn("Invalid request: invalid base64 encoded URL")
		rw.WriteHeader(404)
		return
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		log.Warn("Invalid request: unsupported url scheme %s", uri.Scheme)
		rw.WriteHeader(404)
		return
	}

	pr, err := http.NewRequest(r.Method, uri.String(), nil)
	if err != nil {
		log.PError("Error forming HTTP request", map[string]interface{}{
			"url":   uri.String(),
			"error": err.Error(),
		})
		rw.WriteHeader(500)
		return
	}
	for key, value := range r.Header {
		if key == "Host" {
			value = []string{uri.Host}
		}
		pr.Header.Set(key, value[0])
	}

	resp, err := http.DefaultClient.Do(pr)
	if err != nil {
		log.PError("Network error performing request", map[string]interface{}{
			"url":   uri.String(),
			"error": err.Error(),
		})
		rw.WriteHeader(500)
		return
	}

	if resp.ContentLength > 52428800 {
		log.Warn("Invalid request: response too large %d (%s)", resp.ContentLength, logtic.FormatBytesB(uint64(resp.ContentLength)))
		rw.WriteHeader(413)
		return
	}

	rw.WriteHeader(resp.StatusCode)
	for key, value := range resp.Header {
		rw.Header().Set(key, value[0])
	}
	rw.Header().Set("X-Imgproxy-Version", Version)

	var length int64 = 0
	if r.Method == "GET" {
		wrote, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.PError("Error writing response bytes", map[string]interface{}{
				"url":   uri.String(),
				"error": err.Error(),
			})
		}
		length = wrote
	}
	log.PInfo("Proxied request", map[string]interface{}{
		"remote_addr": r.RemoteAddr,
		"method":      r.Method,
		"url":         uri.String(),
		"length_b":    length,
		"length":      logtic.FormatBytesB(uint64(length)),
	})
}

func main() {
	if len(os.Args) > 2 {
		if os.Args[1] == "-u" {
			fmt.Printf("%s\n", base64.URLEncoding.EncodeToString([]byte(os.Args[2])))
			os.Exit(0)
		}
	}

	fmt.Printf("imgproxy v%s, built on %s, runtime %s\n", Version, BuildDate, runtime.Version())

	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Open()

	s := &httpServer{}

	httpsListener, err := startHTTPSListener()
	if err != nil {
		log.PFatal("Error starting HTTPS server", map[string]interface{}{
			"error": err.Error(),
		})
	}
	httpListener, err := startHTTPListener()
	if err != nil {
		log.PFatal("Error starting HTTP server", map[string]interface{}{
			"error": err.Error(),
		})
	}

	go func() {
		if err := http.Serve(httpsListener, s); err != nil {
			log.Panic("error starting https server: %s", err.Error())
		}
	}()
	go func() {
		if err := http.Serve(httpListener, s); err != nil {
			log.Panic("error starting http server: %s", err.Error())
		}
	}()

	for {
		time.Sleep(1 * time.Minute)
	}
}
