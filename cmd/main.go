package main

import (
	"fmt"
	"os"

	"github.com/leosimoesp/go-metric/api"
	"github.com/leosimoesp/go-metric/config"
	"github.com/leosimoesp/go-metric/pkg/log"
	"github.com/valyala/fasthttp"
)

func init() {
	checkEnvVars([]string{config.Port, config.RemoveMetricsIntervalInMin})
}

func main() {

	port := os.Getenv(config.Port)

	handler := api.NewApiHandler()
	r := handler.GetRouter()

	r.POST("/metric/{key}", handler.PostMetric)

	r.GET("/metric/{key}/sum", handler.GetMetricSum)

	log.Logger().Infof("Starting server at port %s", port)

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), r.Handler); err != nil {
		log.Logger().Fatal(err)
	}
}

// validate all env vars by name into os path
func checkEnvVars(envVarsNames []string) {
	for _, v := range envVarsNames {
		if os.Getenv(v) == "" {
			log.Logger().Panicf("Env var %s is mandatory at path app", v)
		}
	}
}
