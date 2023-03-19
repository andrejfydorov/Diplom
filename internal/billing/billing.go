package billing

import (
	"Diplom/internal/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"unicode"
)

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

var billing *BillingData

func GetData() (BillingData, error) {
	res := billing

	if res == nil {
		var r BillingData
		return r, errors.New("billing service failed")
	}

	return *res, nil
}

func PrintData() {
	fmt.Println(billing)
}

func LoadData() error {
	file, err := os.Open("simulator/billing.data")
	if err != nil {
		log.Println("Unable to open file:", err)
		log.Println(err)
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}

	contentStr := string(content)

	if len(contentStr) < 6 {
		log.Println("Invalid input string")
		return errors.New("Invalid input string")
	}

	for _, v := range contentStr {
		if !unicode.IsDigit(v) {
			log.Println("Invalid input string")
			return errors.New("Invalid input string")
		}
	}

	var bytes [6]byte

	for i, v := range contentStr {
		num, err := strconv.Atoi(string(v))
		if err != nil {
			log.Println(err)
			return err
		}
		bytes[i] = byte(num)
	}

	billing = new(BillingData)

	bytesSum := utils.BitsToUint8(bytes)

	billing.CreateCustomer = utils.Itob((bytesSum >> 0) & 1 & bytes[5])
	billing.Purchase = utils.Itob((bytesSum >> 1) & 1 & bytes[4])
	billing.Payout = utils.Itob((bytesSum >> 2) & 1 & bytes[3])
	billing.Recurring = utils.Itob((bytesSum >> 3) & 1 & bytes[2])
	billing.FraudControl = utils.Itob((bytesSum >> 4) & 1 & bytes[1])
	billing.CheckoutPage = utils.Itob((bytesSum >> 5) & 1 & bytes[0])

	return nil
}
