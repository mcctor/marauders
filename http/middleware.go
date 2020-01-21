package http

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mcctor/marauders/db"
)

const emptyString = ""

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func ApplyGzipCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
		}
		writer.Header().Set("Content-Encoding", "gzip")
		gzipper := gzip.NewWriter(writer)
		defer gzipper.Close()
		zipperRespWriter := gzipResponseWriter{Writer: gzipper, ResponseWriter: writer}
		next.ServeHTTP(zipperRespWriter, request)
	})
}

func ApplyContentTypeCollectionsJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/vnd.collection+json")
		next.ServeHTTP(writer, request)
	})
}

func ApplyOwnerPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		requestingUser, err := db.GetUser(vars["username"])
		if err != nil {
			http.Error(writer, "{\"status\": \"no user with given username\"}", http.StatusNotFound)
			return
		}
		bearerToken := request.Header.Get("Token")
		if bearerToken == emptyString {
			writer.Header().Set("WWW-Authenticate", "Token")
			http.Error(writer, "{\"status\": \"request unauthorized, no token in request header\"}", http.StatusUnauthorized)
			return
		}
		correctUserAuthToken, err := requestingUser.AuthToken()
		if err != nil {
			http.Error(writer,
				"{\"status\": \"user has no associated authentication token\"}", http.StatusInternalServerError)
			return
		}
		if bearerToken == correctUserAuthToken.Token {
			next.ServeHTTP(writer, request)
		} else {
			http.Error(writer, "{\"status\": \"the bearer token does not match the user's\"}", http.StatusForbidden)
		}
	})
}
