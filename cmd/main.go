package main

import "Diplom/internal/incident"

func main() {
	//sms, _ := sms.New()
	//sms.ReplaceCountries()
	//sms.SortWithCountry()
	//sms.PrintData()

	//mms, _ := mms.New()
	//mms.ReplaceCountries()
	//mms.SortWithCountry()
	//mms.PrintData()

	//email, _ := email.New()
	//email.PrintData()
	//fmt.Println("==================================\n")
	//emails := email.GetThreeFast()
	//for _, e := range emails {
	//	fmt.Println(e)
	//}
	//fmt.Println("==================================\n")
	//emails = email.GetThreeSlow()
	//for _, e := range emails {
	//	fmt.Println(e)
	//}

	//vc, _ := voice.New()
	//vc.ReplaceCountries()
	//vc.SortWithCountry()
	//vc.PrintData()

	//bill, _ := billing.New()
	//bill.PrintData()

	//support := &support.Repo{}
	//support.Open()

	inc, _ := incident.New()
	inc.PrintData()

	//service.NewService()

}
