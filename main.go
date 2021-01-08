package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/", forward)

	if err := http.ListenAndServe(":4000", proxyMux); err != nil {
		log.Fatalln("Server has failed to start:", err.Error())
	}
}

func forward(res http.ResponseWriter, req *http.Request) {
	isBackend := strings.HasPrefix(req.RequestURI, "/api/") // Here we assume the only backend is "api"
	if isBackend {
		req.URL.Path = strings.SplitN(req.RequestURI, "api", 2)[1]
		reverseProxyTo("8080", res, req)
		return
	} else {
		reverseProxyTo("3000", res, req)
		return
	}
}

func reverseProxyTo(port string, res http.ResponseWriter, req *http.Request) {
	origin, _ := url.Parse("http://localhost:" + port)
	proxy := httputil.NewSingleHostReverseProxy(origin)

	log.Println("Forwarding to:", origin, req.URL.Path)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
