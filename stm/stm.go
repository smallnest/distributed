package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

var (
	addr = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
)

func main() {
	flag.Parse()

	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 设置5个账户，每个账号都有100元，总共500元
	totalAccounts := 5
	for i := 0; i < totalAccounts; i++ {
		k := fmt.Sprintf("accts/%d", i)
		if _, err = cli.Put(context.TODO(), k, "100"); err != nil {
			log.Fatal(err)
		}
	}

	// STM的应用函数，主要的事务逻辑
	exchange := func(stm concurrency.STM) error {
		// 随机得到两个转账账号
		from, to := rand.Intn(totalAccounts), rand.Intn(totalAccounts)
		if from == to {
			// 自己不和自己转账
			return nil
		}
		// 读取账号的值
		fromK, toK := fmt.Sprintf("accts/%d", from), fmt.Sprintf("accts/%d", to)
		fromV, toV := stm.Get(fromK), stm.Get(toK)
		fromInt, toInt := 0, 0
		fmt.Sscanf(fromV, "%d", &fromInt)
		fmt.Sscanf(toV, "%d", &toInt)

		// 把源账号的一半转账给目标账号
		xfer := fromInt / 2
		fromInt, toInt = fromInt-xfer, toInt+xfer

		// 把转账后的值写回
		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
		stm.Put(toK, fmt.Sprintf("%d", toInt))
		return nil
	}

	// 启动10个goroutine进行转账操作
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if _, serr := concurrency.NewSTM(cli, exchange); serr != nil {
					log.Fatal(serr)
				}
			}
		}()
	}
	wg.Wait()

	// 检查账号最后的数目
	sum := 0
	accts, err := cli.Get(context.TODO(), "accts/", clientv3.WithPrefix()) // 得到所有账号
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range accts.Kvs { // 遍历账号的值
		v := 0
		fmt.Sscanf(string(kv.Value), "%d", &v)
		sum += v
		log.Printf("account %s: %d", kv.Key, v)
	}

	log.Println("account sum is", sum) // 总数
}
