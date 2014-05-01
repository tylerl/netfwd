A simple TCP forwarder with optional encryption
===============================================

`netfwd` is a general-purpose TCP forwarder with optional encryption.

Specify a listening port and an address to connect to, and it will
accept inbound connections, then connect to the remote address and
forward traffic between the two sockets. TLS encryption is supported
on both inbound and outbound connections. 

### Usage

        -l   [addr:]port   Address and port to listen on, addr defaults to 0.0.0.0.
        -r   addr:port     Remote address to connect to.
        -c   certfile      Certificate(s) for inbound encryption [turns on inbound TLS].
        -k   keyfile       Certificate's private key. Defaults to using <certfile>.
        -ca  certs         File containing valid CAs for outbound TLS. "-" turns off verification
        -rt                Turns on TLS for outbound connection.
        -v                 Verbose mode.
        -d                 Dump all traffic to STDOUT.
        -h -? -help        Full help text.

The `-l` and `-r` parameters are required. The inbound address is assumed to be `0.0.0.0`
(i.e. ANY IPv4) if you don't specify. If you specify only an outbound address, then it
uses the same port as for inbound. If you instead specify only an outbound port, then it
assumes you're connecting to `127.0.0.1`. It tells the difference by looking for a `.`. When
in doubt, specify both the address and port.

You turn on inbound encryption by specifying a certificate (the `-c` option), and
outbound encryption by specifying the `-rt` option.

You can control the valid CA certificates checked for outbound connections using
the `-ca` option to point to a file containing one or more PEM certificates. If you
specify `-` as the file, then validation is disabled. If you don't use the `-ca`
option, then it attempts to use your OS's default certificate list.

The `-v` option instructs `netfwd` to display more information; useful for when running interactively
or for logging activity.

The `-d` option (useful for debugging other services) will dump all forwarded traffic to STDOUT.
