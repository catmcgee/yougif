package model

import "encoding/xml"

type Transcript struct {
	XMLName xml.Name `xml:"transcript"`
	Text    []Text   `xml:"text"`
}

type Text struct {
	Text     string  `xml:",chardata"`
	Start    float64 `xml:"start,attr"`
	Duration float64 `xml:"dur,attr"`
}
