package ethereum

type Eip1559DynamicFeeTx struct {
	ChainId              string `json:"chain_id"`                 // 链ID (如 Ethereum 主网=1, Goerli=5)
	Nonce                uint64 `json:"nonce"`                    // 发送者的交易序号
	FromAddress          string `json:"from_address"`             // 发送方地址
	ToAddress            string `json:"to_address"`               // 接收方地址
	GasLimit             uint64 `json:"gas_limit"`                // Gas 上限
	MaxFeePerGas         string `json:"max_fee_per_gas"`          // 用户愿意付的每单位 Gas 的最高费用（wei）
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"` // 给矿工的小费（tip）
	Amount               string `json:"amount"`                   // 转账金额（wei，建议用字符串避免精度丢失）
	ContractAddress      string `json:"contract_address"`         // 合约地址（如果是合约调用）
}

type LegacyFeeTx struct {
	ChainId         string `json:"chain_id"`
	Nonce           uint64 `json:"nonce"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	GasLimit        uint64 `json:"gas_limit"`
	GasPrice        uint64 `json:"gas_price"`
	Amount          string `json:"amount"`
	ContractAddress string `json:"contract_address"`
}

type EthereumSchema struct {
	RequestId    string              `json:"request_id"`
	DynamicFeeTx Eip1559DynamicFeeTx `json:"dynamic_fee_tx"`
	ClassicFeeTx LegacyFeeTx         `json:"classic_fee_tx"`
}
