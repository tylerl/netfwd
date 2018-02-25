// Copyright (c) 2014 Tyler Larson, MIT license applies. See LICENSE

package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type config struct {
	listen   string
	remote   string
	ltls     bool
	rtls     bool
	cacerts  string
	certfile string
	keyfile  string
	verbose  bool
	dump     bool
}

var conf config
var ltlsconf tls.Config
var rtlsconf tls.Config
var _cacheRemoteName string

func main() {
	getconfig()
	if conf.verbose {
		ps := map[bool]string{false: "PLAIN", true: "TLS"}
		info("Listen on %s [%s]", conf.listen, ps[conf.ltls])
		info("Connect to %s [%s]", conf.remote, ps[conf.rtls])
	}

	lsock, err := net.Listen("tcp", conf.listen)
	if lsock == nil {
		fatal("Cannot listen on %s [%v]", conf.listen, err)
	}

	for {
		conn, err := lsock.Accept()
		if conn == nil {
			fatal("Accept failed [%v]", err)
		} else if conf.verbose {
			info("[Connected] %s", conn.RemoteAddr())
		}
		go communicate(conn)
	}
}

func communicate(local net.Conn) {
	var c1, c2 chan int

	remoteAddr := local.RemoteAddr()

	if conf.verbose {
		defer func() {
			info("[Disconnected] %s", remoteAddr)
		}()
		c1 = make(chan int)
		c2 = make(chan int)
	}

	if conf.ltls {
		local = tls.Server(local, &ltlsconf)
	}

	remote, err := net.Dial("tcp", conf.remote)
	if remote == nil {
		nonfatal("Failed to connect to %s [%v]", conf.remote, err)
		local.Close()
		return
	}

	if conf.rtls {
		tlsremote := tls.Client(remote, &rtlsconf)
		if _cacheRemoteName == "" {
			_cacheRemoteName, _, _ = net.SplitHostPort(conf.remote)
		}
		err := tlsremote.Handshake()
		if err == nil && !rtlsconf.InsecureSkipVerify {
			err = tlsremote.VerifyHostname(_cacheRemoteName)
		}
		if err != nil {
			nonfatal("SSL Certificate not validated for %s [%v]", _cacheRemoteName, err)
			tlsremote.Close()
			local.Close()
			return
		}
		remote = tlsremote
	}

	doCopy := func(in, out net.Conn, c chan int) {
		var reader io.Reader = in
		if conf.dump {
			reader = io.TeeReader(in, os.Stdout)
		}
		io.Copy(out, reader)

		out.Close()
		if c != nil {
			close(c)
		}
	}

	go doCopy(local, remote, c1)
	go doCopy(remote, local, c2)

	if c1 != nil && c2 != nil {
		// we use these channels to prevent this fn from going out
		// of scope till the goroutines are done, thus preventing the
		// defered function from running till then.
		<-c1
		<-c2
	}
}

func getconfig() {
	parseoptions()

	if conf.listen == "" {
		usagefatal("missing local <address>:<port>")
	}
	if conf.remote == "" {
		usagefatal("missing remote <address>:<port>")
	}

	if !strings.Contains(conf.listen, ":") {
		conf.listen = "0.0.0.0:" + conf.listen
	}

	if !strings.Contains(conf.remote, ":") {
		lh, lp, _ := net.SplitHostPort(conf.listen)
		if strings.Contains(conf.remote, ".") {
			// we have IP address; so add same port as conf.listen
			conf.remote = conf.remote + ":" + lp
		} else {
			// else we have port, so assume connect localhost
			if lh == "0.0.0.0" {
				lh = "127.0.0.1"
			}
			conf.remote = lh + ":" + conf.remote
		}
	}
	if conf.certfile != "" {
		conf.ltls = true
		if conf.keyfile == "" {
			conf.keyfile = conf.certfile
		}
		kp, err := tls.LoadX509KeyPair(conf.certfile, conf.keyfile)
		if err != nil {
			fatal("Failed to load keypair [%v]", err)
		}
		ltlsconf.Certificates = []tls.Certificate{kp}
	}
	if conf.rtls {
		rtlsconf.ServerName = strings.Split(conf.remote, ":")[0]
		if conf.cacerts != "" {
			if conf.cacerts == "-" {
				rtlsconf.InsecureSkipVerify = true
			} else {
				rtlsconf.RootCAs = x509.NewCertPool()
				data, err := ioutil.ReadFile(conf.cacerts)
				if err == nil {
					ok := rtlsconf.RootCAs.AppendCertsFromPEM(data)
					if !ok {
						err = errors.New("Invalid certificate data")
					}

				}
				if err != nil {
					fatal("Failed to load CA certificates [%v]", err)
				}
			}
		}
	}
}

func fatal(s string, a ...interface{}) {
	log.Fatalf(s, a...)
}
func nonfatal(s string, a ...interface{}) {
	log.Printf(s, a...)
}
func info(s string, a ...interface{}) {
	log.Printf(s, a...)
}
