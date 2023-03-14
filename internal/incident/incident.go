package incident

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы active и closed
}

type Repo struct {
	mutex     sync.Mutex
	incidents []*IncidentData
}

type IncidentService interface {
	SortWithStatus()
	GetData() ([]IncidentData, error)
	PrintData()
}

func New() (IncidentService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("Incident service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]IncidentData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]IncidentData, len(r.incidents))

	for i, inc := range r.incidents {
		res[i] = *inc
	}

	if len(res) == 0 {
		return nil, errors.New("Incident service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, inc := range r.incidents {
		log.Println(inc)
	}
}

func (r *Repo) SortWithStatus() {
	len := len(r.incidents)
	for i := 0; i < len-1; i++ {
		for j := 1; j < len; j++ {
			if r.incidents[i].Status != "active" && r.incidents[j].Status == "active" {
				r.incidents[i], r.incidents[j] = r.incidents[j], r.incidents[i]
			}
		}
	}
}

func (r *Repo) LoadData() error {
	resp, err := http.Get("http://127.0.0.1:8383/accendent") //127.0.0.1:8585
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

		var incident = []IncidentData{}

		err = json.Unmarshal(body, &incident)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Println(incident)
	}

	return nil
}
