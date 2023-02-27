package db

import (
	"context"
	"fmt"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/pkg/apiBody"
	"github.com/Template7/common/structs"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var (
	testUser = structs.User{
		UserId: "testUserId",
		Mobile: "testMobile",
		Email:  "testEmail@test.com",
		BasicInfo: structs.UserInfo{
			NickName: "testUserNickName",
			Bio:      "testBio",
		},
	}

	testSender = structs.User{
		UserId: "testSender",
		Mobile: "senderMobile",
		Email:  "senderEmail@test.com",
		BasicInfo: structs.UserInfo{
			NickName: "testSenderNickName",
			Bio:      "senderBio",
		},
	}
	testReceiver = structs.User{
		UserId: "testReceiver",
		Mobile: "receiverMobile",
		Email:  "receiverEmail@test.com",
		BasicInfo: structs.UserInfo{
			NickName: "testReceiverNickName",
			Bio:      "receiverBio",
		},
	}
	testMoneyNtd = structs.Money{
		Currency: structs.CurrencyNTD,
		Amount:   100,
		Unit:     structs.UnitPico,
	}
	testMoneyUsd = structs.Money{
		Currency: structs.CurrencyUSD,
		Amount:   100,
		Unit:     structs.UnitPico,
	}
)

func TestMain(m *testing.M) {
	viper.AddConfigPath("../../../configs")
	c := config.New()
	db := fmt.Sprintf("template7")
	c.Mongo.Db = db
	c.MySql.Db = db
	c.MySql.ConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.MySql.Username, c.MySql.Password, c.MySql.Host, c.MySql.Port, c.MySql.Db)
	code := m.Run()

	teardown(db)
	os.Exit(code)
}

func Test_dbClient(t *testing.T) {
	t.Run("createUser", func(t *testing.T) {
		if err := New().CreateUser(testUser); err != nil {
			t.Errorf("CreateUser() error = %v", err)
			return
		}
		if err := New().CreateUser(testSender); err != nil {
			t.Errorf("CreateUser() error = %v", err)
			return
		}
		if err := New().CreateUser(testReceiver); err != nil {
			t.Errorf("CreateUser() error = %v", err)
			return
		}
	})

	t.Run("getUserById", func(t *testing.T) {
		got, err := New().GetUserById(testUser.UserId)
		if err != nil {
			t.Errorf("GetUserById() error = %v", err)
			return
		}
		if !reflect.DeepEqual(testUser.UserId, got.UserId) {
			t.Errorf("GetUserById() gotData = %v, want %v", got, testUser)
			return
		}
	})

	t.Run("getUserBasicInfo", func(t *testing.T) {
		got, err := New().GetUserBasicInfo(testUser.UserId)
		if err != nil {
			t.Errorf("GetUserBasicInfo() error = %v", err)
			return
		}
		if !reflect.DeepEqual(testUser.BasicInfo.NickName, got.NickName) {
			t.Errorf("GetUserById() gotData = %v, want %v", got, testUser)
		}
	})

	var senderWalletData, receiverWalletData structs.WalletData
	t.Run("getWallet", func(t *testing.T) {
		_, err := New().GetWallet(testUser.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}
		senderWalletData, err = New().GetWallet(testSender.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}
		receiverWalletData, err = New().GetWallet(testReceiver.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}
	})

	t.Run("deposit", func(t *testing.T) {
		data := DepositData{
			DepositReq: apiBody.DepositReq{
				WalletId: senderWalletData.Id,
				Money:    testMoneyNtd,
			},
			CreatedAt: time.Now(),
		}
		data.DepositId = uuid.New().String()
		if err := New().Deposit(data); err != nil {
			t.Errorf("Deposit() error = %v", err)
			return
		}
		data.DepositId = uuid.New().String()
		if err := New().Deposit(data); err != nil {
			t.Errorf("Deposit() error = %v", err)
			return
		}
		data.Money = testMoneyUsd
		data.DepositId = uuid.New().String()
		if err := New().Deposit(data); err != nil {
			t.Errorf("Deposit() error = %v", err)
			return
		}
		data.DepositId = uuid.New().String()
		if err := New().Deposit(data); err != nil {
			t.Errorf("Deposit() error = %v", err)
			return
		}
	})

	t.Run("verify_deposit", func(t *testing.T) {
		walletData, err := New().GetWallet(testSender.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %s", err)
			return
		}
		expectedMoney := []structs.Money{
			{
				Currency: structs.CurrencyNTD,
				Amount:   200,
				Unit:     structs.UnitPico,
			},
			{
				Currency: structs.CurrencyUSD,
				Amount:   200,
				Unit:     structs.UnitPico,
			},
		}
		sort.Slice(expectedMoney, func(i, j int) bool {
			return strings.Compare(string(expectedMoney[i].Currency), string(expectedMoney[j].Currency)) == 1
		})
		sort.Slice(walletData.Balance, func(i, j int) bool {
			return strings.Compare(string(walletData.Balance[i].Currency), string(walletData.Balance[j].Currency)) == 1
		})
		if !reflect.DeepEqual(expectedMoney, walletData.Balance) {
			t.Errorf("unexpected wallet balance: %v", walletData.Balance)
			return
		}
	})

	t.Run("withdraw_normal", func(t *testing.T) {
		data := WithdrawData{
			WithdrawReq: WithdrawReq{
				Target:   "fakeTarget",
				WalletId: senderWalletData.Id,
				Money:    testMoneyNtd,
			},
			WithdrawId: uuid.New().String(),
		}
		if err := New().Withdraw(data); err != nil {
			t.Errorf("Withdraw() error = %v", err)
			return
		}
		data.Money = testMoneyUsd
		if err := New().Withdraw(data); err != nil {
			t.Errorf("Withdraw() error = %v", err)
			return
		}
	})

	t.Run("verify_withdraw", func(t *testing.T) {
		walletData, err := New().GetWallet(testSender.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}

		expectedMoney := []structs.Money{
			{
				Currency: structs.CurrencyNTD,
				Amount:   100,
				Unit:     structs.UnitPico,
			},
			{
				Currency: structs.CurrencyUSD,
				Amount:   100,
				Unit:     structs.UnitPico,
			},
		}
		sort.Slice(expectedMoney, func(i, j int) bool {
			return strings.Compare(string(expectedMoney[i].Currency), string(expectedMoney[j].Currency)) == 1
		})
		sort.Slice(walletData.Balance, func(i, j int) bool {
			return strings.Compare(string(walletData.Balance[i].Currency), string(walletData.Balance[j].Currency)) == 1
		})
		if !reflect.DeepEqual(expectedMoney, walletData.Balance) {
			t.Errorf("unexpected wallet balance: %v", walletData.Balance)
			return
		}
	})

	t.Run("withdraw_negative", func(t *testing.T) {
		data := WithdrawData{
			WithdrawReq: WithdrawReq{
				Target:   "fakeTarget",
				WalletId: senderWalletData.Id,
				Money:    testMoneyUsd,
			},
			WithdrawId: uuid.New().String(),
		}
		if err := New().Withdraw(data); err != nil {
			t.Errorf("Withdraw() error = %v", err)
			return
		}
		if err := New().Withdraw(data); err != gorm.ErrRecordNotFound {
			t.Errorf("Withdraw() want: %v, got: %v", gorm.ErrRecordNotFound, err)
			return
		}
	})

	testTransferData := TransactionData{
		TransactionReq: apiBody.TransactionReq{
			FromWalletId: senderWalletData.Id,
			ToWalletId:   receiverWalletData.Id,
			Money:        testMoneyNtd,
		},
		TransactionId: "testTransactionId",
	}

	t.Run("transfer", func(t *testing.T) {
		if err := New().Transfer(testTransferData); err != nil {
			t.Errorf("Transfer() error = %v", err)
			return
		}
	})

	t.Run("verify_transfer_sender", func(t *testing.T) {
		walletData, err := New().GetWallet(testSender.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}
		expectedMoney := []structs.Money{
			{
				Currency: structs.CurrencyNTD,
				Amount:   0,
				Unit:     structs.UnitPico,
			},
			{
				Currency: structs.CurrencyUSD,
				Amount:   0,
				Unit:     structs.UnitPico,
			},
		}
		sort.Slice(expectedMoney, func(i, j int) bool {
			return strings.Compare(string(expectedMoney[i].Currency), string(expectedMoney[j].Currency)) == 1
		})
		sort.Slice(walletData.Balance, func(i, j int) bool {
			return strings.Compare(string(walletData.Balance[i].Currency), string(walletData.Balance[j].Currency)) == 1
		})
		if !reflect.DeepEqual(expectedMoney, walletData.Balance) {
			t.Errorf("unexpected wallet balance: %v", walletData.Balance)
			return
		}
	})

	t.Run("verify_transfer_receiver", func(t *testing.T) {
		walletData, err := New().GetWallet(testReceiver.UserId)
		if err != nil {
			t.Errorf("GetWallet() error = %v", err)
			return
		}
		expectedMoney := []structs.Money{
			{
				Currency: structs.CurrencyNTD,
				Amount:   100,
				Unit:     structs.UnitPico,
			},
		}
		sort.Slice(expectedMoney, func(i, j int) bool {
			return strings.Compare(string(expectedMoney[i].Currency), string(expectedMoney[j].Currency)) == 1
		})
		sort.Slice(walletData.Balance, func(i, j int) bool {
			return strings.Compare(string(walletData.Balance[i].Currency), string(walletData.Balance[j].Currency)) == 1
		})
		if !reflect.DeepEqual(expectedMoney, walletData.Balance) {
			t.Errorf("unexpected wallet balance: %v", walletData.Balance)
			return
		}
	})

	t.Run("getTransactions", func(t *testing.T) {
		data, err := New().GetTransactions(testSender.UserId)
		if err != nil {
			t.Errorf("GetTransactions() error = %v", err)
			return
		}
		if data == nil {
			t.Errorf("GetTransactions() empty data")
			return
		}
		t.Log(data)
	})
}

//func Test_t(t *testing.T) {
//	sqlDb, err := gorm.Open(mysql.Open(config.New().MySql.ConnectionString), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//	}
//	//b453e4d7-bf09-4923-b7a0-3143054a0a83
//	w := structs.Wallet{
//		Id: "55f0d350-6cb3-4397-ae3b-faee0eab6470",
//	}
//	sqlDb.Preload("Balance").Find(&w)
//
//	_ = w
//}

func teardown(db string) {
	_ = instance.mongo.client.Database(db).Drop(context.Background())
	instance.mysql.db.Model(&structs.Wallet{}).Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", db))
	//instance.mysql.db.Delete(&structs.Balance{})
}
