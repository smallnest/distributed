package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	recipe "github.com/coreos/etcd/contrib/recipes"
)

var (
	addr     = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	lockName = flag.String("name", "my-test-lock", "lock name")
	action   = flag.String("rw", "w", "r means acquiring read lock, w means acquiring write lock")
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
	m1 := recipe.NewRWMutex(s1, *lockName)

	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "w": // 请求写锁
			testWriteLocker(m1)
		case "r": // 请求读锁
			testReadLocker(m1)
		default:
			fmt.Println("unknown action")
		}
	}
}

func testWriteLocker(m1 *recipe.RWMutex) {
	// 请求写锁
	log.Println("acquiring write lock")
	if err := m1.Lock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired write lock")

	// 等待一段时间
	time.Sleep(30 * time.Second)

	// 释放写锁
	if err := m1.Unlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released write lock")
}

func testReadLocker(m1 *recipe.RWMutex) {
	// 请求读锁
	log.Println("acquiring read lock")
	if err := m1.RLock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired read lock")

	// 等待一段时间
	time.Sleep(30 * time.Second)

	// 释放写锁
	if err := m1.RUnlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released read lock")
}
