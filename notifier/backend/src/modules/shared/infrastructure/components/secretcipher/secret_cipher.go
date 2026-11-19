package secretcipher

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"

    "notifier-backend/src/modules/shared/domain/enums"
    "notifier-backend/src/modules/shared/infrastructure/repositories/configuration"
)

// minimo de caracteres de APP_SECRETS_KEY: una clave corta se adivina aunque luego se derive con sha256
const minSecretsKeyLength = 32

// SecretCipher cifra los secretos de los perfiles de envio (app password, api key, token de bot) con AES-256-GCM.
// La clave AES es el sha256 de APP_SECRETS_KEY. El texto cifrado es base64(nonce + cifrado + tag): cada cifrado lleva
// su nonce aleatorio, asi que el mismo secreto cifrado dos veces no se parece.
type SecretCipher struct {
    secretsKey string
}

func GetInstance() *SecretCipher {
    return &SecretCipher{secretsKey: configuration.GetInstance().GetSecretsKey()}
}

// SecretCipherFromKey para los tests: sin pasar por el .env.
func SecretCipherFromKey(secretsKey string) *SecretCipher {
    return &SecretCipher{secretsKey: secretsKey}
}

func (c *SecretCipher) GetEncrypted(plainText string) (string, error) {
    aead, err := c.getAead()
    if err != nil {
        return "", err
    }
    nonce := make([]byte, aead.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", err
    }
    sealed := aead.Seal(nonce, nonce, []byte(plainText), nil)
    return base64.StdEncoding.EncodeToString(sealed), nil
}

func (c *SecretCipher) GetDecrypted(encryptedText string) (string, error) {
    aead, err := c.getAead()
    if err != nil {
        return "", err
    }
    sealed, err := base64.StdEncoding.DecodeString(encryptedText)
    if err != nil || len(sealed) < aead.NonceSize() {
        return "", errors.New(enums.SharedMessageEnum.SECRET_CANNOT_BE_DECRYPTED)
    }
    nonce, cipherText := sealed[:aead.NonceSize()], sealed[aead.NonceSize():]
    plainText, err := aead.Open(nil, nonce, cipherText, nil)
    if err != nil {
        return "", errors.New(enums.SharedMessageEnum.SECRET_CANNOT_BE_DECRYPTED)
    }
    return string(plainText), nil
}

func (c *SecretCipher) getAead() (cipher.AEAD, error) {
    if len(c.secretsKey) < minSecretsKeyLength {
        return nil, errors.New(enums.SharedMessageEnum.SECRETS_KEY_TOO_SHORT)
    }
    aesKey := sha256.Sum256([]byte(c.secretsKey))
    block, err := aes.NewCipher(aesKey[:])
    if err != nil {
        return nil, err
    }
    return cipher.NewGCM(block)
}
