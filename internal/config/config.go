package config

import (
	"os"
	"strconv"
)

type BrownNoiseConfig struct {
	Damping  float64
	Gain     float64
	StepSize float64
}

type ServerConfig struct {
	Port      string
	Host      string
	StaticDir string
}

func GetBrownNoiseConfig() BrownNoiseConfig {
	config := BrownNoiseConfig{
		Damping:  0.90,
		Gain:     0.5,
		StepSize: 0.02,
	}

	if damping := os.Getenv("BROWN_NOISE_DAMPING"); damping != "" {
		if val, err := strconv.ParseFloat(damping, 64); err == nil {
			config.Damping = val
		}
	}

	if gain := os.Getenv("BROWN_NOISE_GAIN"); gain != "" {
		if val, err := strconv.ParseFloat(gain, 64); err == nil {
			config.Gain = val
		}
	}

	if stepSize := os.Getenv("BROWN_NOISE_STEP_SIZE"); stepSize != "" {
		if val, err := strconv.ParseFloat(stepSize, 64); err == nil {
			config.StepSize = val
		}
	}

	return config
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		Port:      getEnvOrDefault("PORT", "8080"),
		Host:      getEnvOrDefault("HOST", "localhost"),
		StaticDir: getEnvOrDefault("STATIC_DIR", "internal/server/static"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
