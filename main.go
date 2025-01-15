// main.go
package main

func main() {
	// 创建区块链
	bc := NewBlockchain()

	// 创建交易池
	tp := NewTransactionPool()

	// 模拟一些交易
	tp.AddTransaction(&Transaction{Sender: "Alice", Recipient: "Bob", Amount: 10})
	tp.AddTransaction(&Transaction{Sender: "Bob", Recipient: "Charlie", Amount: 5})

	// 打包交易到新区块
	tp.MineBlock(bc)

	// 打印区块链
	for _, block := range bc.Blocks {
		println("Block Index:", block.Index)
		println("Block Hash:", string(block.Hash))

		// 反序列化交易数据
		transactions, err := DeserializeTransactions(block.Data)
		if err != nil {
			println("Failed to deserialize transactions:", err.Error())
			continue
		}

		// 打印交易信息
		for _, tx := range transactions {
			println("Transaction:", tx.Sender, "->", tx.Recipient, "Amount:", tx.Amount)
		}
	}
}
