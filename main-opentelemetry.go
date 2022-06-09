// +build ignore

// NOTE: I think this works, but needs testing. Will only do things if the incoming request has a well-formed traceparent header; can't throw any old string in there
// NOTE: Only forwards OpenTelemetry headers ie traceparent and tracespan. Istio does not use those, it uses zipkin's b3-* ones.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func handle(w http.ResponseWriter, r *http.Request) {
	spew.Dump(r)

	prop := otel.GetTextMapPropagator()
	ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	spew.Dump(ctx)

	msg := os.Getenv("CHAIN_OUTPUT")
	if msg == "" {
		log.Fatal("Must supply output to return")
	}
	fmt.Fprintf(w, "%s\n", msg)

	next := os.Getenv("CHAIN_NEXT")
	if next == "" {
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+next+"/", nil)
	if err != nil {
		log.Fatal(err)
	}
	prop.Inject(ctx, propagation.HeaderCarrier(req.Header))
	spew.Dump(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(body))
}

func main() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handle), "root"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
