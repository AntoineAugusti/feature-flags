package http

import (
	"log"
	"net"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			getIPAddress(r),
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func getIPAddress(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	ips := []string{
		// Forwarded headers first
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("x-forwarded-for"),
		r.Header.Get("X-FORWARDED-FOR"),
		// Client IP address by default
		ip,
	}

	for _, ip := range ips {
		if len(ip) > 0 {
			return ip
		}
	}
	return ""
}
