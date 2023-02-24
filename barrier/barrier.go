package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

var (
	addr        = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	barrierName = flag.String("name", "my-test-barrier", "barrier name")
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
		case "hold":
			holdBarrier(b)
		case "release":
			releaseBarrier(b)
		case "wait":
			waitBarrier(b)
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}

func holdBarrier(b *recipe.Barrier) {
	b.Hold()
	fmt.Println("hold")
}

func releaseBarrier(b *recipe.Barrier) {
	b.Release()
	fmt.Println("released")
}

func waitBarrier(b *recipe.Barrier) {
	b.Wait()
	fmt.Println("after wait")
}
