package support

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

var support_structs []*SupportData

func GetData() ([]SupportData, error) {
	var res = make([]SupportData, len(support_structs))

	for i, sup := range support_structs {
		res[i] = *sup
	}

	if len(res) == 0 {
		return nil, errors.New("Support service failed")
	}

	return res, nil
}

func PrintData() {
	for _, sup := range support_structs {
		log.Println(sup)
	}
}

func GetCalculatedData() []int {
	res := make([]int, 2)

	var openTicketsCount int = 0

	for _, sup := range support_structs {
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

func LoadData() error {
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

		err = json.Unmarshal(body, &support_structs)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
