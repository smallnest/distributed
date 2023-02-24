package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

var (
	addr        = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	barrierName = flag.String("name", "my-test-doublebarrier", "barrier name")
	count       = flag.Int("c", 2, "")
)

func main() {
	flag.Parse()

	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")

	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// 创建session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	// 创建/获取队列
	b := recipe.NewDoubleBarrier(s1, *barrierName, *count)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "enter": // 持有这个barrier
			b.Enter()
			fmt.Println("enter")
		case "leave": // 释放这个barrier
			b.Leave()
			fmt.Println("leave")
		case "quit", "exit": //退出
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
