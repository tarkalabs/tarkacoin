package transaction

import (
	"errors"
)

type Transaction struct {
	Sender      string
	Inputs      []Transaction
	inputSum    int64
	Outputs     []TransactionOutput
	blockNumber int64
}

type TransactionOutput struct {
	Receiver string
	Amount   int64 // amount in satoshis
	Spent    bool
}

func (t *Transaction) isConfirmed() bool {
	return t.blockNumber >= 0
}

func NewUnconfirmedTransaction() *Transaction {
	return &Transaction{blockNumber: -1}
}

func (t *Transaction) AddInput(ti Transaction) error {
	// Checkif ti.Outpus.filter(&:Receiver == t.Sender)
	if !ti.isConfirmed() {
		return errors.New("The input transaction is not confirmed")
	}
	for _, outputT := range ti.Outputs {
		if outputT.Receiver == t.Sender && !outputT.Spent {
			t.Inputs = append(t.Inputs, ti)
			t.inputSum += outputT.Amount
			return nil
		} else {
			continue
		}
	}
	return errors.New("Unable to find a appropriate transaction output in the transaction")
}
func (t *Transaction) GetTotalOutput() int64 {
	var totalOutput int64
	for _, to := range t.Outputs {
		totalOutput += to.Amount
	}
	return totalOutput
}

func (t *Transaction) AdjustFees(fees int64) error {
	totalOutput := t.GetTotalOutput()
	for _, to := range t.Outputs {
		totalOutput += to.Amount
	}
	remaining := t.inputSum - totalOutput
	if remaining < fees {
		return errors.New("Not enough balance to add fees")
	}
	if remaining-fees > 0 {
		selfOutput := TransactionOutput{Receiver: t.Sender, Amount: remaining - fees, Spent: false}
		t.Outputs = append(t.Outputs, selfOutput)
		return nil
	} else {
		return errors.New("Not enough balance to add fees")
	}
}

func (t *Transaction) GetFees() int64 {
	return t.inputSum - t.GetTotalOutput()

}
