package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v\n", err)
	}
	defer cli.Close()
	fmt.Println("connect to etcd success")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := cli.Put(ctx, "set_kv", "gogo bebe")
	defer cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
	fmt.Printf("put resp: %v", resp)
	res, err := cli.Get(context.Background(), "set_kv")
	if err != nil {
		fmt.Printf("get from etcd failed, err: %v\n", err)
	}
	fmt.Printf("get resp: %v\n", res)
}
