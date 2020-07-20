package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMetaTags(t *testing.T) {
	var tests = []struct {
		url    string
		ogData ogMetaData
	}{
		{url: "https://www.facebook.com", ogData: ogMetaData{site: "https://www.facebook.com", ogImage: "https://www.facebook.com/images/fb_icon_325x325.png", ogDescription: "", ogWidth: 0, ogHeight: 0}},
	}

	for _, test := range tests {
		ogResponse := findMetaTags(test.url)
		assert.Equal(t, ogResponse.site, test.ogData.site)
		assert.Equal(t, ogResponse.ogDescription, test.ogData.ogDescription)
		assert.Equal(t, ogResponse.ogImage, test.ogData.ogImage)
		assert.Equal(t, ogResponse.ogHeight, test.ogData.ogHeight)
		assert.Equal(t, ogResponse.ogWidth, test.ogData.ogWidth)
	}

}
