// Package vote /**
package vote

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/wsw365904/dpos/storage"
	"github.com/wsw365904/wswlog/wlogging"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var logger = wlogging.MustGetLoggerWithoutName()

// NodeVote 节点投票命令
var NodeVote = cli.Command{
	Name:  "vote",
	Usage: "vote for node",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "node name",
		},
		cli.IntFlag{
			Name:  "v",
			Value: 0,
			Usage: "vote number",
		},
	},
	Action: func(context *cli.Context) error {
		if err := Vote(context); err != nil {
			return err
		}
		return nil
	},
}

// Vote for node. The votes of node is origin vote plus new vote.
// votes = originVote + vote
func Vote(context *cli.Context) error {
	name := context.String("name")
	vote := context.Int("v")

	if name == "" {
		logger.Errorf("节点名称不能为空")
		return errors.Errorf("节点名称不能为空")
	}

	if vote < 1 {
		logger.Errorf("最小投票数目为1")
		return errors.Errorf("最小投票数目为1")
	}

	f, err := ioutil.ReadFile(storage.FileName)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	res := strings.Split(string(f), "\n")

	voteMap := make(map[string]string)
	for _, node := range res {
		nodeSplit := strings.Split(node, ":")
		if len(nodeSplit) > 1 {
			voteMap[nodeSplit[0]] = fmt.Sprintf("%s", nodeSplit[1])
		}
	}

	originVote, err := strconv.Atoi(voteMap[name])
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	votes := originVote + vote
	voteMap[name] = fmt.Sprintf("%d", votes)

	logger.Infof("节点%s新增票数%d", name, vote)
	str := ""
	for k, v := range voteMap {
		str += k + ":" + v + "\n"
	}

	file, err := os.OpenFile(storage.FileName, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	_, err = file.WriteString(str)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Error(err)
		}
	}(file)

	return nil
}
