package main

import (
	"log"

	"github.com/leosimoesp/go-metric/api"
	"github.com/valyala/fasthttp"
)

func main() {
	handler := api.NewApiHandler()
	r := handler.GetRouter()

	r.POST("/metric/{key}", handler.PostMetric)

	r.GET("/metric/{key}/sum", handler.GetMetricSum)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
