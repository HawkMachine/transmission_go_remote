package main

import (
	"flag"
	"fmt"

	"github.com/HawkMachine/transmission_go_remote"
)

var (
	a = flag.String("a", "", "Address")
	u = flag.String("u", "", "Username")
	p = flag.String("p", "", "Password")
)

func main() {
	flag.Parse()

	r, _ := remote.New(*a, *u, *p, "")
	ts, err := r.ListAll()
	if err != nil {
		fmt.Printf("Failed to list torrents: %v", err)
		return
	}
	for _, t := range ts {
		fmt.Printf("%d: %s", t.Id, t.Name)
	}
}
