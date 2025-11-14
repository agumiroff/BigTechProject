package config

type InventoryConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DBName() string
	MigrationPath() string
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}
