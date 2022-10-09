package message

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type VoteInfo struct {
	NodeName string `json:"nodeName"`
	VoteNum  uint64 `json:"voteNum"`
}

type VoteInfos struct {
	VoteInfos []*VoteInfo `json:"voteInfos"`
}

func newVoteInfos() *VoteInfos {
	return &VoteInfos{
		VoteInfos: make([]*VoteInfo, 0),
	}
}

// Block struct, A block contain 以下信息:
// Index 索引、Timestamp(时间戳)、BPM、Hash(自己的hash值)、PreHash(上一个块的Hash值)、validator(此区块的生产者信息)
type Block struct {
	Index     int    `json:"index"`
	Timestamp string `json:"timestamp"`
	BPM       int    `json:"BPM"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash"`
	Validator string `json:"validator"`
}

type Blocks struct {
	Blocks []*Block `json:"blocks"`
}

func newBlocks() *Blocks {
	return &Blocks{
		Blocks: make([]*Block, 0),
	}
}

// CalculateHash 计算string的hash值
func CalculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// CalculateBlockHash 计算Block的hash值
func CalculateBlockHash(block *Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	return CalculateHash(record)
}

// GenerateBlock 根据上一个区块信息，生成新的区块
func GenerateBlock(oldBlock *Block, BPM int, address string) (*Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.BPM = BPM
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateBlockHash(&newBlock)
	newBlock.Validator = address

	return &newBlock, nil
}

// IsBlockValid 校验区块是否合法
func IsBlockValid(newBlock, oldBlock *Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}
