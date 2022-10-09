package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/wsw365904/dpos/consensus"
	"github.com/wsw365904/dpos/message"
	"github.com/wsw365904/dpos/storage"
	"github.com/wsw365904/wswlog/wlogging"
	"io"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli"
)

const (
	// DefaultVote 节点默认的票数
	DefaultVote = 10
)

var logger = wlogging.MustGetLoggerWithoutName()

type P2p struct {
	mutex      sync.Mutex
	blockChain []*message.Block
	buffer     *message.Buffer
}

func NewP2p() *P2p {
	return &P2p{
		buffer:     message.NewBuffer(),
		blockChain: make([]*message.Block, 0),
	}
}

const (
	newCmd           = "new"
	portFlagForNew   = "port"
	targetFlagForNew = "target"
	secioFlagForNew  = "secio"
	seedFlagForNew   = "seed"
)

func CreateNewNodeCmd(p *P2p) cli.Command {
	// NewNode 创建新的节点加入到P2P网络
	nodeCMd := cli.Command{
		Name:  newCmd,
		Usage: "add a new node to p2p network",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  portFlagForNew,
				Value: 3000,
				Usage: "port of new node",
			},
			cli.StringFlag{
				Name:  targetFlagForNew,
				Value: "",
				Usage: "target node",
			},
			cli.BoolFlag{
				Name:  secioFlagForNew,
				Usage: "enable secio",
			},
			cli.Int64Flag{
				Name:  seedFlagForNew,
				Value: 0,
				Usage: "generate random",
			},
		},
		Action: func(context *cli.Context) error {
			if err := p.Run(context); err != nil {
				return err
			}
			return nil
		},
	}
	return nodeCMd
}

// MakeBasicHost 构建P2P网络
func MakeBasicHost(listenPort int, secio bool, randSeed int64) (host.Host, error) {
	var r io.Reader

	if randSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randSeed))
	}

	// 生产一对公私钥
	pri, _, err := crypto.GenerateKeyPairWithReader(crypto.ECDSA, 0, r)
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(pri),
	}

	if !secio {
		opts = append(opts, libp2p.NoSecurity)
	}
	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multi address
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiAddress to reach this host
	// by encapsulating both addresses;
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	logger.Infof("我是: %s\n", fullAddr)
	SavePeer(basicHost.ID().Pretty())

	if secio {
		logger.Infof("现在在一个新终端运行命令: './dpos new --port %d --target %s -secio' \n", listenPort+1, fullAddr)
	} else {
		logger.Infof("现在在一个新的终端运行命令: './dpos new --port %d --target %s' \n", listenPort+1, fullAddr)
	}
	return basicHost, nil
}

// HandleStream  handler stream info
func (p *P2p) HandleStream(s network.Stream) {
	logger.Infof("receive new connect: %s", s.Conn().RemotePeer().Pretty())
	// 将连接加入到
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go p.readData(rw)
	go p.writeData(rw)
}

// readData 读取数据输出到客户端
func (p *P2p) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logger.Errorf(err.Error())
		}

		if str == "" {
			return
		}
		if str != "\n" {
			chain := make([]*message.Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				logger.Errorf(err.Error())
			}
			p.mutex.Lock()
			if len(chain) > len(p.blockChain) {
				p.blockChain = chain
				bytes, err := json.MarshalIndent(p.blockChain, "", " ")
				if err != nil {
					logger.Errorf(err.Error())
				}

				logger.Infof("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			p.mutex.Unlock()
		}
	}
}

// writeData 将客户端数据处理写入BlockChain
func (p *P2p) writeData(rw *bufio.ReadWriter) {
	// 启动一个协程处理终端同步
	go func() {
		timer := time.NewTimer(2 * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				p.mutex.Lock()
				bytes, err := json.Marshal(p.blockChain)
				if err != nil {
					logger.Errorf(err.Error())
				}
				p.mutex.Unlock()
				p.mutex.Lock()
				_, err = rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
				if err != nil {
					logger.Errorf(err.Error())
				}
				err = rw.Flush()
				if err != nil {
					logger.Errorf(err.Error())
				}
				p.mutex.Unlock()
				timer.Reset(2 * time.Second)
			}
		}
	}()

	stdReader := bufio.NewReader(os.Stdin)
	for {
		logger.Info(">")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			logger.Errorf(err.Error())
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		bpm, err := strconv.Atoi(sendData)
		if err != nil {
			logger.Errorf(err.Error())
		}

		// pick选择block生产者
		address := dpos.PickWinner()
		logger.Infof("******节点 %s 获得了记账权利******", address)
		lastBlock := p.blockChain[len(p.blockChain)-1]
		newBlock, err := message.GenerateBlock(lastBlock, bpm, address)
		if err != nil {
			logger.Errorf(err.Error())
		}

		if message.IsBlockValid(newBlock, lastBlock) {
			p.mutex.Lock()
			p.blockChain = append(p.blockChain, newBlock)
			p.mutex.Unlock()
		}

		spew.Dump(p.blockChain)

		bytes, err := json.Marshal(p.blockChain)
		if err != nil {
			logger.Errorf(err.Error())
		}
		p.mutex.Lock()
		_, err = rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		if err != nil {
			return
		}
		err = rw.Flush()
		if err != nil {
			return
		}
		p.mutex.Unlock()
	}
}

// Run 函数
func (p *P2p) Run(ctx *cli.Context) error {
	t := time.Now()
	genesisBlock := &message.Block{}
	genesisBlock = &message.Block{Timestamp: t.String(), Hash: message.CalculateBlockHash(genesisBlock)}
	p.blockChain = append(p.blockChain, genesisBlock)

	// 命令行传参
	port := ctx.Int(portFlagForNew)
	target := ctx.String(targetFlagForNew)
	secio := ctx.Bool(secioFlagForNew)
	seed := ctx.Int64(seedFlagForNew)

	if port == 0 {
		logger.Fatal("请提供一个端口号")
	}
	// 构造一个host 监听地址
	ha, err := MakeBasicHost(port, secio, seed)
	if err != nil {
		return err
	}

	if target == "" {
		logger.Info("等待节点连接...")
		ha.SetStreamHandler("/p2p/1.0.0", p.HandleStream)
		select {}
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", p.HandleStream)
		ipfsAddr, err := ma.NewMultiaddr(target)
		if err != nil {
			return err
		}
		pid, err := ipfsAddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			return err
		}

		peerid, err := peer.Decode(pid)
		if err != nil {
			return err
		}

		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
		targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

		// 现在我们有一个peerID和一个targetaddr，所以我们添加它到peerstore中。 让libP2P知道如何连接到它。
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
		logger.Info("打开Stream")

		// 构建一个新的stream从hostB到hostA
		// 使用了相同的/p2p/1.0.0 协议
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			return err
		}

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go p.writeData(rw)
		go p.readData(rw)
		select {}
	}
}

// SavePeer 将加入到网络中的节点信息保存到配置文件中，方便后续投票与选择
func SavePeer(name string) {
	vote := DefaultVote // 默认的投票数目
	f, err := os.OpenFile(storage.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Errorf(err.Error())
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {

		}
	}(f)

	_, err = f.WriteString(name + ":" + strconv.Itoa(vote) + "\n")
	if err != nil {
		return
	}

}
