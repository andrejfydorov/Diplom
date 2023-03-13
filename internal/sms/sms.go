package sms

import (
	"Diplom/internal/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type SMSData struct {
	Сountry      string `json:"сountry"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

var providers = []string{"Topolo", "Rond", "Kildy"}

type Repo struct {
	mutex sync.Mutex
	smses []*SMSData
}

type SmsService interface {
	ReplaceCountries()
	SortWithCountry()
	SortWithProvider()
	GetData() ([]SMSData, error)
	PrintData()
}

func New() (SmsService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("sms service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]SMSData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]SMSData, len(r.smses))

	for i, sms := range r.smses {
		res[i] = *sms
	}

	if len(res) == 0 {
		return nil, errors.New("sms service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, sms := range r.smses {
		log.Println(sms)
	}
}

func (r *Repo) SortWithCountry() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.smses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.smses[j].Сountry, r.smses[minIdx].Сountry) == -1 {
				minIdx = j
			}
		}
		r.smses[i], r.smses[minIdx] = r.smses[minIdx], r.smses[i]
	}
}

func (r *Repo) SortWithProvider() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.smses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.smses[j].Provider, r.smses[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		r.smses[i], r.smses[minIdx] = r.smses[minIdx], r.smses[i]
	}
}

func (r *Repo) ReplaceCountries() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.smses)
	for i := 0; i < size; i++ {
		r.smses[i].Сountry = utils.Alpha_2[r.smses[i].Сountry]
	}
}

func (r *Repo) LoadData() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.Open("resources/sms.data")
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

		if len(str) < 4 {
			continue
		}

		if _, ok := utils.Alpha_2[str[0]]; !ok {
			continue
		}

		if ok := utils.Contains(&providers, str[3]); !ok {
			continue
		}

		var sms SMSData

		sms.Сountry = str[0]
		sms.Bandwidth = str[1]
		sms.ResponseTime = str[2]
		sms.Provider = str[3]

		r.smses = append(r.smses, &sms)

	}

	if len(r.smses) == 0 {
		return errors.New("sms service failed")
	}

	return nil
}
