package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fasthttp/router"

	"github.com/leosimoesp/go-metric/internal/app/metricdata"
	"github.com/leosimoesp/go-metric/internal/app/service"
	"github.com/leosimoesp/go-metric/pkg/errorwrapper"
	"github.com/leosimoesp/go-metric/pkg/log"
	"github.com/valyala/fasthttp"
)

type routerImpl struct {
	router    *router.Router
	metricSrv service.MetricService
}

func NewApiHandler() routerImpl {
	return routerImpl{
		metricSrv: service.NewMetricService(),
		router:    router.New(),
	}
}

//PostMetric save a new metric
func (r routerImpl) PostMetric(ctx *fasthttp.RequestCtx) {
	key := ctx.UserValue("key")

	body := ctx.Request.Body()

	var input metricdata.MetricInput

	err := json.Unmarshal(body, &input)

	if err != nil {
		errResp := metricdata.ErrorResponse{
			Message: err.Error(),
		}
		errJson, _ := json.Marshal(errResp)
		ctx.Error(string(errJson), fasthttp.StatusUnprocessableEntity)
	}

	if key != nil {
		log.Logger().Infof("PostMetric %s", key)
		_, err := r.metricSrv.Save(key.(string), input)

		var erroWrap errorwrapper.ErrorWrapper

		if err != nil && errors.As(err, &erroWrap) {
			errResp := metricdata.ErrorResponse{
				Message: erroWrap.Message,
			}
			errJson, _ := json.Marshal(errResp)
			ctx.Error(string(errJson), erroWrap.Code)
			log.Logger().Error(err)
		}

		fmt.Fprintf(ctx, "%s\n", "{}")
		ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	}
}

//PostMetric save a new metric
func (r routerImpl) GetMetricSum(ctx *fasthttp.RequestCtx) {

	key := ctx.UserValue("key")

	if key != nil {
		log.Logger().Infof("GetMetricSum %s", key)
		sum, err := r.metricSrv.CalculateSumMetrics(key.(string))

		var erroWrap errorwrapper.ErrorWrapper

		if err != nil && errors.As(err, &erroWrap) {
			errResp := metricdata.ErrorResponse{
				Message: erroWrap.Message,
			}
			errJson, _ := json.Marshal(errResp)
			ctx.Error(string(errJson), erroWrap.Code)
			log.Logger().Error(err)
		}

		sumResponse := metricdata.MetricSumResponse{Value: int(sum)}
		respson, _ := json.Marshal(sumResponse)

		fmt.Fprintf(ctx, "%s\n", string(respson))
		ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	}

}

func (r routerImpl) GetRouter() *router.Router {
	return r.router
}
