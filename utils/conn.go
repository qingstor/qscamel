package utils

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

var connPool sync.Pool

type netConn net.Conn
type netDialer net.Dialer

// Dialer is wrapped dialer provided by qingstor go sdk.
//
// We provide this dialer wrapper ReadTimeout & WriteTimeout attributes into connection object.
// This timeout is for individual buffer I/O operation like other language (python, perl... etc),
// so don't bother with SetDeadline or stupid nethttp.Client timeout.
type Dialer struct {
	*net.Dialer
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewDialer will create a new dialer.
func NewDialer(connTimeout, readTimeout, writeTimeout time.Duration) *Dialer {
	d := &net.Dialer{
		DualStack: false,
		Timeout:   connTimeout,
	}
	return &Dialer{d, readTimeout, writeTimeout}
}

// Dial connects to the address on the named network.
func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	c, err := d.Dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	conn := NewConn(c)
	conn.readTimeout = d.ReadTimeout
	conn.writeTimeout = d.WriteTimeout
	return conn, nil
}

// DialContext connects to the address on the named network using
// the provided context.
func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	c, err := d.Dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	conn := NewConn(c)
	conn.readTimeout = d.ReadTimeout
	conn.writeTimeout = d.WriteTimeout
	return conn, nil
}

// Conn is a generic stream-oriented network connection.
type Conn struct {
	netConn
	readTimeout  time.Duration
	writeTimeout time.Duration
	timeoutFunc  func() bool
}

// NewConn will create a new conn.
func NewConn(c netConn) *Conn {
	conn, ok := c.(*Conn)
	if ok {
		return conn
	}
	conn, ok = connPool.Get().(*Conn)
	if !ok {
		conn = new(Conn)
	}
	conn.netConn = c
	return conn
}

// SetReadTimeout will set the conn's read timeout.
func (c *Conn) SetReadTimeout(d time.Duration) {
	if c.readTimeout > 0 {
		c.netConn.SetReadDeadline(time.Time{})
	}
	c.readTimeout = d
}

// SetWriteTimeout will set the conn's write timeout.
func (c *Conn) SetWriteTimeout(d time.Duration) {
	if c.writeTimeout > 0 {
		c.netConn.SetWriteDeadline(time.Time{})
	}
	c.writeTimeout = d
}

// Read will read from the conn.
func (c Conn) Read(buf []byte) (n int, err error) {
	if c.readTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	n, err = c.netConn.Read(buf)
	return
}

// Write will write into the conn.
func (c Conn) Write(buf []byte) (n int, err error) {
	if c.writeTimeout > 0 {
		c.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	n, err = c.netConn.Write(buf)
	return
}

// Close will close the conn.
func (c Conn) Close() (err error) {
	if c.netConn == nil {
		return nil
	}
	err = c.netConn.Close()
	connPool.Put(c)
	c.netConn = nil
	c.readTimeout = 0
	c.writeTimeout = 0
	return
}

// IsTimeoutError will check whether the err is a timeout error.
func IsTimeoutError(err error) bool {
	e, ok := err.(net.Error)
	if ok {
		return e.Timeout()
	}
	return false
}

// DefaultDialer is the default dialer for qscamel.
var DefaultDialer = &Dialer{
	&net.Dialer{
		DualStack: false,
		Timeout:   time.Second * 30,
	},
	time.Second * 30,
	time.Second * 30,
}

// DefaultClient is the default HTTP client for qscamel.
var DefaultClient = &http.Client{
	// We do not use the timeout in http client,
	// because this timeout is for the whole http body read/write,
	// it's unsuitable for various length of files and network condition.
	// We provide a wraper in utils/conn.go of net.Dialer to make io timeout to the http connection
	// for individual buffer I/O operation,
	Timeout: 0,
	Transport: &http.Transport{
		DialContext:           DefaultDialer.DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       time.Second * 20,
		TLSHandshakeTimeout:   time.Second * 10, //Default
		ExpectContinueTimeout: 2 * time.Second,
	},
}
