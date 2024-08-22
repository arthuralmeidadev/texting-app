package utils

type cryptoMng struct {
    publicKeyPath string
    privateKeyPath string
}

func (c *cryptoMng) ReadPublicKey(path string) error {
    return nil
}

func (c *cryptoMng) ReadPrivateKey(path string) error {
    return nil
}

func (c *cryptoMng) Encrypt(value string, factor string) (*string, error) {
    return nil, nil
}

func (c *cryptoMng) Decrypt(ciphertext string, factor string) (*string, error) {
    return nil, nil
}

func NewCryptoManager() *cryptoMng {
    return &cryptoMng{}
}
