package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timeStamp     int64
	transactions  []string
}

func NewBlock(nonce int,previousHash [32]byte) *Block {
	b := new(Block)
	b.timeStamp= time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash

	return b
}

func (b Block) Hash () [32]byte{
	m,_ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b Block) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        Nonce        int `json:"nonce"` 
        PreviousHash [32]byte `json:"previous_hash"` 
        Timestamp    int64  `json:"timestamp"` 
        Transactions []string `json:"transactions"` 
    }{
        Nonce:        b.nonce,                  
        PreviousHash: b.previousHash, 
        Transactions: b.transactions,
    })
}

func (bc *Block) Print() {
	fmt.Println("nonce	: ", bc.nonce)
	fmt.Printf("previousHash %x	:", bc.previousHash)
	fmt.Println("timeStamp	: ", bc.timeStamp)
	fmt.Println("nonce	:", bc.transactions)
}

type Blockchain struct {
	transactionPool [] string
	chain 			[] *Block
}

func NewBlockChain() *Blockchain {
	b:= &Block{}
	bc := new (Blockchain)
	bc.CreateBlock(0,b.Hash())

	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b  := NewBlock(nonce,previousHash)
	bc.chain =append(bc.chain, b)
	return b 
}

func (bc *Blockchain) LastBlock() *Block {	
	return bc.chain[len(bc.chain)-1]	 
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n",strings.Repeat("=",25),i,strings.Repeat("=",25))
		block.Print()
		fmt.Printf("%s\n",strings.Repeat("*",25))
	}
}

type Transaction struct {
	senderBlockChainAddress string
	recipientTransactions string
	value float32
}

func newTransaction(sender string,recipient string,value  float32)*Transaction {
	return &Transaction{
		sender,recipient,value,
	}
}

func (t *Transaction) Print() {	
		fmt.Printf("%s\n",strings.Repeat("-",40))
		fmt.Printf("sender_blockchain_address %s\n", t.senderBlockChainAddress)
		fmt.Printf("recipeient_blockchain_address %s\n", t.recipientTransactions)
		fmt.Printf("value %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct {
        SenderBlockChainAddress string `json:"sender_blockchain_address"` 
		RecipientTransactions string `json:"recipeient_blockchain_address"` 
		Value float32 `json:"value"` 
    }{
        SenderBlockChainAddress: t.senderBlockChainAddress,                  
        RecipientTransactions: t.recipientTransactions, 
        Value: t.value,
    })
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	chain:= NewBlockChain()
	chain.Print()


	chain.CreateBlock(5,chain.LastBlock().Hash())
	chain.Print()

	chain.CreateBlock(2,chain.LastBlock().Hash())
	chain.Print()
}