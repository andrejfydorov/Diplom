package voice

import (
	"Diplom/internal/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type VoiceCallData struct {
	Сountry            string  `json:"сountry"`
	Bandwidth          int     `json:"bandwidth"`
	ResponseTime       int     `json:"response_time"`
	Provider           string  `json:"provider"`
	ConnStability      float64 `json:"conn_stability"`
	PurityOfTTFB       int     `json:"purity_of_ttfb"`
	MedianCallDuration int     `json:"median_call_duration"`
}

var providers = []string{"TransparentCalls", "E-Voice", "JustPhone"}

type Repo struct {
	mutex      sync.Mutex
	voiceCalls []*VoiceCallData
}

type VoiceCallService interface {
	ReplaceCountries()
	SortWithCountry()
	SortWithProvider()
	GetData() ([]VoiceCallData, error)
	PrintData()
}

func New() (VoiceCallService, error) {
	var r = Repo{}
	err := r.LoadData()
	if err != nil {
		return nil, errors.New("voice call service failed")
	}
	return &r, nil
}

func (r *Repo) GetData() ([]VoiceCallData, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var res = make([]VoiceCallData, len(r.voiceCalls))

	for i, vc := range r.voiceCalls {
		res[i] = *vc
	}

	if len(res) == 0 {
		return nil, errors.New("voice call service failed")
	}

	return res, nil
}

func (r *Repo) PrintData() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, vc := range r.voiceCalls {
		log.Println(vc)
	}
}

func (r *Repo) SortWithCountry() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.voiceCalls)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.voiceCalls[j].Сountry, r.voiceCalls[minIdx].Сountry) == -1 {
				minIdx = j
			}
		}
		r.voiceCalls[i], r.voiceCalls[minIdx] = r.voiceCalls[minIdx], r.voiceCalls[i]
	}
}

func (r *Repo) SortWithProvider() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.voiceCalls)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(r.voiceCalls[j].Provider, r.voiceCalls[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		r.voiceCalls[i], r.voiceCalls[minIdx] = r.voiceCalls[minIdx], r.voiceCalls[i]
	}
}

func (r *Repo) ReplaceCountries() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	size := len(r.voiceCalls)
	for i := 0; i < size; i++ {
		r.voiceCalls[i].Сountry = utils.Alpha_2[r.voiceCalls[i].Сountry]
	}
}

func (r *Repo) LoadData() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.Open("resources/voice.data")
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

		if len(str) < 7 {
			continue
		}

		if ok := utils.Contains(&providers, str[3]); !ok {
			continue
		}

		if _, ok := utils.Alpha_2[str[0]]; !ok {
			continue
		}

		var vc VoiceCallData

		vc.Сountry = str[0]

		i, err := strconv.Atoi(str[1])
		if err != nil {
			return err
		}
		vc.Bandwidth = i

		i, err = strconv.Atoi(str[2])
		if err != nil {
			return err
		}
		vc.ResponseTime = i

		vc.Provider = str[3]

		f, err := strconv.ParseFloat(str[4], 64)
		if err != nil {
			return err
		}
		vc.ConnStability = f

		i, err = strconv.Atoi(str[5])
		if err != nil {
			return err
		}
		vc.PurityOfTTFB = i

		i, err = strconv.Atoi(str[6])
		if err != nil {
			return err
		}
		vc.MedianCallDuration = i

		r.voiceCalls = append(r.voiceCalls, &vc)

	}

	if len(r.voiceCalls) == 0 {
		return errors.New("voice call service failed")
	}

	return nil
}
