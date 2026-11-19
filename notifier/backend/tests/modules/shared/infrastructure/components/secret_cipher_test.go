package components

import (
    "testing"

    "notifier-backend/tests"

    "notifier-backend/src/modules/shared/infrastructure/components/secretcipher"
)

const testSecretsKey = "0123456789abcdef0123456789abcdef"

func TestSecretRoundTripWithRandomNonce(t *testing.T) {
    secretCipher := secretcipher.SecretCipherFromKey(testSecretsKey)
    encrypted1, err := secretCipher.GetEncrypted("abcd efgh ijkl mnop")
    tests.AssertNil(t, err)
    encrypted2, _ := secretCipher.GetEncrypted("abcd efgh ijkl mnop")
    tests.AssertNotEquals(t, encrypted1, encrypted2)

    decrypted, err := secretCipher.GetDecrypted(encrypted1)
    tests.AssertNil(t, err)
    tests.AssertEquals(t, "abcd efgh ijkl mnop", decrypted)
}

func TestWrongKeyOrShortKeyFails(t *testing.T) {
    encrypted, _ := secretcipher.SecretCipherFromKey(testSecretsKey).GetEncrypted("secreto")
    _, err := secretcipher.SecretCipherFromKey(testSecretsKey + "x").GetDecrypted(encrypted)
    tests.AssertNotEquals(t, nil, err)
    _, err = secretcipher.SecretCipherFromKey("corta").GetEncrypted("secreto")
    tests.AssertNotEquals(t, nil, err)
}
