package agent

import "encoding/base64"

// ObfuscationResult represents the result of an obfuscation request
type ObfuscationResult struct {
	ObfuscatedValue string `json:"obfuscatedValue"`
}

// ObfuscateStringWithKey obfuscates a string using a key
// It XORs each byte of the value using part of the key
// and converts it to a UTF8-string value.
// This is useful for obfuscating configuration values
func obfuscateStringWithKey(textToObfuscate string, encodingKey string) string {

	encodingKeyBytes := []byte(encodingKey)
	encodingKeyLen := len(encodingKeyBytes)

	textToObfuscateBytes := []byte(textToObfuscate)
	textToObfuscateLen := len(textToObfuscate)

	if encodingKeyLen == 0 || textToObfuscateLen == 0 {
		return ""
	}

	obfuscatedTextBytes := make([]byte, textToObfuscateLen)

	for i := 0; i < textToObfuscateLen; i++ {
		obfuscatedTextBytes[i] = textToObfuscateBytes[i] ^ encodingKeyBytes[i%encodingKeyLen]
	}

	obfuscatedText := base64.StdEncoding.EncodeToString(obfuscatedTextBytes)

	return obfuscatedText
}
