package model

type LatestBlock struct {
	BaseRND       string        `json:"baseRND"`
	Hash          string        `json:"hash"`
	Number        string        `json:"number"`
	PrevBlockHash string        `json:"prevBlockHash"`
	Signatures    []interface{} `json:"signatures"`
	TxHashes      []struct {
		ExecHash string `json:"execHash"`
		Hash     string `json:"hash"`
	} `json:"txHashes"`
}

type TxResponse struct {
	Transactions []Tx `json:"transactions"`
}

type Tx struct {
	Hash      string `json:"hash"`
	Block     string `json:"block"`
	Nonce     string `json:"nonce"`
	Vm        string `json:"vm"`
	Sender    string `json:"sender"`
	Signature string `json:"signature"`
	Message   string `json:"message"`
	Exec      struct {
		Hash       string     `json:"hash"`
		VmResponse VmResponse `json:"vmResponse"`
	} `json:"exec"`
	FeeCurrency string `json:"feeCurrency"`
}

type VmResponse struct {
	C interface{}            `json:"C"`
	D map[string]interface{} `json:"D"`
	R map[string]interface{} `json:"R"`
	T map[string]interface{} `json:"T"`
	V map[string]interface{} `json:"V"`
}

type Message struct {
}

type Amount struct {
}

type Json map[string]interface{}
