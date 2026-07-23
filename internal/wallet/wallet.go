package wallet

type Wallet interface {
	Address(index uint32) (string, error)
	PrivateKey(index uint32) (string, error)
}
