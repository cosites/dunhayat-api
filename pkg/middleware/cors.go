package middleware

import (
	"net/http"
	"strings"

	"dunhayat-api/pkg/config"
)

func CORS(corsConfig *config.CORSConfig) func(http.Handler) http.Handler {
	if corsConfig == nil {
		corsConfig = &config.CORSConfig{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization"},
			AllowCredentials: false,
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Method == http.MethodOptions {
				origin := r.Header.Get("Origin")
				if origin != "" && isOriginAllowed(
					origin,
					corsConfig.AllowedOrigins,
				) {
					w.Header().Set(
						"Access-Control-Allow-Origin",
						origin,
					)
				}

				if len(corsConfig.AllowedMethods) > 0 {
					methods := strings.Join(
						corsConfig.AllowedMethods,
						", ",
					)
					w.Header().Set(
						"Access-Control-Allow-Methods",
						methods,
					)
				}

				if len(corsConfig.AllowedHeaders) > 0 {
					headers := strings.Join(
						corsConfig.AllowedHeaders,
						", ",
					)
					w.Header().Set(
						"Access-Control-Allow-Headers",
						headers,
					)
				}

				if corsConfig.AllowCredentials {
					w.Header().Set(
						"Access-Control-Allow-Credentials",
						"true",
					)
				}

				w.WriteHeader(http.StatusOK)
				return
			}

			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(
				origin,
				corsConfig.AllowedOrigins,
			) {
				w.Header().Set(
					"Access-Control-Allow-Origin",
					origin,
				)
			}

			if corsConfig.AllowCredentials {
				w.Header().Set(
					"Access-Control-Allow-Credentials",
					"true",
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
		if domain, ok := strings.CutPrefix(allowed, "*."); ok {
			return strings.HasSuffix(
				origin,
				"."+domain,
			) || origin == domain
		}
	}
	return false
}
