package services

import (
	"pattern/internal/repositories"
	"pattern/internal/repositories/interfaces"
)

type TransactionFactory interface {
	CreateTransaction() interfaces.Transaction
}

type TransferTransactionFactory struct{}

func (f *TransferTransactionFactory) CreateTransaction() interfaces.Transaction {
	return &repositories.TransferTransaction{}
}

type WithdrawalsTransactionFactory struct{}

func (f *WithdrawalsTransactionFactory) CreateTransaction() interfaces.Transaction {
	return &repositories.WithdrawalsTransaction{}
}

type DepositTransactionFactory struct{}

func (f *DepositTransactionFactory) CreateTransaction() interfaces.Transaction {
	return &repositories.DepositTransaction{}
}
