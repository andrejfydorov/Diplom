package support

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

type Repo struct {
	mutex           sync.Mutex
	support_structs []*SupportData
}

type SupportService interface {
	GetData() ([]SupportData, error)
	GetCalculatedData() []int
	PrintData()
}

func New() (SupportService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("Support service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]SupportData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]SupportData, len(r.support_structs))

	for i, sup := range r.support_structs {
		res[i] = *sup
	}

	if len(res) == 0 {
		return nil, errors.New("Support service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, sup := range r.support_structs {
		log.Println(sup)
	}
}

func (r *Repo) GetCalculatedData() []int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	res := make([]int, 2)

	var openTicketsCount int = 0

	for _, sup := range r.support_structs {
		openTicketsCount += sup.ActiveTickets
	}

	res[1] = (60 / 18) * openTicketsCount

	if openTicketsCount < 9 {
		res[0] = 1
	} else if openTicketsCount >= 9 && openTicketsCount <= 16 {
		res[0] = 2
	} else {
		res[0] = 3
	}

	return res
}

func (r *Repo) LoadData() error {
	resp, err := http.Get("http://127.0.0.1:8383/support") //127.0.0.1:8484
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

		var support []*SupportData

		err = json.Unmarshal(body, &support)
		if err != nil {
			log.Println(err)
			return err
		}

		r.support_structs = support
	}

	return nil
}
