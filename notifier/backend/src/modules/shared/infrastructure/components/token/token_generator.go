package token

import (
    "crypto/rand"
    "crypto/sha256"
    "crypto/subtle"
    "encoding/hex"
)

const tokenBytes = 32

// TokenGenerator tokens de API opacos: 32 bytes del CSPRNG en hex (256 bits). Con esa entropia basta un SHA-256 para
// guardarlos (no hace falta un hash lento tipo bcrypt, que es para contraseñas elegidas por humanos).
type TokenGenerator struct{}

func GetInstance() *TokenGenerator {
    return &TokenGenerator{}
}

func (t *TokenGenerator) GetNewToken() string {
    randomBytes := make([]byte, tokenBytes)
    _, _ = rand.Read(randomBytes)
    return hex.EncodeToString(randomBytes)
}

func (t *TokenGenerator) GetTokenHash(token string) string {
    tokenHash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(tokenHash[:])
}

// AreTokensEqual comparacion en tiempo constante (un token fijo de configuracion contra el que llega).
func (t *TokenGenerator) AreTokensEqual(expectedToken string, actualToken string) bool {
    return subtle.ConstantTimeCompare([]byte(expectedToken), []byte(actualToken)) == 1
}
