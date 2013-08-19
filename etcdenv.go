package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"os/exec"
	"strings"
)

var url = flag.String("url", "", "etcd URL")

func main() {
	flag.Parse()
	if flag.NArg() == 0 || *url == "" {
		fmt.Fprintln(os.Stderr, "etcdenv [-url=url] [...]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	client := etcd.NewClient()
	res, err := client.Get(*url)
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
	os.Exit(0)
}
