package mms

import (
	"Diplom/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

var mmses []*MMSData

var providers = []string{"Topolo", "Rond", "Kildy"}

func GetData() ([]MMSData, error) {
	var res = make([]MMSData, len(mmses))

	for i, mms := range mmses {
		res[i] = *mms
	}

	if len(res) == 0 {
		return nil, errors.New("mms service failed")
	}

	return res, nil
}

func PrintData() {
	for _, mms := range mmses {
		log.Println(mms)
	}
}

func SortWithCountry() {
	size := len(mmses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(mmses[j].Country, mmses[minIdx].Country) == -1 {
				minIdx = j
			}
		}
		mmses[i], mmses[minIdx] = mmses[minIdx], mmses[i]
	}
}

func SortWithProvider() {
	size := len(mmses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(mmses[j].Provider, mmses[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		mmses[i], mmses[minIdx] = mmses[minIdx], mmses[i]
	}
}

func ReplaceCountries() {
	size := len(mmses)
	for i := 0; i < size; i++ {
		mmses[i].Country = utils.Alpha_2[mmses[i].Country]
	}
}

func LoadData() error {
	resp, err := http.Get("http://127.0.0.1:8383/mms") //127.0.0.1:8383
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(body, &mmses)
		if err != nil {
			log.Println(err)
			return err
		}

		for i, mms := range mmses {
			ok1 := utils.Contains(&providers, mms.Provider)
			_, ok2 := utils.Alpha_2[mms.Country]

			if !ok2 || !ok1 {
				mmses = utils.Remove(mmses, i)
			}
		}

	}

	if len(mmses) == 0 {
		return errors.New("mms service failed")
	}

	return nil
}
