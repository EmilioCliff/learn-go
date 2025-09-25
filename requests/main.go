package main

import (
	"fmt"
	"log"

	"github.com/EmilioCliff/learn-go/requests/myservice"
	"github.com/hooklift/gowsdl/soap"
)

func main() {
	countrySeviceClient := soap.NewClient("http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso")
	countryService := myservice.NewCountryInfoServiceSoapType(countrySeviceClient)

	capitalCityName, err := countryService.CapitalCity(&myservice.CapitalCity{SCountryISOCode: "KE"})
	if err != nil {
		log.Fatalf("Error calling CapitalCity: %v", err)
	}
	log.Println(capitalCityName.CapitalCityResult)

	fullCountryInfoRequest := &myservice.CountriesUsingCurrency{SISOCurrencyCode: "KES"}
	fullCountryInfoResponse, err := countryService.CountriesUsingCurrency(fullCountryInfoRequest)
	if err != nil {
		log.Fatalf("Error calling FullCountryInfo: %v", err)
	}

	for i := range fullCountryInfoResponse.CountriesUsingCurrencyResult.TCountryCodeAndName {
		fmt.Printf("Country Code: %s, Country Name: %s\n",
			fullCountryInfoResponse.CountriesUsingCurrencyResult.TCountryCodeAndName[i].SISOCode,
			fullCountryInfoResponse.CountriesUsingCurrencyResult.TCountryCodeAndName[i].SName)
	}

}
