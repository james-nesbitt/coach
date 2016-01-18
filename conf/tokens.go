package conf

import (
	"strings"
)

const (
	TOKEN_KEY_PREFIX = "%"
	TOKEN_KEY_SUFFIX = ""
)

func MakeTokens() Tokens {
	return Tokens{}
}

type Tokens map[string]string

// Set a Token
func (tokens Tokens) SetToken(key string, value string) {
	tokens[key] = value
}

// Replce any tokens in the string with tokens from the token map
func (tokens *Tokens) TokenReplace(text string) string {
	for key, value := range *tokens {
		key = TOKEN_KEY_PREFIX+key+TOKEN_KEY_SUFFIX
		text = strings.Replace(text, key, value, -1)
	}
	return text
}
