package Authenticator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func loadTokensFromYaml(tokenYamlPath string) (tokenYaml, error) {
	var tokens tokenYaml

	data, err := os.ReadFile(tokenYamlPath)
	if err != nil {
		return tokenYaml{tokens: make(map[int]string)},
			fmt.Errorf("Error while rading yaml: %w", err)
	}

	err = yaml.Unmarshal(data, &tokens)
	if err != nil {
		return tokenYaml{tokens: make(map[int]string)},
			fmt.Errorf("Error while loading tokens from yaml: %w", err)
	}

	return tokens, nil
}
