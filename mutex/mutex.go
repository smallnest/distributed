package main

import (
	"context"
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

	useMutex(cli)
}

func useMutex(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, *lockName)

	log.Printf("before acquiring. key: %s", m1.Key())
	// 请求锁
	log.Println("acquiring lock")
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Printf("acquired lock. key: %s", m1.Key())

	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)

	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("released lock")
}
