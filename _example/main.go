package main

import (
	"flag"
	"fmt"
	go_config "github.com/pefish/go-config"
	"os"
)

func main() {
	go_config.Config.MustLoadConfig(go_config.Configuration{
		SecretFilepath: "./_example/test.yaml",
	})
	var flagSet flag.FlagSet
	flagSet.String("abcd", "haha", "")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	go_config.Config.MergeFlagSet(&flagSet)
	go_config.Config.MergeEnvs(map[string]string{
		"ABC": "abcd",
	})
	fmt.Println(go_config.Config.MustGetString("abcd"))
}

// go run ./_example/
// Output:
// 123

// go run ./_example/ --abc=124
// Output:
// 124

// ABC=125 go run ./_example/
// Output:
// 125

// ABC=125 go run ./_example/ --abc=126
// Output:
// 126
