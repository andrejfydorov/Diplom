package sms

import (
	"Diplom/internal/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type SMSData struct {
	Сountry      string `json:"сountry"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

var providers = []string{"Topolo", "Rond", "Kildy"}

var smses []*SMSData

func GetData() ([]SMSData, error) {
	var res = make([]SMSData, len(smses))

	for i, sms := range smses {
		res[i] = *sms
	}

	if len(res) == 0 {
		return nil, errors.New("sms service failed")
	}

	return res, nil
}

func PrintData() {
	for _, sms := range smses {
		log.Println(sms)
	}
}

func SortWithCountry() {
	size := len(smses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(smses[j].Сountry, smses[minIdx].Сountry) == -1 {
				minIdx = j
			}
		}
		smses[i], smses[minIdx] = smses[minIdx], smses[i]
	}
}

func SortWithProvider() {
	size := len(smses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(smses[j].Provider, smses[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		smses[i], smses[minIdx] = smses[minIdx], smses[i]
	}
}

func ReplaceCountries() {
	size := len(smses)
	for i := 0; i < size; i++ {
		smses[i].Сountry = utils.Alpha_2[smses[i].Сountry]
	}
}

func LoadData() error {
	file, err := os.Open("simulator/sms.data")
	if err != nil {
		fmt.Println("Unable to open file:", err)
		log.Println(err)
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

		smses = append(smses, &sms)

	}

	if len(smses) == 0 {
		return errors.New("sms service failed")
	}

	return nil
}
