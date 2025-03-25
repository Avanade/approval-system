package fallback

import "net/http"

type FallbackController interface {
	NotFound(w http.ResponseWriter, r *http.Request)
}
