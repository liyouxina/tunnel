package main

import (
	"flag"
	"github.com/liyouxina/tunnel/pkg/agent"
	"github.com/liyouxina/tunnel/pkg/logger"
)

var tunnelServer = flag.String("tunnelServer", "localhost:8091", "tunnelServer")
var targetServer = flag.String("targetServer", "localhost:8080", "targetServer")
var maxAgentCnt = flag.Int("maxAgentCnt", 100, "maxAgentCnt")

var log = logger.Logger

func init() {
	flag.Parse()
}

func main() {
	agent.StartAgents(*tunnelServer, *targetServer, *maxAgentCnt)
	log.Infof("启动成功")
	select {}
}
