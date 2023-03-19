package email

import (
	"Diplom/internal/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type EmailData struct {
	Country      string `json:"country,omitempty"`
	Provider     string `json:"provider,omitempty"`
	DeliveryTime int    `json:"deliveryTime,omitempty"`
}

var providers = []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast",
	"AOL", "Live", "RediffMail", "GMX", "Proton Mail", "Yandex", "Mail.ru"}

var emails []*EmailData

func GetData() ([]EmailData, error) {
	var res = make([]EmailData, len(emails))

	for i, email := range emails {
		res[i] = *email
	}

	if len(res) == 0 {
		return nil, errors.New("email service failed")
	}

	return res, nil
}

func PrintData() {
	for _, email := range emails {
		log.Println(email)
	}
}

func SortWithCountry() {
	size := len(emails)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(emails[j].Country, emails[minIdx].Country) == -1 {
				minIdx = j
			}
		}
		emails[i], emails[minIdx] = emails[minIdx], emails[i]
	}
}

func SortWithProvider() {
	size := len(emails)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(emails[j].Provider, emails[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		emails[i], emails[minIdx] = emails[minIdx], emails[i]
	}
}

func ReplaceCountries() {
	size := len(emails)
	for i := 0; i < size; i++ {
		emails[i].Country = utils.Alpha_2[emails[i].Country]
	}
}

func GetThreeFast() map[string][]EmailData {
	var res = make(map[string][]EmailData)

	var email1 *EmailData = emails[0]
	var email2 *EmailData = emails[1]
	var email3 *EmailData = emails[2]

	for i, email := range emails {
		if email.Country == email1.Country {
			if email.DeliveryTime > email1.DeliveryTime {
				email1 = email
			} else if email.DeliveryTime > email2.DeliveryTime {
				email2 = email
			} else if email.DeliveryTime > email3.DeliveryTime {
				email3 = email
			}
		} else {
			res[email1.Country] = []EmailData{*email1, *email2, *email3}
			email1 = email
			email2 = emails[i+1]
			email3 = emails[i+2]
		}
	}

	return res
}

func GetThreeSlow() map[string][]EmailData {
	var res = make(map[string][]EmailData)

	var email1 *EmailData = emails[0]
	var email2 *EmailData = emails[0]
	var email3 *EmailData = emails[0]

	for _, email := range emails {
		if email.Country == email1.Country {
			if email.DeliveryTime < email1.DeliveryTime {
				email1 = email
			} else if email.DeliveryTime < email2.DeliveryTime {
				email2 = email
			} else if email.DeliveryTime < email3.DeliveryTime {
				email3 = email
			}
		} else {
			res[email1.Country] = []EmailData{*email1, *email2, *email3}
			email1 = email
			email2 = email
			email3 = email
		}
	}

	return res
}

func LoadData() error {
	file, err := os.Open("simulator/email.data")
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
	lines := strings.Split(contentStr, "\n")

	for _, line := range lines {

		str := strings.Split(line, ";")

		if len(str) < 3 {
			continue
		}

		if ok := utils.Contains(&providers, str[1]); !ok {
			continue
		}

		if _, ok := utils.Alpha_2[str[0]]; !ok {
			continue
		}

		var email EmailData

		email.Country = str[0]
		email.Provider = str[1]

		i, err := strconv.Atoi(str[2])
		if err != nil {
			log.Println(err)
			return err
		}
		email.DeliveryTime = i

		emails = append(emails, &email)

	}

	return nil
}
