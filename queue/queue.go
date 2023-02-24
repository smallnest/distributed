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
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	queueName = flag.String("name", "my-test-queue", "queue name")
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

	// 创建/获取队列
	q := recipe.NewQueue(cli, *queueName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push": // 加入队列
			if len(items) != 2 {
				fmt.Println("must set value to push")
				continue
			}
			q.Enqueue(items[1]) // 入队
		case "pop": // 从队列弹出
			v, err := q.Dequeue() // 出队
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(v) // 输出出队的元素
		case "quit", "exit": //退出
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
