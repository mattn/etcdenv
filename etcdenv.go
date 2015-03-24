package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/coreos/go-etcd/etcd"
)

var sep = flag.Bool("s", false, "separate arguments with spaces")
var rec = flag.Bool("r", false, "recursively fetch child values")
var key = flag.String("key", "", "etcd key")
var host = flag.String("host", "", "etcd host")
var hosts []string
var envs []string

func prefixInSlice(list []string, s string) int {
	for i, entry := range list {
		if strings.HasPrefix(entry, s) {
			return i
		}
	}
	return -1
}

func handleNode(n *etcd.Node) {
	if !n.Dir {
		key := strings.Split(n.Key, "/")
		k, v := strings.ToUpper(key[len(key)-1]), n.Value
		// if k already exists and in recursive mode, append v with comma
		if i := prefixInSlice(envs, k+"="); i != -1 && *rec {
			envs[i] = fmt.Sprint(envs[i], ",", v)
		} else {
			envs = append(envs, k+"="+v)
		}
	} else if *rec {
		for _, n2 := range n.Nodes {
			handleNode(n2)
		}
	}
}

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
	if *host == "" {
		*host = os.Getenv("ETCDENV_HOST")
	}
	if *host != "" {
		hosts = strings.Split(*host, ",")
	}
	client := etcd.NewClient(hosts)
	res, err := client.Get(*key, true, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "etcdenv: %s\n", err)
		os.Exit(1)
	}

	envs = os.Environ()
	for _, n := range res.Node.Nodes {
		handleNode(n)
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
