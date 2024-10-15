package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/amrshaban2005/go-blockchain/blockchain"
	"github.com/amrshaban2005/go-blockchain/utils"
	"github.com/amrshaban2005/go-blockchain/wallet"
)

const pathToTemplateDir = "templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(pathToTemplateDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Printf("Error: invalid http method")
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := json.Marshal(myWallet)
		io.Writer.Write(w, m[:])
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: invalid http method")
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		//io.Writer.Write(w, utils.JsonStatus("test"))
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("Error: %v", err)
			io.Writer.Write(w, utils.JsonStatus("fail"))
			return
		}
		if !t.Validate() {
			log.Printf("Error: miss field(s)")
			io.Writer.Write(w, utils.JsonStatus("fail"))
			return
		}

		publicKey := utils.String2PublicKey(*t.SenderPublicKey)
		privateKey := utils.String2PrivateKey(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("Error: parse error")
			io.Writer.Write(w, utils.JsonStatus("fail"))
			return
		}
		value32 := float32(value)

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &blockchain.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if resp.StatusCode == 201 {
			io.Writer.Write(w, utils.JsonStatus("success"))
			return
		}
		io.Writer.Write(w, utils.JsonStatus("fail"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: invalid http method")
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transaction", ws.CreateTransaction)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.Port())), nil))
}
