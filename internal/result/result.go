package result

import (
	"Diplom/internal/billing"
	"Diplom/internal/email"
	"Diplom/internal/incident"
	"Diplom/internal/mms"
	"Diplom/internal/sms"
	"Diplom/internal/voice"
	"log"
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

func GetResultData() (ResultSetT, ResultT) {
	var rst ResultSetT
	var rt = ResultT{Status: true, Data: rst}

	//========start sms region===================
	sms, err := sms.New()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}
	sms.ReplaceCountries()
	sms.SortWithCountry()
	sms1, err := sms.GetData()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}

	sms.SortWithProvider()
	sms2, err := sms.GetData()
	if err != nil {
		rt.Error = err.Error()
	}
	//========end sms region===================

	//========start mms region===================
	mms, err := mms.New()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}
	mms.ReplaceCountries()
	mms.SortWithCountry()
	mms1, err := mms.GetData()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}

	mms.SortWithProvider()
	mms2, err := mms.GetData()
	if err != nil {
		rt.Error = err.Error()
	}
	//========end mms region===================
	//======start voice call region=========
	vc, err := voice.New()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}

	vc1, err := vc.GetData()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}

	//======end voice call region=========
	//======start email region=========
	email, err := email.New()
	if err != nil {
		log.Println(err)
		rt.Status = false
		rt.Error = err.Error()
		return rst, rt
	}

	max_emails := email.GetThreeFast()
	mini_emails := email.GetThreeSlow()

	var emails map[string][][]email.EmailData
	for i, _ := range max_emails {
		emails[i][0] = max_emails[i]
		emails[i][1] = mini_emails[i]
	}
	if len(emails) == 0 {
		log.Println(err)
		rt.Status = false
		rt.Error = "email service failed"
		return rst, rt
	}

	//======end email region=========
	//========== assembling ResultSetT ResultT structures ================

	rst.SMS[0] = sms2
	rst.SMS[1] = sms1

	rst.MMS[0] = mms2
	rst.MMS[1] = mms1

	rst.VoiceCall = vc1

	rst.Email = emails

	//============== return structures=============

	return rst, rt

}
