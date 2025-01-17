package blockchain

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

const (
	Protocol      = "tcp"
	CommandLength = 12
)

type Node struct {
	Address string
	Conn    net.Conn
}

type Network struct {
	Nodes []*Node
	mu    sync.Mutex
}

func NewNetwork() *Network {
	return &Network{
		Nodes: []*Node{},
	}
}

func (n *Network) AddNode(address string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	conn, err := net.Dial(Protocol, address)
	if err != nil {
		fmt.Println("Failed to connect to node:", err)
		return
	}

	node := &Node{
		Address: address,
		Conn:    conn,
	}
	n.Nodes = append(n.Nodes, node)
}

func (n *Network) BroadcastBlock(block *Block) {
	n.mu.Lock()
	defer n.mu.Unlock()

	data, err := json.Marshal(block)
	if err != nil {
		fmt.Println("Failed to marshal block:", err)
		return
	}

	for _, node := range n.Nodes {
		_, err := node.Conn.Write(data)
		if err != nil {
			fmt.Println("Failed to send block to node:", err)
		}
	}
}

func (n *Network) StartServer(port string, bc *Blockchain) {
	listener, err := net.Listen(Protocol, ":"+port)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go n.handleConnection(conn, bc)
	}
}

func (n *Network) handleConnection(conn net.Conn, bc *Blockchain) {
	defer conn.Close()

	var block Block
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Failed to decode block:", err)
		return
	}

	bc.AddBlock(block.Data)
	fmt.Println("Received new block and added to the blockchain")
}
