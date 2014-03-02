package main

import (
	"flag"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var sep = flag.Bool("s", false, "separate arguments with spaces")
var key = flag.String("key", "", "etcd key")
var host = flag.String("host", "", "etcd host")
var hosts = []string{}

func main() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-s ") {
		args := []string{}
		for _, arg := range os.Args {
			args = append(args, strings.Split(arg, " ")...)
		}
		os.Args = args
	}
	flag.Parse()

	if *key == "" {
		*key = os.Getenv("ETCDENV_KEY")
	}
	if *key == "" {
		fmt.Fprintln(os.Stderr, "etcdenv [-key=key] [...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *host == "" {
		*host = os.Getenv("ETCDENV_HOST")
	}
	if *host != "" {
		hosts = []string{*host}
	}
	client := etcd.NewClient(hosts)
	res, err := client.Get(*key, true, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "etcdenv: %s\n", err)
		os.Exit(1)
	}

	envs := os.Environ()
	for _, n := range res.Node.Nodes {
		key := strings.Split(n.Key, "/")
		k, v := strings.ToUpper(key[len(key)-1]), n.Value
		envs = append(envs, k+"="+v)
	}

	if flag.NArg() == 0 {
		for _, env := range envs {
			line := fmt.Sprintf("%q", env)
			fmt.Println(line[1 : len(line)-1])
		}
		os.Exit(0)
	}

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Env = envs
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if cmd.Process == nil {
		fmt.Fprintf(os.Stderr, "etcdenv: %s\n", err)
		os.Exit(1)
	}
	os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
}
