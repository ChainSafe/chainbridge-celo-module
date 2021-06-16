package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

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
}

type RawCeloConfig struct {
	chains.SharedEVMConfig `mapstructure:",squash"`
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
