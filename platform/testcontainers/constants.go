package testcontainers

const (
	// MongoDB container constants
	MongoContainerName = "mongo"
	MongoPort          = "27017"

	// MongoDB environment variables
	MongoImageNameKey  = "MONGO_IMAGE_NAME"
	MongoHostKey       = "MONGO_HOST"
	MongoPortKey       = "MONGO_PORT"
	MongoDatabaseKey   = "MONGO_INITDB_DATABASE"
	MongoUsernameKey   = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey   = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey     = "MONGO_AUTH_DB"
	MongoMigrationsDir = "MONGO_MIGRATIONS_DIR"
	MongoLoggerLevel   = "LOGGER_LEVEL"
	MongoLoggerAsJSON  = "LOGGER_AS_JSON"

	// PostgreSQL container constants
	PostgresContainerName = "postgres"
	PostgresDefaultPort   = "5432"

	// PostgreSQL environment variables (from order .env file)
	PostgresImageNameKey  = "POSTGRES_IMAGE_NAME"
	PostgresHostKey       = "POSTGRES_HOST"
	PostgresPortKey       = "POSTGRES_PORT"
	PostgresDatabaseKey   = "POSTGRES_DB"
	PostgresUserKey       = "POSTGRES_USER"
	PostgresPasswordKey   = "POSTGRES_PASSWORD" //nolint:gosec
	PostgresMigrationPath = "POSTGRES_MIGRATION_PATH"
)
