package blockchain

type Block struct {
	Hash     []byte // hash of this block - derived from the Data and prevHash
	PrevHash []byte // hash of the previous block
	Data     []byte
	Nonce     int
	// Timestamp int
}

type BlockChain struct {
	Blocks []*Block
}


func CreateBlock(data string, prevHash []byte) *Block {
	b := &Block{
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}
	pow := NewProofOfWork(b)
	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash[:]
	return b
}

func GenesisBlock() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	bc := &BlockChain{
		Blocks: []*Block{GenesisBlock()},
	}
	return bc
}

func (bc *BlockChain) AddBlock(data string) {
	prveBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(data, prveBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
