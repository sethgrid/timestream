package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// endless streaming content
type Content struct {
	f         os.File
	uri       string
	startTime int64
	ticker    *time.Ticker
}

// read as long as the ticker is active.
// the ticker will stop when the connection is closed.
func (c Content) Read(p []byte) (n int, err error) {

	for _ = range c.ticker.C {
		// don't adjust from 0 bytes read. Eventually, we will chew through some internal buffer and panic
		return 0, nil
	}

	// this is never run
	log.Print("stop read")
	return 0, io.EOF
}

// dummy func to continually seek
func (c Content) Seek(offset int64, whence int) (int64, error) {
	return c.f.Seek(0, 0)
}

// times an open connection until client disconnects
func dripHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Don't forget to close the connection:
	defer conn.Close()

	// content keeps a ticker and logs start time.
	content := &Content{}
	content.ticker = time.NewTicker(time.Millisecond * 500)
	content.startTime = time.Now().Unix()
	content.uri = uri

	// strip leading slash and trailing extension
	if len(content.uri) > 1 {
		content.uri = uri[1:strings.Index(uri, ".")]
	}

	// don't block on serving content
	log.Print("New Connection")
	go http.ServeContent(w, r, "", time.Now(), content)

	// reads from hijacked connection. Closed connections result in EOF.
	_, err = bufrw.ReadString('\n')
	if err.Error() == "EOF" {
		info, err := base64.StdEncoding.DecodeString(content.uri)
		if err != nil {
			info = []byte(content.uri)
			log.Print(err)
		}
		log.Printf("duration-seconds: %d %s", time.Now().Unix()-content.startTime, info)
		content.ticker.Stop()
		return
	} else if err != nil {
		log.Printf("error reading from connection: %v", err)
		return
	}
}

func Serve() {
	http.HandleFunc("/", dripHandler)
	http.ListenAndServe(":9999", nil)
	// This does not do what I thought it would :(
	// I can get a panic in ServeContent that terminates the
	// whole process (if the Read method returns anything other than 0 bytes read)
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Recovered from Panic()")
			Serve()
		}
	}()
}

func main() {
	log.Print("starting...")
	Serve()
}
