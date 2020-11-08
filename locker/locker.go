package main

import (
	"flag"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

var (
	addr     = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	lockName = flag.String("name", "my-test-lock", "lock name")
)

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	// etcd地址
	endpoints := strings.Split(*addr, ",")
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	useLock(cli)
}

func useLock(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	locker := concurrency.NewLocker(s1, *lockName)

	// 请求锁
	log.Println("acquiring lock")
	locker.Lock()
	log.Println("acquired lock")

	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	locker.Unlock()

	log.Println("released lock")
}
