package main

import (
	. "catalog/src"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	//kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	//httptransport "github.com/go-kit/kit/transport/http"
	//stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os/signal"
	"syscall"
)

const (
	defaultPort = "8080"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var (
		addr     = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	//fieldKeys := []string{"method", "error"}
	//requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
	//	Namespace: "my_group",
	//	Subsystem: "catalog_service",
	//	Name:      "request_count",
	//	Help:      "Number of requests received.",
	//}, fieldKeys)
	//requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//	Namespace: "my_group",
	//	Subsystem: "order_service",
	//	Name:      "request_latency_microseconds",
	//	Help:      "Total duration of requests in microseconds.",
	//}, fieldKeys)
	//countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
	//	Namespace: "my_group",
	//	Subsystem: "order_service",
	//	Name:      "count_result",
	//	Help:      "The result of each count method.",
	//}, []string{})

	db := GetMongoDB()
	var svc GolfService
	{
		repository, err := NewRepo(db, logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = NewService(repository, logger)
	}
	//svc = loggingMiddleware{logger, svc}
	//svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

	mux := http.NewServeMux()
	httpLogger := log.With(logger, "component", "http")
	mux.Handle("/catalog", MakeHandler(svc, httpLogger))
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
