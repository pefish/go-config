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
	flagSet.String("abcde", "123456", "")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	go_config.Config.MergeFlagSet(&flagSet)
	go_config.Config.MergeEnvs(map[string]string{
		"ABCD": "abcd",
		"ABCDE": "abcde",
	})
	fmt.Println(go_config.Config.MustGetString("abcd"))
	fmt.Println(go_config.Config.MustGetString("abcde"))
}

// go run ./_example/
// Output:
// 123
// 123456

// go run ./_example/ --abcd=124
// Output:
// 124
// 123456

// ABCD=125 go run ./_example/
// Output:
// 125
// 123456

// ABCD=125 go run ./_example/ --abcd=126
// Output:
// 126
// 123456

// ABCDE=127 go run ./_example/
// Output:
// 123
// 127
