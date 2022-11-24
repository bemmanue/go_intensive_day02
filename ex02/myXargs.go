package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func execCommand(arg string, group *sync.WaitGroup) {
	cmd := exec.Command(os.Args[1], os.Args[2], arg)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(stdout))
	group.Done()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Pass command like 'wc -l' or 'ls -la'\n")
		os.Exit(1)
	}

	argv, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	args := strings.Fields(string(argv))

	var group sync.WaitGroup
	for i := range args {
		group.Add(1)
		go execCommand(args[i], &group)
	}
	group.Wait()
}
