package main

import (
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/handlers"
	"github.com/kelseyhightower/envconfig"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/client/selector/static"
	metrics "github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/oschwald/geoip2-golang"
	geoip "github.com/paysuper/geoip-service/pkg"
	"github.com/paysuper/geoip-service/pkg/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/Clever/pathio.v3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	GeoIpDbPath   string `envconfig:"MAXMIND_GEOIP_DB_PATH" required:"true"`
	MetricsPort   int    `envconfig:"METRICS_PORT" required:"false" default:"8080"`
	MicroSelector string `envconfig:"MICRO_SELECTOR" required:"false"`
}

type customHealthCheck struct{}

func main() {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Config init failed with error: %s\n", err)
	}

	pathioReader, err := pathio.Reader(cfg.GeoIpDbPath)
	if err != nil {
		log.Fatal(err)
	}

	geoipBuf, err := ioutil.ReadAll(pathioReader)
	if err != nil {
		log.Fatal(err)
	}

	db, err := geoip2.FromBytes(geoipBuf)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	dbMeta := db.Metadata()
	dbBuildDate := time.Unix(int64(dbMeta.BuildEpoch), 0).UTC()
	dbInfo := "Loaded database info:\n" +
		"\tFilename: %s\n" +
		"\tVersion: %d.%d\n" +
		"\tType: %s\n" +
		"\tBuild date: %s\n"

	log.Printf(dbInfo, cfg.GeoIpDbPath, dbMeta.BinaryFormatMajorVersion, dbMeta.BinaryFormatMinorVersion, dbMeta.DatabaseType, dbBuildDate)

	log.Println("Initialize micro service")

	options := []micro.Option{
		micro.Name(geoip.ServiceName),
		micro.Version(geoip.Version),
		micro.WrapHandler(metrics.NewHandlerWrapper()),
	}

	if cfg.MicroSelector == "static" {
		log.Println(`Use micro selector "static"`)
		options = append(options, micro.Selector(static.NewSelector()))
	}

	service := micro.NewService(options...)
	service.Init()

	err = proto.RegisterGeoIpServiceHandler(service.Server(), &geoip.Service{GeoReader: db})

	if err != nil {
		log.Fatal(err)
	}

	initHealth(cfg)
	initPrometheus()

	go func() {
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.MetricsPort), nil); err != nil {
			log.Fatal("Metrics listen failed")
		}
	}()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initHealth(cfg *Config) {
	h := health.New()
	err := h.AddChecks([]*health.Config{
		{
			Name:     "health-check",
			Checker:  &customHealthCheck{},
			Interval: time.Duration(1) * time.Second,
			Fatal:    true,
		},
	})

	if err != nil {
		log.Fatal("Health check register failed")
	}

	log.Printf("Health check listening on :%d", cfg.MetricsPort)

	if err = h.Start(); err != nil {
		log.Fatal("Health check start failed")
	}

	http.HandleFunc("/health", handlers.NewJSONHandlerFunc(h, nil))
}

func initPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
}

func (c *customHealthCheck) Status() (interface{}, error) {
	return "ok", nil
}
