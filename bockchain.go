package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MIN_DIFFICULTY = 4
	MINING_SENDER = "THE BLOCKCHAIN"
	MINING_REWARD = 5.0
)

type Block struct {
	timeStamp     int64
	nonce        int
	previousHash [32]byte	
	transactions  [] *Transaction
}

func NewBlock(nonce int,previousHash [32]byte, transactions[]*Transaction ) *Block {
	b := new(Block)
	b.timeStamp= time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions

	return b
}

func (bc *Block) Print() {
	fmt.Println("nonce	: ", bc.nonce)
	fmt.Printf("previousHash %x	:", bc.previousHash)
	fmt.Println("timeStamp	: ", bc.timeStamp)
	
	for _, row := range bc.transactions {
		fmt.Printf("%s\n",strings.Repeat("-",40))
		fmt.Println("Recipient 	: ", row.recipientTransactions)
		fmt.Println("Sender	: ", row.senderBlockChainAddress)
		fmt.Println("Value	: ", row.value)
		fmt.Printf("%s\n",strings.Repeat("-",40))
	}
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
        Transactions [] *Transaction `json:"transactions"` 
    }{
        Nonce:        b.nonce,                  
        PreviousHash: b.previousHash, 
        Transactions: b.transactions,
    })
}

type Blockchain struct {
	transactionPool [] *Transaction
	chain 			[] *Block
	blockchainAddress string
}

func NewBlockChain(blockchainAddress string) *Blockchain {
	b:= &Block{}
	bc := new (Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0,b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte ) *Block {
	b  := NewBlock(nonce,previousHash, bc.transactionPool)
	bc.chain =append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
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

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender,recipient,value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([] *Transaction,0)
	for  _,t := range bc.transactionPool{
		transactions = append(transactions, NewTransaction(t.recipientTransactions,t.senderBlockChainAddress,t.value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeors := strings.Repeat("0",difficulty)
	guessBlock := Block{0,nonce,previousHash,transactions}
	hashStr := fmt.Sprintf("%x",guessBlock.Hash()) 
	return  hashStr[:difficulty] == zeors
}

func  (bc *Blockchain) ProofOfWork() int{
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce,previousHash,transactions, MIN_DIFFICULTY) {
		nonce += 1
	}

	return nonce
}

func  (bc *Blockchain) Mining() bool{
	bc.AddTransaction(MINING_SENDER,bc.blockchainAddress,MINING_REWARD)	
	bc.CreateBlock(bc.ProofOfWork(),bc.LastBlock().Hash())
	log.Println("action=mining, status=sucesss")
	return true
}

type Transaction struct {
	senderBlockChainAddress string
	recipientTransactions string
	value float32
}

func NewTransaction(sender string,recipient string,value  float32)*Transaction {
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
	chain:= NewBlockChain("my_address")

	chain.AddTransaction("B","A",100)	
	chain.Mining()
	chain.Print()

	chain.AddTransaction("C","D",34.5)
	chain.AddTransaction("G","H",67)
	chain.Mining()
	chain.Print()
}