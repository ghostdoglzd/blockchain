package blockchain

import (
	"encoding/json"
)

// 定义交易结构
type Transaction struct {
	Sender    string  // 发送者
	Recipient string  // 接收者
	Amount    float64 // 交易金额
}

// 定义交易池
type TransactionPool struct {
	Transactions []*Transaction // 待处理的交易列表
}

// 创建新的交易池
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Transactions: []*Transaction{},
	}
}

// 添加交易到交易池
func (tp *TransactionPool) AddTransaction(tx *Transaction) {
	tp.Transactions = append(tp.Transactions, tx)
}

// 清空交易池
func (tp *TransactionPool) Clear() {
	tp.Transactions = []*Transaction{}
}

// 将交易列表序列化为 []byte
func SerializeTransactions(transactions []*Transaction) ([]byte, error) {
	return json.Marshal(transactions)
}

// 打包交易到新区块
func (tp *TransactionPool) MineBlock(bc *Blockchain) {
	if len(tp.Transactions) == 0 {
		println("No transactions to mine.")
		return
	}

	// 将交易列表序列化为 []byte
	data, err := SerializeTransactions(tp.Transactions)
	if err != nil {
		println("Failed to serialize transactions:", err.Error())
		return
	}

	// 获取上一个区块
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// 使用原来的 NewBlock 函数创建新区块
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)

	// 将新区块添加到区块链
	bc.Blocks = append(bc.Blocks, newBlock)

	// 清空交易池
	tp.Clear()

	println("New block mined with", len(tp.Transactions), "transactions.")
}

// 反序列化交易数据
func DeserializeTransactions(data []byte) ([]*Transaction, error) {
	var transactions []*Transaction
	err := json.Unmarshal(data, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
