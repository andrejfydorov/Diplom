package incident

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` // возможные статусы active и closed
}

var incidents []*IncidentData

func GetData() ([]IncidentData, error) {
	var res = make([]IncidentData, len(incidents))

	for i, inc := range incidents {
		res[i] = *inc
	}

	if len(res) == 0 {
		return nil, errors.New("Incident service failed")
	}

	return res, nil
}

func PrintData() {
	for _, inc := range incidents {
		fmt.Println(inc)
	}
}

func SortWithStatus() {
	len := len(incidents)
	for i := 0; i < len-1; i++ {
		for j := 1; j < len; j++ {
			if incidents[i].Status != "active" && incidents[j].Status == "active" {
				incidents[i], incidents[j] = incidents[j], incidents[i]
			}
		}
	}
}

func LoadData() error {
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

		err = json.Unmarshal(body, &incidents)
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}
