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
	if err != nil {
		fmt.Fprintf(os.Stderr, "etcdenv: %s\n", err)
		os.Exit(1)
	}

	envs := []string{}
	for _, n := range res {
		key := strings.Split(n.Key, "/")
		k, v := strings.ToUpper(key[len(key)-1]), n.Value
		envs = append(envs, k + "=" + v)
	}
	if flag.NArg() == 0 {
		for _, env := range envs {
			line := fmt.Sprintf("%q", env)
			fmt.Println(line[1:len(line)-1])
		}
		os.Exit(0)
	}
	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Env = envs
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
