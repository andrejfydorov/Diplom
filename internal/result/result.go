package result

import (
	"Diplom/internal/billing"
	"Diplom/internal/email"
	"Diplom/internal/incident"
	"Diplom/internal/mms"
	"Diplom/internal/sms"
	"Diplom/internal/support"
	"Diplom/internal/voice"
	"errors"
	"log"
	"sync"
)

type ResultT struct {
	Status bool       `json:"status"` // True, если все этапы сбора данных прошли успешно, False во всех остальных случаях
	Data   ResultSetT `json:"data"`   // Заполнен, если все этапы сбора  данных прошли успешно, nil во всех остальных случаях
	Error  string     `json:"error"`  // Пустая строка, если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки
}

type ResultSetT struct {
	SMS       [][]sms.SMSData                `json:"sms"`
	MMS       [][]mms.MMSData                `json:"mms"`
	VoiceCall []voice.VoiceCallData          `json:"voice_call"`
	Email     map[string][][]email.EmailData `json:"email"`
	Billing   billing.BillingData            `json:"billing"`
	Support   []int                          `json:"support"`
	Incidents []incident.IncidentData        `json:"incident"`
}

var mutex sync.Mutex

func assemblingSms() (*[][]sms.SMSData, error) {
	var smses = make([][]sms.SMSData, 2)
	smses[0] = make([]sms.SMSData, 1)
	smses[1] = make([]sms.SMSData, 1)

	err := sms.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sms.ReplaceCountries()
	sms.SortWithCountry()

	sms1, err := sms.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sms.SortWithProvider()

	sms2, err := sms.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	smses[0] = sms1
	smses[0] = sms2

	if len(smses) == 0 {
		log.Println("sms service failed")
		return nil, errors.New("sms service failed")
	}

	return &smses, nil
}

func assemblingMms() (*[][]mms.MMSData, error) {
	var mmses = make([][]mms.MMSData, 2)
	mmses[0] = make([]mms.MMSData, 1)
	mmses[1] = make([]mms.MMSData, 1)

	err := mms.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	mms.ReplaceCountries()
	mms.SortWithCountry()

	mms1, err := mms.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	mms.SortWithProvider()

	mms2, err := mms.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	mmses[0] = mms1
	mmses[0] = mms2

	if len(mmses) == 0 {
		log.Println("mms service failed")
		return nil, errors.New("mms service failed")
	}

	return &mmses, err
}

func asseblingVoiceCall() (*[]voice.VoiceCallData, error) {
	err := voice.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	vc, err := voice.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(vc) == 0 {
		log.Println("voice service failed")
		return nil, errors.New("voice service failed")
	}

	return &vc, nil
}

func assemblingEmail() (*map[string][][]email.EmailData, error) {
	var emails = make(map[string][][]email.EmailData)

	err := email.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	max_emails := email.GetThreeFast()
	min_emails := email.GetThreeSlow()

	for i, _ := range max_emails {
		emails[i] = make([][]email.EmailData, 2)
		emails[i][0] = make([]email.EmailData, 1)
		emails[i][1] = make([]email.EmailData, 1)

		emails[i][0] = max_emails[i]
		emails[i][1] = min_emails[i]
	}

	if len(emails) == 0 {
		log.Println("email service failed")
		return nil, errors.New("email service failed")
	}

	return &emails, nil
}

func assemblingBilling() (*billing.BillingData, error) {
	err := billing.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	billing, err := billing.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &billing, nil
}

func assemblingSupport() (*[]int, error) {
	err := support.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	supports := support.GetCalculatedData()

	if len(supports) == 0 {
		log.Println("support service failed")
		return nil, errors.New("support service failed")
	}

	return &supports, nil
}

func assemlingIncident() (*[]incident.IncidentData, error) {
	err := incident.LoadData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	incident.SortWithStatus()

	incidents, err := incident.GetData()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(incidents) == 0 {
		log.Println("incident service failed")
		return nil, errors.New("incident service failed")
	}

	return &incidents, nil
}

func GetResultData() *ResultT {
	mutex.Lock()
	defer mutex.Unlock()

	var rst ResultSetT
	var rt ResultT

	smses, err := assemblingSms()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	mmses, err := assemblingMms()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	voices, err := asseblingVoiceCall()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	ems, err := assemblingEmail()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	bil, err := assemblingBilling()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	sup, err := assemblingSupport()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	inc, err := assemlingIncident()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		rt.Data = rst
		return &rt
	}

	//if all ok return struct
	rst.SMS = *smses
	rst.MMS = *mmses
	rst.VoiceCall = *voices
	rst.Email = *ems
	rst.Billing = *bil
	rst.Support = *sup
	rst.Incidents = *inc

	rt.Status = true
	rt.Error = ""
	rt.Data = rst

	return &rt
}
