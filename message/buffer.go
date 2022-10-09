package message

import "sync"

type Buffer struct {
	voteInfoLock   sync.RWMutex
	VoteInfoBuffer *VoteInfos
	Blocks         *Blocks
	blockLock      sync.RWMutex
}

func NewBuffer() *Buffer {
	return &Buffer{
		VoteInfoBuffer: newVoteInfos(),
		Blocks:         newBlocks(),
	}
}

func (b *Buffer) BufferVoteInfo(voteInfo *VoteInfo) {
	b.voteInfoLock.Lock()
	defer b.voteInfoLock.Unlock()
	b.VoteInfoBuffer.VoteInfos = append(b.VoteInfoBuffer.VoteInfos, voteInfo)
}

func (b *Buffer) GetBufferVoteInfos() *VoteInfos {
	b.voteInfoLock.Lock()
	defer b.voteInfoLock.Unlock()
	return b.VoteInfoBuffer
}

func (b *Buffer) BufferBlock(block *Block) {
	b.blockLock.Lock()
	defer b.blockLock.Unlock()
	b.Blocks.Blocks = append(b.Blocks.Blocks, block)
}

func (b *Buffer) GetBufferBlocks() *Blocks {
	b.blockLock.Lock()
	defer b.blockLock.Unlock()
	return b.Blocks
}
