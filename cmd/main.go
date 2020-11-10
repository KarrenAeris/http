package main

import (
	/*
	"bytes"
	"fmt"
	"io"
	"io/ioutil" 
	"strings"
	"time"

	"strconv" */
	"github.com/KarrenAeris/http/pkg/server"
	"net"
	"os"
	//"log"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) error {
	
	srv := server.NewServer(net.JoinHostPort(host, port))

	// srv.Register("/", srv.Requesting("Hi there")) 
	// srv.Register("/about", srv.Requesting("have you kept the promise?")) 
	return srv.Start()
	
}



