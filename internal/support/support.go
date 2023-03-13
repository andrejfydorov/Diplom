package support

import (
	"encoding/json"
	"fmt"
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
	mutex sync.Mutex
	mmses []*SupportData
}

func (r *Repo) Open() {
	resp, err := http.Get("http://127.0.0.1:8383/support") //127.0.0.1:8484
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}

		var support = []SupportData{}

		err = json.Unmarshal(body, &support)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(support)

	}

}
