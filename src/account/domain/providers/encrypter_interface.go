package providers

type Encrypter interface {
	Encrypt(password string) (string, error)
	Compare(password, hash string) bool
}
