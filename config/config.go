package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ChainSafe/chainbridge-core/config"
	"github.com/spf13/viper"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1
const DefaultBlockConfirmations = 10

type CeloConfig struct {
	SharedEVMConfig config.SharedEVMConfig
}

type RawCeloConfig struct {
	config.RawSharedEVMConfig `mapstructure:",squash"`
}

func GetConfig(path string, name string) (*RawCeloConfig, error) {
	config := &RawCeloConfig{}

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in the config file, error: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into struct, error: %w", err)
	}

	if err := config.GeneralChainConfig.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseConfig(rawConfig *RawCeloConfig) (*CeloConfig, error) {

	cfg, err := rawConfig.RawSharedEVMConfig.ParseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse shared evm config, error: %w", err)
	}

	config := &CeloConfig{
		SharedEVMConfig: *cfg,
	}

	return config, nil
}

func (c *RawCeloConfig) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(&c); err != nil {
		fmt.Println("error marshalling json", "err", err)
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		fmt.Println("error creating config file", "err", err)
	}
	_, err = newFile.Write(raw)
	if err != nil {
		fmt.Println("error writing to config file", "err", err)
	}

	if err := newFile.Close(); err != nil {
		fmt.Println("failed to unmarshal config into struct", "err", err)
	}
	return newFile
}
