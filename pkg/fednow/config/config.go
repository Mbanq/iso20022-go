package config

import (
	"encoding/json"
	"os"

	"github.com/Mbanq/iso20022-go/ISO20022/head_001_001_02"
	"github.com/Mbanq/iso20022-go/ISO20022/pacs_008_001_08"
)

type Config struct {
	MarketPractice         head_001_001_02.Max2048Text                               `json:"marketPractice"`
	MarketPracticeRegistry head_001_001_02.Max350Text                                `json:"marketPracticeRegistry"`
	LocalInstrument        pacs_008_001_08.LocalInstrument2Choice                    `json:"localInstrument"`
	SettlementMethod       pacs_008_001_08.SettlementMethod1Code                     `json:"settlementMethod"`
	ClearingSystemId       pacs_008_001_08.ExternalCashClearingSystem1Code           `json:"clearingSystemId"`
	ChargeBearer           pacs_008_001_08.ChargeBearerType1Code                     `json:"chargeBearer"`
	Currency               pacs_008_001_08.ActiveCurrencyCode                        `json:"currency"`
	ClearingSystem         head_001_001_02.ExternalClearingSystemIdentification1Code `json:"clearingSystem"`
	FrbId                  head_001_001_02.Max35Text                                 `json:"frbId"`
	IspId                  head_001_001_02.Max35Text                                 `json:"ispId"`
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
