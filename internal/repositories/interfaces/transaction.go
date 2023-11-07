package interfaces

type Transaction interface {
	CreateTransaction() Transaction
	Execute()
}
