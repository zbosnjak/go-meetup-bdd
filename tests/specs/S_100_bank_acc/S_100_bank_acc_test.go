package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/kelseyhightower/envconfig"
	"github.com/zbosnjak/go-meetup-bdd/internal/rest"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type testState struct {
	appAddress          string
	httpClient          http.Client
	savingAcc           int32
	currentAcc          int32
	responseMessage     string
	balanceWithInterest float32
}

var config struct {
	AppAddress string `envconfig:"app_addr" default:"127.0.0.1:9099"`
}

func (h *testState) myCurrentAccountHasABalanceOf(currentAcc int32) error {
	h.currentAcc = currentAcc
	return nil
}

func (h *testState) mySavingsAccountHasABalanceOf(savingAcc int32) error {
	h.savingAcc = savingAcc
	bankAcc := rest.BankAcc{
		SavingAcc:  savingAcc,
		CurrentAcc: h.currentAcc,
	}

	url := fmt.Sprintf("http://%v/account", h.appAddress)
	marshal, err := json.Marshal(bankAcc)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var acc rest.BankAcc
	err = json.Unmarshal(body, &acc)
	if err != nil {
		return err
	}
	if bankAcc != acc {
		return fmt.Errorf("current and saving account not set!")
	}
	return nil
}

func (h *testState) iTransferFromMyCurrentAccountToMySavingsAccount(amountToTransfer int32) error {
	bankAcc := rest.BankAcc{
		SavingAcc:  amountToTransfer,
		CurrentAcc: amountToTransfer * -1,
	}

	url := fmt.Sprintf("http://%v/account", h.appAddress)
	marshal, err := json.Marshal(bankAcc)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var acc rest.BankAcc
		err = json.Unmarshal(body, &acc)
		if err != nil {
			return err
		}
		if bankAcc != acc {
			return fmt.Errorf("current and saving account not set!")
		}
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		h.responseMessage = string(body)
	}

	return nil
}

func (h *testState) iShouldHaveInMyCurrentAccount(currentAcc int32) error {

	url := fmt.Sprintf("http://%v/account", h.appAddress)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var acc rest.BankAcc
	err = json.Unmarshal(body, &acc)
	if err != nil {
		return err
	}
	if acc.CurrentAcc != currentAcc {
		return fmt.Errorf("current acc: expected %v - actual %v", currentAcc, acc.CurrentAcc)
	}
	return nil
}

func (h *testState) iShouldHaveInMySavingsAccount(savingAcc int32) error {
	url := fmt.Sprintf("http://%v/account", h.appAddress)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var acc rest.BankAcc
	err = json.Unmarshal(body, &acc)
	if err != nil {
		return err
	}
	if acc.SavingAcc != savingAcc {
		return fmt.Errorf("saving acc: expected %v - actual %v", savingAcc, acc.SavingAcc)
	}
	return nil
}

func (h *testState) iShouldReceiveAnError(error string) error {
	if !strings.Contains(h.responseMessage, error) {
		return fmt.Errorf("Expected response message is %v - actual %v", error, h.responseMessage)
	}
	return nil
}

func (h *testState) itheMonthlyInterestIsCalculated(rate float32) error {

	calWithInterestedRate := rest.CalculateBalWithInterestRate{
		NumberOfMonths: 12,
		Rate:           rate,
	}

	url := fmt.Sprintf("http://%v/interestrate", h.appAddress)
	marshal, err := json.Marshal(calWithInterestedRate)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	response, err := h.httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var bal rest.BalanceWithRate
		err = json.Unmarshal(body, &bal)
		if err != nil {
			return err
		}
		h.balanceWithInterest = bal.Result
	} else {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		h.responseMessage = string(body)
	}

	return nil
}

func (h *testState) iShouldHaveAtEndOfYear(expectedResultInt, decimal int32) error {

	decimalValue := float32(decimal)
	for decimalValue > 1 {
		decimalValue = decimalValue / 10
	}
	expectedResult := float32(expectedResultInt) + decimalValue
	if h.balanceWithInterest != expectedResult {
		return fmt.Errorf("Expected balance %v - actual %v", expectedResult, h.balanceWithInterest)
	}
	return nil
}

func iWantToManageMyAccounts() error {
	fmt.Println("First setup method!")
	return nil
}

func transferMoneyFromOneToAnother() error {
	fmt.Println("Second setup method!")
	return nil
}

func (h *testState) iPerformActionWithAmount(arg1 *gherkin.DataTable) error {
	// index starts from 1, row 0 are headers
	for index := 1; index < len(arg1.Rows); index++ {
		row := arg1.Rows[index]
		action := row.Cells[0].Value
		amountStr := row.Cells[1].Value
		fmt.Printf("performing action %v - amount %v\n", action, amountStr)
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return err
		}
		var bankAcc = rest.BankAcc{}
		if action == "transfer" {
			bankAcc = rest.BankAcc{
				SavingAcc:  int32(amount),
				CurrentAcc: int32(amount) * -1,
			}
		} else if action == "deposit" {
			bankAcc = rest.BankAcc{
				SavingAcc:  0,
				CurrentAcc: int32(amount),
			}
		}

		url := fmt.Sprintf("http://%v/account", h.appAddress)
		marshal, err := json.Marshal(bankAcc)
		if err != nil {
			return err
		}
		request, err := http.NewRequest("PUT", url, bytes.NewReader(marshal))
		if err != nil {
			return err
		}
		response, err := h.httpClient.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode == http.StatusOK {
			fmt.Println("successful action")
		} else {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
			fmt.Println(string(body))
		}
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	envconfig.Process("", &config)
	//
	holder := testState{
		appAddress: config.AppAddress,
		httpClient: http.Client{},
	}
	//TODO
	//s.AfterSuite(func() {
	//	fmt.Println("AfterSuite")
	//})
	//s.BeforeScenario(func(interface{}) {
	//	fmt.Println("BeforeScenario")
	//})
	//
	//s.AfterScenario(func(interface{}, error) {
	//	fmt.Println("AfterScenario")
	//})

	s.Step(`^my Current account has a balance of (\d+)\.(\d+)$`, holder.myCurrentAccountHasABalanceOf)
	s.Step(`^my Savings account has a balance of (\d+)\.(\d+)$`, holder.mySavingsAccountHasABalanceOf)
	s.Step(`^I transfer (\d+)\.(\d+) from my Current account to my Savings account$`, holder.iTransferFromMyCurrentAccountToMySavingsAccount)
	s.Step(`^I should have (\d+)\.(\d+) in my Current account$`, holder.iShouldHaveInMyCurrentAccount)
	s.Step(`^I should have (\d+)\.(\d+) in my Savings account$`, holder.iShouldHaveInMySavingsAccount)
	s.Step(`^I should receive an "([^"]*)" error$`, holder.iShouldReceiveAnError)

	s.Step(`^I have a saving account with a balance of (\d+)$`, holder.mySavingsAccountHasABalanceOf)
	s.Step(`^the monthly interest (\d+) is calculated$`, holder.itheMonthlyInterestIsCalculated)
	s.Step(`^I should have (\d+),(\d+) at end of year$`, holder.iShouldHaveAtEndOfYear)

	s.Step(`^I perform <action> with <amount>$`, holder.iPerformActionWithAmount)

	s.Step(`^I want to manage my accounts$`, iWantToManageMyAccounts)
	s.Step(`^transfer money from one to another$`, transferMoneyFromOneToAnother)

}
