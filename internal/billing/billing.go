package billing

import (
	"Diplom/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

type Repo struct {
	mutex  sync.Mutex
	bilngs []*BillingData
}

func (r *Repo) Open() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.Open("resources/billing.data")
	if err != nil {
		log.Println("Unable to open file:", err)
		log.Fatalln(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	//fmt.Println(lines)

MainLoop:
	for _, line := range lines {

		if len(line) < 6 {
			continue
		}

		for _, v := range line {
			if !unicode.IsDigit(v) {
				continue MainLoop
			}
		}

		var bytes [6]byte

		for i, v := range line {
			num, err := strconv.Atoi(string(v))
			if err != nil {
				log.Fatalln(err)
			}
			bytes[i] = byte(num)

			//bol, err := strconv.ParseBool(string(v))
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//bools = append(bools, bol)

		}

		bytesSum := utils.BitsToUint8(bytes)

		//bytes = utils.BitsToBytes(bytes)

		var billing BillingData

		billing.CreateCustomer = utils.Itob((bytesSum >> 0) & 1 & bytes[5])
		billing.Purchase = utils.Itob((bytesSum >> 1) & 1 & bytes[4])
		billing.Payout = utils.Itob((bytesSum >> 2) & 1 & bytes[3])
		billing.Recurring = utils.Itob((bytesSum >> 3) & 1 & bytes[2])
		billing.FraudControl = utils.Itob((bytesSum >> 4) & 1 & bytes[1])
		billing.CheckoutPage = utils.Itob((bytesSum >> 5) & 1 & bytes[0])

		r.bilngs = append(r.bilngs, &billing)

	}

	for _, bill := range r.bilngs {
		log.Println(bill)
	}

}
