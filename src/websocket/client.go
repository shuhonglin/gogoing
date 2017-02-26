package websocket

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"net/http"
)

var ErrBadHandshake = errors.New("websocket: bad handshake")
var errInvalidCompression = errors.New("websocket: invalid compression negotiation")

func NewClient(netConn net.Conn, u *url.URL, requestHeader http.Header, readBufSize, writeBufSize int) (c *Conn)


