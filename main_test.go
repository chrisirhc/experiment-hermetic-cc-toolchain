package main

import (
	"log"
	"net"
	"net/http"
	"testing"

	ieproxy "github.com/mattn/go-ieproxy"
)

func TestPacfile(t *testing.T) {
	listener, err := listenAndServeWithClose("127.0.0.1:0", http.FileServer(http.Dir("pacfile_examples")))
	serverBase := "http://" + listener.Addr().String() + "/"
	if err != nil {
		t.Fatal(err)
	}

	// test inactive proxy
	proxy := ieproxy.ProxyScriptConf{
		Active:           false,
		PreConfiguredURL: serverBase + "simple.pac",
	}
	out := proxy.FindProxyForURL("http://google.com")
	if out != "" {
		t.Error("Got: ", out, "Expected: ", "")
	}
	proxy.Active = true

	pacSet := []struct {
		pacfile  string
		url      string
		expected string
	}{
		{
			serverBase + "direct.pac",
			"http://google.com",
			"",
		},
		{
			serverBase + "404.pac",
			"http://google.com",
			"",
		},
		{
			serverBase + "simple.pac",
			"http://google.com",
			"127.0.0.1:8",
		},
		{
			serverBase + "simple.pac",
			"https://google.com",
			"127.0.0.1:8",
		},
		{
			serverBase + "multiple.pac",
			"http://google.com",
			"127.0.0.1:8081",
		},
		{
			serverBase + "except.pac",
			"http://imgur.com",
			"localhost:9999",
		},
		{
			serverBase + "except.pac",
			"http://example.com",
			"",
		},
		{
			"",
			"http://example.com",
			"",
		},
		{
			" ",
			"http://example.com",
			"",
		},
		{
			"wrong_format",
			"http://example.com",
			"",
		},
	}
	for _, p := range pacSet {
		proxy.PreConfiguredURL = p.pacfile
		out := proxy.FindProxyForURL(p.url)
		if out != p.expected {
			t.Error("Got: ", out, "Expected: ", p.expected)
		}
	}
	listener.Close()
}

var multipleMap map[string]string

func init() {
	multipleMap = make(map[string]string)
	multipleMap["http"] = "127.0.0.1"
	multipleMap["ftp"] = "128"
}

func listenAndServeWithClose(addr string, handler http.Handler) (net.Listener, error) {

	var (
		listener net.Listener
		err      error
	)

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(listener.(*net.TCPListener))
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	return listener, nil
}
