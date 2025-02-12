package middlewares

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP overrides the RemoteAddr field of the request with data from the
// relevant headers, if present.
func ClientIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip := extractIP(r); ip != nil {
			r.RemoteAddr = ip.String()
		}

		next.ServeHTTP(w, r)
	})
}

// we use infra-client-ip primarily, to avoid issues where we need to figure out
// who owns each IP address in the X-Forwarded-For header
//
// but, as a backup, we'll also accept the first in the X-Forwarded-For chain
func extractIP(r *http.Request) net.IP {
	strIP := strings.TrimSpace(r.Header.Get("Infra-Client-Ip"))
	if ip := net.ParseIP(strIP); ip != nil {
		return ip
	}

	strIP, _, _ = strings.Cut(r.Header.Get("X-Forwarded-For"), ",")
	strIP = strings.TrimSpace(strIP)
	if ip := net.ParseIP(strIP); ip != nil {
		return ip
	}

	return nil
}
