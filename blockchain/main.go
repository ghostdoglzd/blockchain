package main

import (
	"blockchain/blockchain/blockchain"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 创建区块链
	bc := blockchain.NewBlockchain()

	// 创建交易池
	tp := blockchain.NewTransactionPool()

	// 创建网络
	network := blockchain.NewNetwork()

	// 启动服务器
	go network.StartServer("8080", bc)

	// 启动 HTTP 服务器
	go bc.StartHTTPServer("8081", tp)

	// 提供静态文件服务（HTML 页面）
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// 启动静态文件服务器
	go func() {
		fmt.Println("Starting static file server on port 8082")
		if err := http.ListenAndServe(":8082", nil); err != nil {
			fmt.Println("Failed to start static file server:", err)
		}
	}()

	// 添加其他节点
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			network.AddNode(arg)
		}
	}

	// 模拟一些交易
	tp.AddTransaction(&blockchain.Transaction{Sender: "Alice", Recipient: "Bob", Amount: 10})
	tp.AddTransaction(&blockchain.Transaction{Sender: "Bob", Recipient: "Charlie", Amount: 5})

	// 打包交易到新区块
	tp.MineBlock(bc)

	// 广播新区块
	network.BroadcastBlock(bc.Blocks[len(bc.Blocks)-1])

	// 打印区块链
	for _, block := range bc.Blocks {
		fmt.Println("Block Index:", block.Index)
		fmt.Println("Block Hash:", string(block.Hash))

		// 反序列化交易数据
		transactions, err := blockchain.DeserializeTransactions(block.Data)
		if err != nil {
			fmt.Println("Failed to deserialize transactions:", err.Error())
			continue
		}

		// 打印交易信息
		for _, tx := range transactions {
			fmt.Println("Transaction:", tx.Sender, "->", tx.Recipient, "Amount:", tx.Amount)
		}
	}

	// 防止主程序退出
	go func() {
		fmt.Println("Starting static file server on port 8082")
		if err := http.ListenAndServe(":8082", nil); err != nil {
			fmt.Println("Failed to start static file server:", err)
		}
	}()
	select {}
}
