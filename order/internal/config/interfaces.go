package config

import "time"

type OrderConfig interface {
	Address() string
	Timeout() time.Duration
	InventoryAddr() string
	PaymentAddr() string
}

type PostgresConfig interface {
	DSN() string
	DBName() string
	MigPath() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
