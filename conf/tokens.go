package conf

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
	return text
}
