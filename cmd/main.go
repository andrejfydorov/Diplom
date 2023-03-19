package main

import (
	"Diplom/internal/result"
	"fmt"
)

func main() {
	//sms.LoadData()
	//sms.ReplaceCountries()
	//sms.SortWithCountry()
	//sms.PrintData()

	//mms.LoadData()
	//mms.ReplaceCountries()
	//mms.SortWithCountry()
	//mms.PrintData()

	//email.LoadData()
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

	//voice.LoadData()
	//voice.ReplaceCountries()
	//voice.SortWithCountry()
	//voice.PrintData()

	//billing.LoadData()
	//billing.PrintData()

	//support.LoadData()
	//support.PrintData()
	//s := support.GetCalculatedData()
	//for _, i := range s {
	//	fmt.Print(i, "  ")
	//}
	//fmt.Print("\n")

	//incident.LoadData()
	//incident.PrintData()

	result := result.GetResultData()
	fmt.Println(*result)

	//service.NewService()

}
