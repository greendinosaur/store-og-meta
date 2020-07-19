package main

import "testing"

func TestFindMetaTags12(t *testing.T) {
	var tests = []struct {
		url    string
		ogData ogMetaData
	}{
		{url: "https://www.facebook.com", ogData: ogMetaData{site: "https://www.facebook.com", ogImage: "https://www.facebook.com/images/fb_icon_325x325.png", ogDescription: "", ogWidth: 0, ogHeight: 0}},
		{url: "https://www.google.com", ogData: ogMetaData{site: "https://www.google.com", ogImage: "", ogDescription: "", ogWidth: 0, ogHeight: 0}},
	}

	for _, test := range tests {
		if ogResponse := findMetaTags(test.url); ogResponse != test.ogData {
			t.Errorf("response value didn't match actual value: %v, %v", ogResponse, test.ogData)
		}
	}

}
