package request

import (
	"github.com/imroc/req"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func client(timeout int64) *http.Client {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   2 * time.Minute,
	}
}

func Get(url string, timeout int64, v ...interface{}) (*req.Resp, error) {
	if timeout <= 0 {
		timeout = 5
	}
	return req.Get(url, client(timeout))
}


func Post(url string, timeout int64, v ...interface{}) (*req.Resp, error) {
	if timeout <= 0 {
		timeout = 5
	}
	return req.Post(url, client(timeout))
}

func Delete(url string, timeout int64, v ...interface{}) (*req.Resp, error) {
	if timeout <= 0 {
		timeout = 5
	}
	return req.Delete(url, client(timeout))
}

func Put(url string, timeout int64, v ...interface{}) (*req.Resp, error) {
	if timeout <= 0 {
		timeout = 5
	}
	return req.Put(url, client(timeout))
}

func Head(url string, timeout int64, v ...interface{}) (*req.Resp, error) {
	if timeout <= 0 {
		timeout = 5
	}
	return req.Head(url, client(timeout))
}