package bitcoin

type AmountSat = uint64
type Vin struct {
	Hash    string    `json:"hash"`
	Index   uint32    `json:"index"`
	Amount  AmountSat `json:"amount"`
	Address string    `json:"address"`
}

type Vout struct {
	Address string    `json:"address"`
	Amount  AmountSat `json:"amount"`
	Index   uint32    `json:"index"`
}

type BitcoinSchema struct {
	RequestId string  `json:"request_id"`
	Fee       string  `json:"fee"`
	Vins      []*Vin  `json:"vin"`
	Vouts     []*Vout `json:"vouts"`
}
