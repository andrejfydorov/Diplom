package incident

import (
	"encoding/json"
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
	mutex sync.Mutex
	mmses []*IncidentData
}

func (r *Repo) Open() {
	resp, err := http.Get("http://127.0.0.1:8383/accendent") //127.0.0.1:8585
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		var incident = []IncidentData{}

		err = json.Unmarshal(body, &incident)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(incident)

	}

}
