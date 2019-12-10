package go_config

import (
	"errors"
)

func Parse(arguments []string) (map[string]string, error) {

	result := map[string]string{}
	for {
		seen, err := parseOne(result, &arguments)
		if seen {
			continue
		}
		if err == nil {
			break
		}
		panic(err)
	}
	return result, nil
}


func parseOne(map_ map[string]string, arguments *[]string) (bool, error) {
	args := *arguments
	if len(args) < 2 {
		return false, nil
	}
	target := args[0:2]
	s := target[0]
	if len(s) < 2 || s[0] != '-' {
		return false, errors.New(`syntax error`)
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
	}
	name := s[numMinuses:]

	map_[name] = target[1]

	*arguments = args[2:]

	return true, nil
}

