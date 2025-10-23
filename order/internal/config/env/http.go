package env

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host     string `env:"HTTP_HOST,required"`
	Port     int    `env:"HTTP_PORT,required"`
	Timeout  string `env:"HTTP_TIMEOUT" envDefault:"15s"`
	ShutDown string `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"15s"`

	InventoryAddr string `env:"INVENTORY_SERVICE_ADDRESS,required"`
	PaymentAddr   string `env:"PAYMENT_SERVICE_ADDRESS,required"`
}

type httpConfig struct {
	raw httpEnvConfig
}

// NewHTTPConfig parses environment variables into an httpConfig instance
func NewHTTPConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse HTTP env config: %w", err)
	}
	return &httpConfig{raw: raw}, nil
}

// Address builds the HTTP server address
func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.raw.Host, strconv.Itoa(c.raw.Port))
}

// Timeout returns parsed request timeout
func (c *httpConfig) Timeout() time.Duration {
	d, err := time.ParseDuration(c.raw.Timeout)
	if err != nil {
		return 15 * time.Second // дефолт, если что-то не так
	}
	return d
}

func (c *httpConfig) ShutDown() time.Duration {
	d, err := time.ParseDuration(c.raw.ShutDown)
	if err != nil {
		return 15 * time.Second // дефолт, если что-то не так
	}
	return d
}

func (c *httpConfig) InventoryAddr() string {
	return c.raw.InventoryAddr
}

func (c *httpConfig) PaymentAddr() string {
	return c.raw.PaymentAddr
}
