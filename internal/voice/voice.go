package voice

import (
	"Diplom/internal/utils"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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

var voiceCalls []*VoiceCallData

func GetData() ([]VoiceCallData, error) {
	var res = make([]VoiceCallData, len(voiceCalls))

	for i, vc := range voiceCalls {
		res[i] = *vc
	}

	if len(res) == 0 {
		return nil, errors.New("voice call service failed")
	}

	return res, nil
}

func PrintData() {
	for _, vc := range voiceCalls {
		log.Println(vc)
	}
}

func SortWithCountry() {
	size := len(voiceCalls)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(voiceCalls[j].Сountry, voiceCalls[minIdx].Сountry) == -1 {
				minIdx = j
			}
		}
		voiceCalls[i], voiceCalls[minIdx] = voiceCalls[minIdx], voiceCalls[i]
	}
}

func SortWithProvider() {
	size := len(voiceCalls)
	for i := 0; i < size-1; i++ {
		var minIdx = i
		for j := i; j < size; j++ {
			if strings.Compare(voiceCalls[j].Provider, voiceCalls[minIdx].Provider) == -1 {
				minIdx = j
			}
		}
		voiceCalls[i], voiceCalls[minIdx] = voiceCalls[minIdx], voiceCalls[i]
	}
}

func ReplaceCountries() {
	size := len(voiceCalls)
	for i := 0; i < size; i++ {
		voiceCalls[i].Сountry = utils.Alpha_2[voiceCalls[i].Сountry]
	}
}

func LoadData() error {
	file, err := os.Open("simulator/voice.data")
	if err != nil {
		log.Println("Unable to open file:", err)
		log.Println(err)
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
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

		voiceCalls = append(voiceCalls, &vc)

	}

	if len(voiceCalls) == 0 {
		return errors.New("voice call service failed")
	}

	return nil
}
