package tcmb

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	errUnexpectedStatus = errors.New("tcmb: unexpected status")
)

const url = "https://www.tcmb.gov.tr/kurlar/today.xml"

//Currency ...
type Currency struct {
	CrossOrder      string `xml:"CrossOrder,attr"`
	Kod             string `xml:"Kod,attr"`
	CurrencyCode    string `xml:"CurrencyCode,attr"`
	Unit            string `xml:"Unit"`
	Isim            string `xml:"Isim"`
	CurrencyName    string `xml:"CurrencyName"`
	ForexBuying     string `xml:"ForexBuying"`
	ForexSelling    string `xml:"ForexSelling"`
	BanknoteBuying  string `xml:"BanknoteBuying"`
	BanknoteSelling string `xml:"BanknoteSelling"`
	CrossRateUSD    string `xml:"CrossRateUSD"`
	CrossRateOther  string `xml:"CrossRateOther"`
}

//Response is the model of the response returned from TCMB API
type Response struct {
	XMLName    xml.Name   `xml:"Tarih_Date"`
	Tarih      string     `xml:"Tarih,attr"`
	Date       string     `xml:"Date,attr"`
	BultenNo   string     `xml:"Bulten_No,attr"`
	Currencies []Currency `xml:"Currency"`
}

//Get fetches latest exhange rates from TCMB
func Get() (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errUnexpectedStatus
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := new(Response)

	err = xml.Unmarshal(bodyBytes, r)

	return r, err
}

// Schedule executes the given function with the given interval
//
// When the signal received from quit chan, the gorouitine quits.
func Schedule(duration time.Duration, quit chan bool, handler func(*Response, error)) {
	ticker := time.NewTicker(duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				handler(Get())
			case <-quit:
				return
			}
		}
	}()
}
