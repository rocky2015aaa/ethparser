package models

// @Description A transaction record on the blockchain
// @Model
type Transaction struct {
	Hash      string `json:"hash" example:"0x009ed951afbef35139089ce03192a5a2d3226c397407c5f39a707c7f3d596bc8"`
	From      string `json:"from" example:"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"`
	To        string `json:"to" example:"0xe688b84b23f322a994a53dbf8e15fa82cdb71127"`
	Value     string `json:"value" example:"0xc047d21ca5f809"`
	BlockHash string `json:"blockHash" example:"0x59d7f3c1cf9ada06b52cd36efea2be5b29dd1e15649aaf92c45d305277ec6693"`
	BlockNum  int    `json:"blockNum" example:"2132563"`
	Timestamp string `json:"timestamp" example:"2024-11-29 09:42:12 +0900 KST"`
}
