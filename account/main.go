package main

import (
	"account/src"
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

const (
	defaultPort		= "8080"
	dbsource		= "mongodb://mongo:27017"
	dbname			= "gokitexample"
)

func main() {
	var (
		addr = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "main module", log.DefaultTimestampUTC)

	err := registerService()
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}

	db, err := src.GetMongoDB(ctx, dbsource, dbname)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	var repo src.Repository
	repo = src.NewRepo(db, logger)

	fieldKeys := []string{"method"}

	var service src.Service
	service = src.NewService(repo, logger)
	service = src.NewLoggingService(log.With(logger, "component", "account"), service)
	service = src.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "account_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		service,
	)

	httpLogger := log.With(logger, "component", "http")
	mux := http.NewServeMux()
	mux.Handle("/user", src.MakeHandler(service, httpLogger))
	mux.Handle("/user/", src.MakeHandler(service, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthcheck", healthcheck)

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

func envString(env string, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user service is good")
}


func registerService() error {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(envString("PORT", defaultPort))
	if err != nil {
		return err
	}

	address, err := hostname()
	if err != nil {
		return err
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "account-service"
	registration.Name = "account-service"
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)

	return nil
}

func hostname() (string, error) {
	return "host.docker.internal", nil
}
