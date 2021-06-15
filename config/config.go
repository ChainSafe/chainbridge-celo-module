package config

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains"
	"github.com/spf13/viper"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1
const DefaultBlockConfirmations = 10

type CeloConfig struct {
	GeneralChainConfig chains.GeneralChainConfig
	Bridge             string
	Erc20Handler       string
	Erc721Handler      string
	GenericHandler     string
	MaxGasPrice        *big.Int
	GasMultiplier      *big.Float
	GasLimit           *big.Int
	Http               bool
	StartBlock         *big.Int
	BlockConfirmations *big.Int
	// EgsApiKey          string // API key for ethgasstation to query gas prices
	// EgsSpeed           string // The speed which a transaction should be processed: average, fast, fastest. Default: fast
}

type RawCeloConfig struct {
	chains.GeneralChainConfig `mapstructure:",squash"`
	Bridge                    string  `mapstructure:"bridge"`
	Erc20Handler              string  `mapstructure:"erc20Handler"`
	Erc721Handler             string  `mapstructure:"erc721Handler"`
	GenericHandler            string  `mapstructure:"genericHandler"`
	MaxGasPrice               int64   `mapstructure:"maxGasPrice"`
	GasMultiplier             float64 `mapstructure:"gasMultiplier"`
	GasLimit                  int64   `mapstructure:"gasLimit"`
	Http                      bool    `mapstructure:"http"`
	StartBlock                int64   `mapstructure:"startBlock"`
	BlockConfirmations        int64   `mapstructure:"blockConfirmations"`
	// EgsApiKey                 string  `mapstructure:"egsApiKey"`
	// EgsSpeed                  string  `mapstructure:"egsSpeed"`
}

func GetConfig(path string, name string) (*CeloConfig, error) {
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

	cfg, err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfig(rawConfig *RawCeloConfig) (*CeloConfig, error) {

	config := &CeloConfig{
		GeneralChainConfig: rawConfig.GeneralChainConfig,
		Erc20Handler:       rawConfig.Erc20Handler,
		Erc721Handler:      rawConfig.Erc721Handler,
		GenericHandler:     rawConfig.GenericHandler,
		GasLimit:           big.NewInt(DefaultGasLimit),
		MaxGasPrice:        big.NewInt(DefaultGasPrice),
		GasMultiplier:      big.NewFloat(DefaultGasMultiplier),
		Http:               rawConfig.Http,
		StartBlock:         big.NewInt(rawConfig.StartBlock),
		BlockConfirmations: big.NewInt(DefaultBlockConfirmations),
	}

	if rawConfig.Bridge != "" {
		config.Bridge = rawConfig.Bridge
	} else {
		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
	}

	if rawConfig.GasLimit != 0 {
		config.GasLimit = big.NewInt(rawConfig.GasLimit)
	}

	if rawConfig.MaxGasPrice != 0 {
		config.MaxGasPrice = big.NewInt(rawConfig.MaxGasPrice)
	}

	if rawConfig.GasMultiplier != 0 {
		config.GasMultiplier = big.NewFloat(rawConfig.GasMultiplier)
	}

	if rawConfig.BlockConfirmations != 0 {
		config.BlockConfirmations = big.NewInt(rawConfig.BlockConfirmations)
	}

	return config, nil
}
