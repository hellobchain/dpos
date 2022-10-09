package leveldb

import (
	"github.com/wsw365904/dpos/pkg/leveldb"
	"github.com/wsw365904/dpos/storage"
)

var _ storage.Storage = (*leveldbImpl)(nil)

type leveldbImpl struct {
	leveldb *leveldb.LevelDbEngine
}

func NewLeveldbImpl(leveldbPath string) (storage.Storage, error) {
	levelDbEngine, err := leveldb.NewLevelDbEngine(leveldbPath)
	if err != nil {
		return nil, err
	}
	return &leveldbImpl{
		leveldb: levelDbEngine,
	}, nil
}

func (l *leveldbImpl) ReadVoteInfo(key string) ([]byte, error) {
	return l.leveldb.Get("consensus.dpos." + key)
}

func (l *leveldbImpl) WriteVoteInfo(key string, bytes []byte) error {
	return l.leveldb.Put("consensus.dpos."+key, bytes)
}
