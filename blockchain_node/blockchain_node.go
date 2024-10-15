package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/amrshaban2005/go-blockchain/blockchain"
	block "github.com/amrshaban2005/go-blockchain/blockchain"
	"github.com/amrshaban2005/go-blockchain/utils"
	"github.com/amrshaban2005/go-blockchain/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain, 0)

type BlockchainNode struct {
	port uint16
}

func NewBlockchianNode(port uint16) *BlockchainNode {
	return &BlockchainNode{port}
}

func (bcn *BlockchainNode) Port() uint16 {
	return bcn.port
}

func (bcn *BlockchainNode) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcn.Port())
		cache["blockchain"] = bc
	}
	return bc

}

func (bcn *BlockchainNode) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		m, _ := json.Marshal(bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: invalid http method")
	}
}

func (bcn *BlockchainNode) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transaction []*block.Transaction `json:"transactions"`
			Length      int                  `json:"length"`
		}{
			Transaction: transactions,
			Length:      len(transactions),
		})
		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("Error: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("Error: missing filed(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		publicKey := utils.String2PublicKey(*t.SenderPublicKey)
		signature := utils.String2Signature(*t.Signature)
		bc := bcn.GetBlockchain()
		isCreated := bc.CreateTransactions(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)
		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}

		io.WriteString(w, string(m))
		return

	default:
		log.Printf("Error: invalid http method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcn *BlockchainNode) Mine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcn.GetBlockchain()
		isMined := bc.Mining()
		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.Writer.Write(w, utils.JsonStatus(string(m)))

	default:
		log.Println("Error: invalid http method")
		io.Writer.Write(w, utils.JsonStatus("fail"))

	}
}

func (bcn *BlockchainNode) StartMine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		bc.StartMining()
		io.Writer.Write(w, utils.JsonStatus("success"))

	default:
		log.Println("Error: invalid http method")
		io.Writer.Write(w, utils.JsonStatus("fail"))
	}
}

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", bcn.GetChain)
	http.HandleFunc("/transactions", bcn.Transactions)
	http.HandleFunc("/mine", bcn.Mine)
	http.HandleFunc("/mine/start", bcn.StartMine)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.Port())), nil))
}
