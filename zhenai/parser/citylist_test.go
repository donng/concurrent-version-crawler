package parser

import (
	"testing"
	"io/ioutil"
)

func TestParseCityList(t *testing.T) {
	bytes, e := ioutil.ReadFile("citylist_test_data.html")
	if e != nil {
		panic(e)
	}


	result := ParseCityList(bytes)

	const cityNum = 470
	if len(result.Requests) != cityNum {
		t.Errorf("result should have %d requests, but have %d", cityNum, len(result.Requests))
	}
}
