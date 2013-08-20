package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"os/exec"
	"strings"
)

var key = flag.String("key", "", "etcd key")

func main() {
	flag.Parse()
	if *key == "" {
		fmt.Fprintln(os.Stderr, "etcdenv [-key=key] [...]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	client := etcd.NewClient()
	res, err := client.Get(*key)

	if flag.NArg() > 0 {
		if err == nil {
			for _, n := range res {
				key := strings.Split(n.Key, "/")
				k, v := strings.ToUpper(key[len(key)-1]), n.Value
				fmt.Println("creating env var:", k)
				os.Setenv(k, v)
			}
		}
		cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			os.Exit(1)
		}
	} else {
		if err == nil {
			for _, n := range res {
				key := strings.Split(n.Key, "/")
				k, v := strings.ToUpper(key[len(key)-1]), n.Value
				fmt.Printf("%s=%s\n", k, v)
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	os.Exit(0)
}
