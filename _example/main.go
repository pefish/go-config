package main

import (
	"flag"
	"fmt"
	go_config "github.com/pefish/go-config"
	"os"
)

func main() {
	go_config.ConfigManagerInstance.MustLoadConfig(go_config.Configuration{
		ConfigFilepath: "./test.yaml",
	})
	var flagSet flag.FlagSet
	flagSet.String("abcd", "haha", "")
	flagSet.String("abcde", "123456", "")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	go_config.ConfigManagerInstance.MergeFlagSet(&flagSet)
	go_config.ConfigManagerInstance.MergeEnvs(map[string]string{
		"ABCD":  "abcd",
		"ABCDE": "abcde",
	})
	fmt.Println(go_config.ConfigManagerInstance.MustGetString("abcd"))
	fmt.Println(go_config.ConfigManagerInstance.MustGetString("abcde"))

	go_config.ConfigManagerInstance.MustLoadConfig(go_config.Configuration{
		SecretFilepath: "./test1.yaml",
	})
	var config struct {
		Test struct {
			Haha string `json:"haha"`
		} `json:"test"`
		Test1 struct {
			Test2 []struct {
				Test3 uint64 `json:"test3"`
			} `json:"test2"`
		} `json:"test1"`
	}
	go_config.ConfigManagerInstance.Unmarshal(&config)
	fmt.Println(config)
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
