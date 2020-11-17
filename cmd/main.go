package main

import (
	"net"
	"net/http"
	"os"

	"github.com/KarrenAeris/http/pkg/banners"

	"github.com/KarrenAeris/http/cmd/app"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	bannerSvc := banners.NewService()
	server := app.NewServer(mux, bannerSvc)
	server.Init()
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	return srv.ListenAndServe()
}
