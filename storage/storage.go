package storage

const (
	// FileName 节点信息保存配置文件
	FileName = "config"
)

type Storage interface {
	ReadVoteInfo(key string) ([]byte, error)
	WriteVoteInfo(key string, bytes []byte) error
}
