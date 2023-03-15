package email

import (
	"Diplom/internal/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type EmailData struct {
	Country      string `json:"country,omitempty"`
	Provider     string `json:"provider,omitempty"`
	DeliveryTime int    `json:"deliveryTime,omitempty"`
}

var providers = []string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast",
	"AOL", "Live", "RediffMail", "GMX", "Proton Mail", "Yandex", "Mail.ru"}

type Repo struct {
	mutex  sync.Mutex
	emails []*EmailData
}

type EmailService interface {
	ReplaceCountries()
	SortWithCountry()
	SortWithProvider()
	GetData() ([]EmailData, error)
	GetThreeFast() map[string][]EmailData
	GetThreeSlow() map[string][]EmailData
	PrintData()
}

func New() (EmailService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("email service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]EmailData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]EmailData, len(r.emails))

	for i, email := range r.emails {
		res[i] = *email
	}

	if len(res) == 0 {
		return nil, errors.New("email service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, email := range r.emails {
		log.Println(email)
	}
}

func (r *Repo) SortWithCountry() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.emails)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.emails[j].Country, r.emails[minIdx].Country) == -1 {
				minIdx = j
			}
		}
		r.emails[i], r.emails[minIdx] = r.emails[minIdx], r.emails[i]
	}
}

func (r *Repo) SortWithProvider() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.emails)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.emails[j].Provider, r.emails[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		r.emails[i], r.emails[minIdx] = r.emails[minIdx], r.emails[i]
	}
}

func (r *Repo) ReplaceCountries() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.emails)
	for i := 0; i < size; i++ {
		r.emails[i].Country = utils.Alpha_2[r.emails[i].Country]
	}
}

func (r Repo) GetThreeFast() map[string][]EmailData {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make(map[string][]EmailData)

	var email1 *EmailData = r.emails[0]
	var email2 *EmailData = r.emails[1]
	var email3 *EmailData = r.emails[2]

	for i, email := range r.emails {
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
			email2 = r.emails[i+1]
			email3 = r.emails[i+2]
		}
	}

	return res
}

func (r Repo) GetThreeSlow() map[string][]EmailData {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make(map[string][]EmailData)

	var email1 *EmailData = r.emails[0]
	var email2 *EmailData = r.emails[0]
	var email3 *EmailData = r.emails[0]

	for _, email := range r.emails {
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

func (r *Repo) LoadData() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.Open("resources/email.data")
	if err != nil {
		log.Println("Unable to open file:", err)
		log.Fatalln(err)
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
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
			log.Fatalln(err)
			return err
		}
		email.DeliveryTime = i

		r.emails = append(r.emails, &email)

	}

	return nil
}
