package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coreos/etcd/clientv3"
	recipe "github.com/coreos/etcd/contrib/recipes"
)

var (
	addr        = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	barrierName = flag.String("name", "my-test-queue", "barrier name")
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

	// 创建/获取栅栏 
	b := recipe.NewBarrier(cli, *barrierName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "hold": // 持有这个barrier
			b.Hold()
			fmt.Println("hold")
		case "release": // 释放这个barrier
			b.Release()
			fmt.Println("released")
		case "wait": // 等待barrier被释放
			b.Wait()
			fmt.Println("after wait")
		case "quit", "exit": //退出
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
