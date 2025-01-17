package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (bc *Blockchain) StartHTTPServer(port string, tp *TransactionPool) {
	// 查询所有区块
	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bc.Blocks)
	})

	// 创建交易
	http.HandleFunc("/createTransaction", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var tx Transaction
		err := json.NewDecoder(r.Body).Decode(&tx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tp.AddTransaction(&tx)
		w.WriteHeader(http.StatusCreated)

		tp.MineBlock(bc)
		fmt.Fprintf(w, "Transaction added successfully")
	})

	// 查询特定交易
	http.HandleFunc("/getTransaction", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		sender := r.URL.Query().Get("sender")
		if sender == "" {
			http.Error(w, "Missing sender parameter", http.StatusBadRequest)
			return
		}

		var transactions []*Transaction
		for _, block := range bc.Blocks {
			txs, err := DeserializeTransactions(block.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, tx := range txs {
				if tx.Sender == sender {
					transactions = append(transactions, tx)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	})

	// 查询所有交易
	http.HandleFunc("/allTransactions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var allTransactions []*Transaction
		for _, block := range bc.Blocks {
			txs, err := DeserializeTransactions(block.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			allTransactions = append(allTransactions, txs...)
		}

		json.NewEncoder(w).Encode(allTransactions)
	})

	fmt.Println("Starting HTTP server on port", port)
	http.ListenAndServe(":"+port, nil)
}
