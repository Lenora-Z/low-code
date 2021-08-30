package utils

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy 创建反向代理处理方法
func NewReverseProxy(urll *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		urll.RawQuery = urll.Query().Encode() //编码特殊字符
		logrus.Infof("destURL=%v", *urll)

		req.Host = urll.Host //重要
		req.URL.Host = urll.Host
		req.URL.Scheme = urll.Scheme
		req.URL.Path = urll.Path
		req.URL.RawQuery = urll.RawQuery
	}
	errHandle := func(rw http.ResponseWriter, req *http.Request, err error) {
		if err != nil {
			logrus.Errorln(*req.URL, ", err= ", err)
			rw.WriteHeader(http.StatusBadGateway)
		}
	}
	return &httputil.ReverseProxy{Director: director, ErrorHandler: errHandle}
}
