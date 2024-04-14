package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func WithCORSMiddleware(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	return c.Handler(h)
}
