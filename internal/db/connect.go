package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"blogo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connectionParams holds all the database connection configuration parameters
// that are parsed from environment variables and used to establish the connection pool
type connectionParams struct {
	Host              string        // Database host address
	Port              string        // Database port number
	User              string        // Database username
	Password          string        // Database password
	Name              string        // Database name
	SSLMode           string        // SSL mode for the connection
	MaxConnections    int32         // Maximum number of connections in the pool
	MinConnections    int32         // Minimum number of connections to maintain
	MaxConnectionLife time.Duration // Maximum lifetime of a connection
	MaxIdleTime       time.Duration // Maximum time a connection can remain idle
	HealthCheckPeriod time.Duration // How often to perform health checks on connections
	ConnectionTimeout time.Duration // Timeout for establishing new connections
}

// Connect establishes a new database connection pool using the provided configuration.
// It creates the pool, configures it with the specified parameters, and verifies
// the connection by performing a ping test.
//
// Parameters:
//   - cfg: Application configuration containing database connection details
//
// Returns:
//   - *pgxpool.Pool: Configured connection pool
//   - error: Any error that occurred during connection setup
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	// Build the pool configuration from the provided config
	poolConfig, err := buildPoolConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("error building pool configuration: %v", err)
	}

	// Create a new connection pool with the built configuration
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}

	// Verify the connection is working by performing a ping test
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return pool, nil
}

// buildPoolConfig creates a pgxpool configuration from the application config.
// It parses the connection parameters and sets up the connection pool settings
// including connection limits, timeouts, and health check intervals.
//
// Parameters:
//   - cfg: Application configuration containing database settings
//
// Returns:
//   - *pgxpool.Config: Configured pool configuration
//   - error: Any error that occurred during configuration building
func buildPoolConfig(cfg *config.Config) (*pgxpool.Config, error) {
	// Parse and validate connection parameters from the config
	params, err := loadConnectionParams(cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading connection parameters: %v", err)
	}

	// Build the connection string in the format expected by pgx
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		params.Host, params.Port, params.User, params.Password, params.Name, params.SSLMode)

	// Parse the connection string into a pgxpool configuration
	poolConfig, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection URL: %v", err)
	}

	// Configure the connection pool settings
	poolConfig.MaxConns = params.MaxConnections                     // Maximum number of connections
	poolConfig.MinConns = params.MinConnections                     // Minimum number of connections to maintain
	poolConfig.MaxConnLifetime = params.MaxConnectionLife           // Connection lifetime limit
	poolConfig.MaxConnIdleTime = params.MaxIdleTime                 // Maximum idle time for connections
	poolConfig.HealthCheckPeriod = params.HealthCheckPeriod         // Health check frequency
	poolConfig.ConnConfig.ConnectTimeout = params.ConnectionTimeout // Connection establishment timeout

	return poolConfig, nil
}

// loadConnectionParams parses the application configuration and converts string values
// to their appropriate types for database connection parameters. It handles parsing
// of numeric values, durations, and validates the configuration.
//
// Parameters:
//   - cfg: Application configuration containing database settings
//
// Returns:
//   - *connectionParams: Parsed and validated connection parameters
//   - error: Any error that occurred during parameter parsing
func loadConnectionParams(cfg *config.Config) (*connectionParams, error) {
	// Parse maximum connections setting
	maxConnections, err := strconv.ParseInt(cfg.DBMaxConn, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing max connections: %v", err)
	}

	// Parse minimum connections setting
	minConnections, err := strconv.ParseInt(cfg.DBMinConn, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing min connections: %v", err)
	}

	// Parse maximum connection lifetime
	maxConnectionLife, err := time.ParseDuration(cfg.DBConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("error parsing max connection lifetime: %v", err)
	}

	// Parse maximum idle time for connections
	maxIdleTime, err := time.ParseDuration(cfg.DBConnMaxIdleLifetime)
	if err != nil {
		return nil, fmt.Errorf("error parsing max idle time: %v", err)
	}

	// Parse health check period
	healthCheckPeriod, err := time.ParseDuration(cfg.DBHealthCheckPeriod)
	if err != nil {
		return nil, fmt.Errorf("error parsing health check period: %v", err)
	}

	// Parse connection timeout
	connectionTimeout, err := time.ParseDuration(cfg.ConnectTimeout)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection timeout: %v", err)
	}

	// Return the parsed connection parameters
	return &connectionParams{
		Host:              cfg.DBHost,
		Port:              cfg.DBPort,
		User:              cfg.DBUser,
		Password:          cfg.DBPassword,
		Name:              cfg.DBName,
		SSLMode:           cfg.SSLMode,
		MaxConnections:    int32(maxConnections),
		MinConnections:    int32(minConnections),
		MaxConnectionLife: maxConnectionLife,
		MaxIdleTime:       maxIdleTime,
		HealthCheckPeriod: healthCheckPeriod,
		ConnectionTimeout: connectionTimeout,
	}, nil
}
