package mw

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := statusRecorder{ResponseWriter: w, statusCode: ""}

		next.ServeHTTP(&rec, r)

		log.Printf("[%s] %s %s \t\t %s", r.Method, r.RequestURI, time.Since(start), rec.statusCode)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode string
}

func (rec *statusRecorder) WriteHeader(code int) {
	color := "\033[42m"
	if code >= 400 {
		color = "\033[41m"
	}
	rec.statusCode = fmt.Sprintf("%s[%d]\033[0m %s", color, code, http.StatusText(code))

	rec.ResponseWriter.WriteHeader(code)
}
