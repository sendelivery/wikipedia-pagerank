package config

import (
	"context"
	"log/slog"

	"github.com/sendelivery/wikipedia-pagerank/internal/logger"
)

type key int

const configKey key = iota

// Config holds the configuration for the application.
type Config struct {
	// NumConcurrentFetches is the number of concurrent fetches to make when scraping Wikipedia
	// articles.
	NumConcurrentFetches int

	// NumPages is the number of pages to scrape.
	NumPages int

	// OutputDir is the directory to write the output to.
	OutputDir string

	// Logger is the logger to use.
	Logger *slog.Logger

	// ConvergenceThreshold is the threshold for the PageRank algorithm to converge.
	ConvergenceThreshold float64

	// DampingFactor is the damping factor for the PageRank algorithm.
	DampingFactor float64
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		NumConcurrentFetches: 50,
		NumPages:             1000,
		OutputDir:            "output",
		Logger:               logger.MakeLogger("log"),
		ConvergenceThreshold: 0.0001,
		DampingFactor:        0.85,
	}
}

// ContextWithConfig returns a new context with the given configuration.
func ContextWithConfig(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

// ConfigFromContext returns the configuration from the given context.
func ConfigFromContext(ctx context.Context) (*Config, bool) {
	cfg, ok := ctx.Value(configKey).(Config)
	return &cfg, ok
}
