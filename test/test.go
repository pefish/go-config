package main

import (
	"fmt"
	"github.com/pefish/go-config"
)

func main() {
	go_config.Config.MustMergeFlag()
	fmt.Println(go_config.Config.GetAll())
}
