package config

import (
	"encoding/json"
	"os"

	head "github.com/mbanq/iso20022-go/ISO20022/head_001_001_02"
	pacs "github.com/mbanq/iso20022-go/ISO20022/pacs_008_001_08"
)

type Config struct {
	MarketPractice         head.Max2048Text                               `json:"marketPractice"`
	MarketPracticeRegistry head.Max350Text                                `json:"marketPracticeRegistry"`
	LocalInstrument        pacs.LocalInstrument2Choice                    `json:"localInstrument"`
	SettlementMethod       pacs.SettlementMethod1Code                     `json:"settlementMethod"`
	ClearingSystemId       pacs.ExternalCashClearingSystem1Code           `json:"clearingSystemId"`
	ChargeBearer           pacs.ChargeBearerType1Code                     `json:"chargeBearer"`
	Currency               pacs.ActiveCurrencyCode                        `json:"currency"`
	ClearingSystem         head.ExternalClearingSystemIdentification1Code `json:"clearingSystem"`
	FrbId                  head.Max35Text                                 `json:"frbId"`
	IspId                  head.Max35Text                                 `json:"ispId"`
}

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
