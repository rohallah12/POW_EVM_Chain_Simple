package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
)

// Balance map
type State struct {
	Balances map[int64]int64
}

const (
	diff = 2 << 10
)

type Block struct {
	Hash      []byte        //previous hash
	PrHash    []byte        //previous hash
	Txs       []Transaction //transactions hash
	StateRoot []byte
	Nonce     int64
}

type Transaction struct {
	Data []byte
}

type Blockchain struct {
	CurrentBlock Block
	CurrentState State
	Blocks       []Block
}

// state_1 : initial state in which evm is performing the transition
// data : machine code to run
// returns a new state, state root and a potential error
//how our EVM works?
/*
	we have 4 basic opcodes : ADD , SUB, SSTORE, PUSH
	ADD(input) : add last number in stack with input
	SUB(input) : sub last number in stack from input
	SSTORE(s) : store stack in storage location s
	PUSH : push data into stack
	SLOAD : load data from storage
*/
const (
	ADD    = 0x1
	SUB    = 0x2
	SSTORE = 0x3
	SLOAD  = 0x4
	PUSH   = 0x5
)

// very very very simple virutal machine
func evm(state_1 State, data []byte) (State, []byte, error) {
	//create stack which is our temporary memory
	state_2 := State{
		Balances: make(map[int64]int64),
	}
	//initialize 1000 slots for secondary state
	for i := 0; i < 1000; i++ {
		state_2.Balances[int64(i)] = state_1.Balances[int64(i)]
	}
	var stack int64
	for index, opcode := range data {
		switch opcode {
		case PUSH:
			arg := int64(data[index+1]) //extract next byte
			stack = arg
		case ADD:
			arg := int64(data[index+1]) //extract next 8 bytes
			result := stack + arg
			stack = result
		case SUB:
			arg := int64(data[index+1]) //extract next 8 bytes
			result := stack - arg
			stack = result
		case SSTORE:
			arg := int64(data[index+1]) //extract next 8 bytes
			state_2.Balances[arg] = stack
		case SLOAD:
			arg := int64(data[index+1]) //extract next 8 bytes
			stack = state_1.Balances[arg]
		}
		index += 1
	}
	stateBytes := encodeState(state_2)
	stateRoot := keccak256(stateBytes)
	return state_2, stateRoot, nil
}

func keccak256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// state to byte
func encodeState(state State) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(state)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// byte to int
func byteToInt(data []byte) int64 {
	var num int64
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &num)
	if err != nil {
		panic(err)
	}
	return num
}

func int64ToBytes(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		log.Fatal("Failed to convert int64 to bytes: ", err)
	}
	return buf.Bytes()
}

// block to byte
func encodeBlock(block Block) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(block)
	if err != nil {
		panic("error")
	}
	return buf.Bytes()
}

// creating a transaction
func CreateTransferTx(from int64, to int64, amount int64) []byte {
	// opcode for transferring tokens
	machineCode := []byte{
		SLOAD, byte(from), SUB, byte(50), SSTORE, byte(from),
		SLOAD, byte(to), ADD, byte(50), SSTORE, byte(to),
	}
	return machineCode
}

func (block *Block) AddTxToBlock(tx []byte) {
	//we should verify Tx first, but we dont do taht for simplicity
	block.Txs = append(block.Txs, Transaction{
		Data: tx,
	})
}

// construct a transaction for rewarding a miner, if is genesis, reward with 200 tokens
func GetMinerRewardOpcodes(miner int64, isGenesis bool) []byte {
	var reward int64
	if isGenesis {
		reward = 200
	} else {
		reward = 10
	}
	return []byte{
		SLOAD, byte(miner), ADD, byte(reward), SSTORE, byte(miner),
	}
}

// we mine one block with each transaction for simplicity
// we give miner 10 tokens for mining a block
// opcode for mining block is
// SLOAD {miner} ADD reward SSTORE {miner}
func (chain *Blockchain) MineBlock(block *Block, miner int64, isGenesis bool) {
	//get previous block
	prBlock := chain.CurrentBlock
	//construct the block
	newBlock := Block{
		PrHash: prBlock.Hash,
		Txs:    block.Txs,
	}
	//get EVM opcodes for giving 10 tokens to miner
	mineTx := GetMinerRewardOpcodes(miner, isGenesis)
	newBlock.Txs = append(newBlock.Txs, Transaction{
		Data: mineTx,
	})
	//execute transactions using our cute little EVM :)
	state := chain.CurrentState
	var root []byte
	var err error
	for _, tx := range newBlock.Txs {
		state, root, err = evm(state, tx.Data)
		if err != nil {
			panic(err)
		}
	}
	//attatch state root
	newBlock.StateRoot = root
	//proof of work
	nonce, hash := Work(newBlock)
	//attach nonce and hash
	newBlock.Hash = hash
	newBlock.Nonce = nonce
	//attach block to chain
	chain.Blocks = append(chain.Blocks, newBlock)
	//set state to final state
	chain.CurrentState = state
	//set block to latest block
	chain.CurrentBlock = newBlock
}

// proof of work, calculate hash for differente nonce untill hash is under difficulty
func Work(block Block) (int64, []byte) {
	var result int64 = diff << 12
	var currentNonce int64
	finalHash := []byte{}
	result = 0
	for result < diff {
		block.Nonce = currentNonce
		blockBytes := encodeBlock(block)
		finalHash = keccak256(blockBytes)
		result = byteToInt(finalHash)
		currentNonce += 1
	}
	return currentNonce, finalHash
}

// verify that hash of a given block is proof of work
func verify(block Block, nonce int64) bool {
	targetNonce, targetHash := Work(block)
	if targetNonce != nonce {
		return false
	}
	if !bytes.Equal(targetHash, block.Hash) {
		return false
	}
	if byteToInt(targetHash) > diff {
		return false
	}
	return true
}

// Printing blocks
func (chain Blockchain) PrintChain() {
	fmt.Printf("\n=========Printing chain==========\n")
	for number, block := range chain.Blocks {
		fmt.Printf("Block Height %d \n", number)
		fmt.Printf("Block Previous hash %x\n", block.PrHash)
		fmt.Printf("Hash %x \n", block.Hash)
		fmt.Printf("===================\n")
	}
}

func main() {
	//Initalizing our Storage
	state1 := State{
		Balances: make(map[int64]int64),
	}
	//initialize 1000 slots
	for i := 0; i < 1000; i++ {
		state1.Balances[int64(i)] = int64(0)
	}

	//initialize our chain
	chain := Blockchain{
		Blocks:       []Block{},
		CurrentBlock: Block{},
		CurrentState: state1,
	}

	//Genesis Block
	genesisBlock := Block{}

	//mine first block and receive genesis tokens (200 tokens, we dont want to exceed 1 byte)
	SatoshiNakomotoAddress := int64(88)

	chain.MineBlock(&genesisBlock, SatoshiNakomotoAddress, true)
	fmt.Printf("Satoshi Balane after Genesis %d\n", chain.CurrentState.Balances[SatoshiNakomotoAddress])

	// We want to transfer 20 tokens from Satoshi to alice
	aliceAddr := int64(33)
	tx := CreateTransferTx(SatoshiNakomotoAddress, aliceAddr, 20)

	//this time miner is bob, who will receive 10 tokens as reward
	bobTheMiner := int64(55)

	//Add transactions to our block
	newBlock := Block{}
	newBlock.AddTxToBlock(tx)

	//mine the block
	chain.MineBlock(&newBlock, bobTheMiner, false)

	//printing chain
	chain.PrintChain()

	//printing balances
	fmt.Printf("Satoshi Balane after sending to alice %d\n", chain.CurrentState.Balances[SatoshiNakomotoAddress])
	fmt.Printf("Alice Balane after receiving %d\n", chain.CurrentState.Balances[aliceAddr])
	fmt.Printf("BOB Balane after mining %d\n", chain.CurrentState.Balances[bobTheMiner])
}
