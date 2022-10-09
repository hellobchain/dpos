package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDbEngine struct {
	Db *leveldb.DB
}

func NewLevelDbEngine(path string) (*LevelDbEngine, error) {
	//创建并打开数据库
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDbEngine{
		Db: db,
	}, nil
}

func (e *LevelDbEngine) Close() error {
	return e.Db.Close()
}

func (e *LevelDbEngine) Get(key string) ([]byte, error) {
	res, err := e.Db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	if err == leveldb.ErrNotFound {
		return []byte("-1"), nil
	}
	return res, nil
}

func (e *LevelDbEngine) Put(key string, value []byte) error {
	return e.Db.Put([]byte(key), value, nil)
}

func (e *LevelDbEngine) Del(key string) error {
	return e.Db.Delete([]byte(key), nil)
}
