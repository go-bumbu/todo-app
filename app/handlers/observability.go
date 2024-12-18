package handlrs

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Observability handler adds an /metrics endpoint for prometheus
// other observability features can and will be added in the future
func Observability() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		content := `
<a href="/metrics">/metrics</a>
`
		_, _ = fmt.Fprint(writer, content)

	})
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
