package repositories

import (
	"fmt"
	"pattern/internal/repositories/interfaces"
)

type TransferTransaction struct {
}

func (t *TransferTransaction) CreateTransaction() interfaces.Transaction {
	return &TransferTransaction{}
}

func (t *TransferTransaction) Execute() {
	fmt.Println("transferred money")
}

type WithdrawalsTransaction struct {
}

func (t *WithdrawalsTransaction) CreateTransaction() interfaces.Transaction {
	return &WithdrawalsTransaction{}
}

func (t *WithdrawalsTransaction) Execute() {
	fmt.Println("money was withdrawn")
}

type DepositTransaction struct {
}

func (t *DepositTransaction) CreateTransaction() interfaces.Transaction {
	return &DepositTransaction{}
}

func (t *DepositTransaction) Execute() {
	fmt.Println("the transaction was carried out in a deposit")
}
