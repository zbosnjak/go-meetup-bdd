package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BankAcc struct {
	SavingAcc  int32 `json:"saving_acc"`
	CurrentAcc int32 `json:"current_acc"`
}

type CalculateBalWithInterestRate struct {
	NumberOfMonths int32   `json:"number_of_months"`
	Rate           float32 `json:"interest_rate"`
}

type BalanceWithRate struct {
	Result float32 `json:"result"`
}

type AccHandler struct {
	acc BankAcc
}

// GetBankAccHandler returns current state of bank acc
func GetBankAccHandler(accounts *AccHandler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("GetAccountsBalance request received")

		acc := accounts.acc
		bytes, err := json.Marshal(acc)
		if err != nil {
			msg := "failed to marshal bank accounts"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			msg := "failed to write response"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// SetBankAccHandler set and returns current state of bank acc
func SetBankAccHandler(accounts *AccHandler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("SetAccountsBalance request received")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		var acc BankAcc
		err = json.Unmarshal(body, &acc)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		accounts.acc = acc
		bytes, err := json.Marshal(acc)
		if err != nil {
			msg := "failed to marshal bank accounts"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			msg := "failed to write response"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateBankAccHandler set and returns current state of bank acc
func UpdateBankAccHandler(accounts *AccHandler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("UpdateAccountsBalance request received")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		var acc BankAcc
		err = json.Unmarshal(body, &acc)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if accounts.acc.CurrentAcc+acc.CurrentAcc < 0 {
			msg := "Insufficient funds on current account"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		if accounts.acc.SavingAcc+acc.SavingAcc < 0 {
			msg := "Insufficient funds on saving account"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}
		accounts.acc.CurrentAcc += acc.CurrentAcc
		accounts.acc.SavingAcc += acc.SavingAcc
		bytes, err := json.Marshal(acc)
		if err != nil {
			msg := "failed to marshal bank accounts"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			msg := "failed to write response"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// CalculateBalanceWithInterestRateHandler return calculate balance with interest rate
func CalculateBalanceWithInterestRateHandler(accounts *AccHandler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("request received")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		var rate CalculateBalWithInterestRate
		err = json.Unmarshal(body, &rate)
		if err != nil {
			msg := fmt.Sprintf("Error reading body: %v\n", err)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		var i int32 = 0
		sum := float32(accounts.acc.SavingAcc)
		for ; i < rate.NumberOfMonths; i++ {
			sum = sum * (1 + rate.Rate/100)
		}
		acc := BalanceWithRate{
			Result: sum,
		}
		bytes, err := json.Marshal(acc)
		if err != nil {
			msg := "failed to marshal bank accounts"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			msg := "failed to write response"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
