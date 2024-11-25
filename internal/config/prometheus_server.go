package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	prometheusServerHostEnvName = "PROMETHEUS_SERVER_HOST"
	prometheusServerPortEnvName = "PROMETHEUS_SERVER_PORT"
)

// PrometheusServerConfig to get server address
type PrometheusServerConfig interface {
	Address() string
}

type prometheusServerConfigConfig struct {
	host string
	port string
}

// NewPrometheusServerConfig prometheus server parameters
func NewPrometheusServerConfig() (PrometheusServerConfig, error) {
	host := os.Getenv(prometheusServerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusServerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("prometheus port not found")
	}

	return &prometheusServerConfigConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusServerConfigConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
