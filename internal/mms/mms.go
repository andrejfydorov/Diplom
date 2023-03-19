package mms

import (
	"Diplom/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

type Repo struct {
	mutex sync.Mutex
	mmses []*MMSData
}

type MMSService interface {
	ReplaceCountries()
	SortWithCountry()
	SortWithProvider()
	GetData() ([]MMSData, error)
	PrintData()
}

var providers = []string{"Topolo", "Rond", "Kildy"}

func New() (MMSService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("mms service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]MMSData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]MMSData, len(r.mmses))

	for i, mms := range r.mmses {
		res[i] = *mms
	}

	if len(res) == 0 {
		return nil, errors.New("mms service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, mms := range r.mmses {
		log.Println(mms)
	}
}

func (r *Repo) SortWithCountry() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.mmses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.mmses[j].Country, r.mmses[minIdx].Country) == -1 {
				minIdx = j
			}
		}
		r.mmses[i], r.mmses[minIdx] = r.mmses[minIdx], r.mmses[i]
	}
}

func (r *Repo) SortWithProvider() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.mmses)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.mmses[j].Provider, r.mmses[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		r.mmses[i], r.mmses[minIdx] = r.mmses[minIdx], r.mmses[i]
	}
}

func (r *Repo) ReplaceCountries() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.mmses)
	for i := 0; i < size; i++ {
		r.mmses[i].Country = utils.Alpha_2[r.mmses[i].Country]
	}
}

func (r *Repo) LoadData() error {
	resp, err := http.Get("http://127.0.0.1:8383/mms") //127.0.0.1:8383
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		r.mutex.Lock()
		defer r.mutex.Unlock()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		var mmses []*MMSData

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

		r.mmses = mmses

	}

	if len(r.mmses) == 0 {
		return errors.New("mms service failed")
	}

	return nil
}
