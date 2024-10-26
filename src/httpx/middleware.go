package httpx

import (
	"net/http"
	"strings"

	"github.com/baiqll/src-http/src/lib"
)

func HostMiddleware(host string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_host := strings.Split(r.Host, ":")[0]
		if !lib.IsIP(host) && req_host != host{
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found!"))
			return
		}
		next.ServeHTTP(w, r)
    })
}