// netfwd is a general-purpose TCP forwarder with optional encryption.
// Specify a listening port and an address to connect to, and it will
// accept inbound connections, then connect to the remote address and
// forward traffic between the two sockets. TLS encryption is supported
// on both inbound and outbound connections. The -d option (useful for
// troubleshooting other services) will dump all forwarded traffic to STDOUT.
//
// Think of it like stunnel, but with the encryption being optional on
// both sides (and also easier to configure).
//
// Copyright (c) 2014 Tyler Larson, MIT license applies. See LICENSE
package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "netfwd -l [local_addr:]local_port -r remote_addr:remote_port\n"+
		"       [-c certfile] [-k keyfile] [-ca authority_certs] [-rt] [-v] [-d]\n"+
		"       [-h|-?]\n")
}

func usagefatal(s string) {
	fmt.Fprintf(os.Stderr, "%s\n\n", s)
	usage()
	os.Exit(1)
}

func help() {
	fmt.Fprintf(os.Stderr,
		"netfwd is a general-purpose TCP forwarder with optional encryption.\n"+
			"Specify a listening port and an address to connect to, and it will\n"+
			"accept inbound connections, then connect to the remote address and\n"+
			"forward traffic between the two sockets. TLS encryption is supported\n"+
			"on both inbound and outbound connections. The -d option (useful for\n"+
			"debugging other services) will dump all forwarded traffic to STDOUT.\n\n")
	fullUsage()
}

func fullUsage() {
	usage()
	fmt.Fprintf(os.Stderr, "\n"+
		" **  -l   [addr:]port   Address and port to listen on, addr defaults to 0.0.0.0.\n"+
		" **  -r   addr:port     Remote address to connect to.\n"+
		"     -c   certfile      Certificate(s) for inbound encryption [turns on inbound TLS].\n"+
		"     -k   keyfile       Certificate's private key. Defaults to using <certfile>.\n"+
		"     -ca  certs         File containing valid CAs for outbound TLS. \"-\" turns off verification\n"+
		"     -rt                Turns on TLS for outbound connection.\n"+
		"     -v                 Verbose mode.\n"+
		"     -d                 Dump all traffic to STDOUT.\n"+
		"     -h -? -help        Full help text.\n")
}

func parseoptions() {
	if len(os.Args) == 1 {
		fullUsage()
		os.Exit(0)
	}
	flag.Usage = usage
	doHelp := false
	flag.StringVar(&conf.listen, "l", "", "local (listen) host:port")
	flag.StringVar(&conf.remote, "r", "", "remote (connect) host:port")
	flag.BoolVar(&conf.rtls, "rt", false, "Use TLS on remote connection")
	flag.BoolVar(&conf.verbose, "v", false, "Verbose mode")
	flag.BoolVar(&conf.dump, "d", false, "Dump all traffic to stdout")
	flag.StringVar(&conf.cacerts, "ca", "", "CA Certificate file, or \"-\" to disable verification")
	flag.StringVar(&conf.certfile, "c", "", "Local certificate file")
	flag.StringVar(&conf.keyfile, "k", "", "Local key file")
	flag.BoolVar(&doHelp, "h", false, "")
	flag.BoolVar(&doHelp, "help", false, "")
	flag.BoolVar(&doHelp, "?", false, "")
	flag.Parse()
	if doHelp {
		help()
		os.Exit(0)
	}
}
