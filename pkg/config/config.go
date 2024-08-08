package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
)

// Predefined errors for the package.
var (
	ErrConfigNotInitialized = errors.New("config not initialized")
)

// Config holds the application configuration settings.
type Config struct {
	Test string `json:"test"`
}

// Global variable to hold the configuration in memory.
var config *Config

// Default configuration values.
var defaultConfig = &Config{
	Test: "123",
}

// Init initializes the configuration by either creating a new config file
// or loading an existing one.
func Initialize() {
	if !doesConfigFileExist() {
		if err := createConfigFile(); err != nil {
			panic(err)
		}
	}
	config = defaultConfig
}

// GetConfig returns the current configuration. It loads the configuration
// from the config file if it has not been initialized yet.
func GetConfig() (*Config, error) {
	if config == nil {
		return nil, ErrConfigNotInitialized
	}

	cfg, err := getDataFromConfigFile()
	if err != nil {
		return nil, err
	}

	config = cfg
	return config, nil
}

// WriteConfig saves the provided configuration to the config file.
func WriteConfig(c *Config) error {
	config = c

	f, err := os.Create(getConfigFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

// getConfigFilePath returns the path to the configuration file.
func getConfigFilePath() string {
	return "config.json"
}

// createConfigFile creates a new config file with default settings.
func createConfigFile() error {
	f, err := os.Create(getConfigFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

// doesConfigFileExist checks if the configuration file exists.
func doesConfigFileExist() bool {
	_, err := os.Stat(getConfigFilePath())
	return !os.IsNotExist(err)
}

// getDataFromConfigFile loads and returns the configuration from the config file.
// If the file is empty, it writes the default configuration to the file.
func getDataFromConfigFile() (*Config, error) {
	f, err := os.Open(getConfigFilePath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Check if the file is empty and write default config if necessary.
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		err = WriteConfig(defaultConfig)
		if err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	var rawConfig map[string]interface{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&rawConfig); err != nil {
		return nil, err
	}

	// Ensure all expected keys are present in the configuration.
	if updated := ensureAllKeysExist(rawConfig); updated {
		if err := writeUpdatedConfigFile(rawConfig); err != nil {
			return nil, err
		}
	}

	config := &Config{}
	if err := mapToStruct(rawConfig, config); err != nil {
		return nil, err
	}

	return config, nil
}

// ensureAllKeysExist checks if all keys in the default configuration are
// present in the provided rawConfig. If any keys are missing, it adds them
// with their zero values.
func ensureAllKeysExist(rawConfig map[string]interface{}) bool {
	updated := false
	expectedKeys := getExpectedKeys()
	for _, key := range expectedKeys {
		if _, ok := rawConfig[key]; !ok {
			rawConfig[key] = getZeroValueForKey(key)
			updated = true
		}
	}
	return updated
}

// writeUpdatedConfigFile writes the updated rawConfig map to the config file.
func writeUpdatedConfigFile(rawConfig map[string]interface{}) error {
	f, err := os.Create(getConfigFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(rawConfig, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	return err
}

// getExpectedKeys returns a list of all expected keys (JSON tags) in the Config struct.
func getExpectedKeys() []string {
	var keys []string
	val := reflect.TypeOf(Config{})
	for i := 0; i < val.NumField(); i++ {
		keys = append(keys, strings.ToLower(val.Field(i).Tag.Get("json")))
	}
	return keys
}

// getZeroValueForKey returns the zero value for the given key in the Config struct.
func getZeroValueForKey(key string) interface{} {
	val := reflect.TypeOf(Config{})
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if strings.ToLower(field.Tag.Get("json")) == key {
			return reflect.Zero(field.Type).Interface()
		}
	}
	return nil
}

// mapToStruct converts a map to a Config struct using JSON marshaling.
func mapToStruct(m map[string]interface{}, s *Config) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}
