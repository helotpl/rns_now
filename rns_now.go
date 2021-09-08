package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"
	"time"
)

type TimeRange struct {
	start, end time.Time
}

func ParseTimeRange(input, separator string) (TimeRange,error) {
	times := strings.Split(input, separator)
	if len(times) != 2 {
		return TimeRange{},fmt.Errorf("no or too many separators found")
	}
	//strStartTime := times[0]
	//strEndTime := times[1]
	//time.Parse("")

	//TODO: unfinished
	return TimeRange{},nil
}

func main() {
	doc, err := htmlquery.LoadURL("https://nowyswiat.online/")
	if err != nil {
		log.Fatal(err)
	}

	data := htmlquery.FindOne(doc, "//*[starts-with(@id,'proradio-upcomingcarousel-c')]")

	list, err := htmlquery.QueryAll(data, "//*[@class=\"proradio-post__card__cap\"]")
	if err != nil {
		log.Fatal(err)
	}
	for _, x := range list {
		//fmt.Println(htmlquery.OutputHTML(x, true))
		nameNode := htmlquery.FindOne(x, "//h3[@class='proradio-post__title proradio-cutme-t-2 proradio-h4']/a")
		name := strings.TrimSpace(htmlquery.InnerText(nameNode))
		timeNode := htmlquery.FindOne(x, "//p[@class='proradio-itemmetas']")
		timeRange := strings.TrimSpace(htmlquery.InnerText(timeNode))
		fmt.Println(name, timeRange)
	}
}
