package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Account struct {
	ID    string `json:"id"`
	Balance int    `json:"balance"`
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, fromID string, toID string, amount int) error {
	fromAccount, err := ctx.GetStub().GetState(fromID)
	if err != nil {
		return fmt.Errorf("Failed to read from world state: %v", err)
	}
	if fromAccount == nil {
		return fmt.Errorf("Account does not exist: %s", fromID)
	}

	toAccount, err := ctx.GetStub().GetState(toID)
	if err != nil {
		return fmt.Errorf("Failed to read from world state: %v", err)
	}
	if toAccount == nil {
		return fmt.Errorf("Account does not exist: %s", toID)
	}

	from := new(Account)
	to := new(Account)

	err = json.Unmarshal(fromAccount, &from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(toAccount, &to)
	if err != nil {
		return err
	}

	if from.Balance < amount {
		return fmt.Errorf("Insufficient balance in account: %s", fromID)
	}

	from.Balance -= amount
	to.Balance += amount

	fromBytes, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(fromID, fromBytes)
	if err != nil {
		return err
	}

	toBytes, err := json.Marshal(to)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(toID, toBytes)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v", err)
	}
}
