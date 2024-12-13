package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dvcrn/pocketsmith-anapay/anapay"

	"github.com/dvcrn/pocketsmith-go"
)

const INSTITUION_NAME = "ANA Pay"
const ACCOUNT_NAME = "ANA Pay"

type Config struct {
	AnapayUsername   string
	AnapayPassword   string
	PocketsmithToken string

	NumTransactions int
}

func getConfig() *Config {
	config := &Config{}

	// Define command-line flags
	flag.StringVar(&config.AnapayUsername, "username", os.Getenv("ANAPAY_USERNAME"), "Anapay username")
	flag.StringVar(&config.AnapayPassword, "password", os.Getenv("ANAPAY_PASSWORD"), "Anapay password")
	flag.StringVar(&config.PocketsmithToken, "token", os.Getenv("POCKETSMITH_TOKEN"), "Pocketsmith API token")
	flag.IntVar(&config.NumTransactions, "num-transactions", 100, "Number of transactions to parse")
	flag.Parse()

	// Validate required fields
	if config.AnapayUsername == "" {
		fmt.Println("Error: Anapay username is required. Set via -username flag or ANAPAY_USERNAME environment variable")
		os.Exit(1)
	}
	if config.AnapayPassword == "" {
		fmt.Println("Error: Anapay password is required. Set via -password flag or ANAPAY_PASSWORD environment variable")
		os.Exit(1)
	}
	if config.PocketsmithToken == "" {
		fmt.Println("Error: Pocketsmith token is required. Set via -token flag or POCKETSMITH_TOKEN environment variable")
		os.Exit(1)
	}

	return config
}

func findOrCreateAccount(ps *pocketsmith.Client, userID int, accountName string) (*pocketsmith.Account, error) {
	account, err := ps.FindAccountByName(userID, accountName)
	if err != nil {
		if err != pocketsmith.ErrNotFound {
			return nil, err
		}

		institution, err := ps.FindInstitutionByName(userID, INSTITUION_NAME)
		if err != nil {
			if err != pocketsmith.ErrNotFound {
				return nil, err
			}

			institution, err = ps.CreateInstitution(userID, INSTITUION_NAME, "jpy")
			if err != nil {
				return nil, err
			}
		}

		account, err := ps.CreateAccount(userID, institution.ID, accountName, "jpy", pocketsmith.AccountTypeCredits)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	return account, nil
}

func main() {
	config := getConfig()

	ps := pocketsmith.NewClient(config.PocketsmithToken)
	res, err := ps.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	account, err := findOrCreateAccount(ps, res.ID, ACCOUNT_NAME)
	if err != nil {
		fmt.Println("could not find or create account")
		panic(err)
	}

	loginResponse, err := anapay.Login(config.AnapayUsername, config.AnapayPassword)
	if err != nil {
		panic(err)
	}

	anaPayAccounts, err := anapay.GetAccounts(loginResponse.AccessToken)
	if err != nil {
		panic(err)
	}

	balanceFloat, err := strconv.ParseFloat(anaPayAccounts.Balance, 64)
	if err != nil {
		panic(err)
	}
	updateAccountRes, err := ps.UpdateTransactionAccount(account.PrimaryTransactionAccount.ID, account.PrimaryTransactionAccount.Institution.ID, balanceFloat, time.Now().Format("2006-01-02"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Updated account balance: ", updateAccountRes.StartingBalance)

	var allTxs []anapay.Transaction
	pageNum := 1
	for {
		txs, err := anapay.GetTransactions(loginResponse.AccessToken, pageNum, 999)
		if err != nil {
			panic(err)
		}

		if len(txs) == 0 {
			break
		}

		allTxs = append(allTxs, txs...)
		if len(allTxs) > config.NumTransactions {
			break
		}

		pageNum++
	}

	fmt.Println("Found ", len(allTxs), " transactions")

	repeatedExistingTransactions := 0
	for _, tx := range allTxs {
		if repeatedExistingTransactions > 10 {
			fmt.Println("Too many repeated existing transactions, exiting")
			break
		}

		amount := 0.0
		if tx.Amount != "" {
			parsedAmount, _ := strconv.ParseFloat(tx.Amount, 64)
			amount = -parsedAmount

			if tx.DealType == "05" || tx.DealType == "06" || tx.DelKbn == "02" || tx.DelKbn == "07" || tx.DelKbn == "08" {
				amount = parsedAmount
			}
		}

		bookingText := tx.DescriptionType
		isTransfer := false
		switch {
		case tx.DealType == "05":
			bookingText = "チャージ"
			isTransfer = true
		case tx.DealType == "06":
			bookingText = "キャッシュバック"
		case tx.DescriptionType == "3001":
			bookingText = "クレジットカード"
		case tx.DescriptionType == "3006":
			bookingText = "Apple Pay"
		case tx.DescriptionType == "3007":
			bookingText = "キャッシュバック"
		case tx.DescriptionType == "3009":
			bookingText = "オートチャージ"
			isTransfer = true
		case tx.DescriptionType == "1017":
			bookingText = "バーチャルプリペイドカード"
		case tx.DescriptionType == "1018":
			bookingText = "VISAタッチ払い"
		case tx.DescriptionType == "1019":
			bookingText = "iDタッチ払い"
		default:
			bookingText = "Unknown Type"
		}

		name := tx.ShopName
		if name == "" {
			name = bookingText
		}

		// parse date from 20241211070348 into golang time
		date, err := time.Parse("20060102150405", tx.SaleDatetime)
		if err != nil {
			fmt.Println("Error parsing date: ", err)
			continue
		}
		createTx := &pocketsmith.CreateTransaction{
			Payee:        strings.TrimSpace(name),
			Amount:       amount,
			Date:         tx.SaleDatetime[0:8], // Convert YYYYMMDDHHMMSS to YYYYMMDD
			IsTransfer:   isTransfer,
			Note:         "",
			Memo:         strings.TrimSpace(bookingText),
			ChequeNumber: tx.WalletSettlementNo,
		}

		searchByChequeResponse, err := ps.SearchTransactionsByChequeNumber(account.PrimaryTransactionAccount.ID, date, tx.WalletSettlementNo)
		if err != nil {
			fmt.Println("Error searching for transaction: ", err)
			continue
		}

		if len(searchByChequeResponse) > 0 {
			fmt.Println("Found transaction already, won't add it again: ", name)
			repeatedExistingTransactions++
			continue
		}

		fmt.Println("Creating transaction with createTx: ", createTx.Payee, createTx.Amount, createTx.Date, createTx.IsTransfer, createTx.Note)
		_, err = ps.AddTransaction(account.TransactionAccounts[0].ID, createTx)
		if err != nil {
			fmt.Printf("Error creating transaction: %v\n", err)
			continue
		}
	}
}
