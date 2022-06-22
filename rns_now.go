package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"
	"time"
)

var replacements map[rune]rune = map[rune]rune{
	'ą': 'a',
	'Ą': 'A',
	'ę': 'e',
	'Ę': 'E',
	'ó': 'o',
	'Ó': 'O',
	'ł': 'l',
	'Ł': 'L',
	'ś': 's',
	'Ś': 'S',
	'ż': 'z',
	'Ż': 'Z',
	'ź': 'z',
	'Ź': 'Z',
	'ć': 'c',
	'Ć': 'C',
	'ń': 'n',
	'Ń': 'N',
}

func UnPlString(s string) string {
	ret := make([]rune, 0, len(s))
	i := 0
	for _, c := range s {
		if r, ok := replacements[c]; ok {
			ret = append(ret, r)
		} else {
			ret = append(ret, c)
		}
		i++
	}
	return string(ret)
}

func FileNameizeString(s string, removePL bool) string {
	if removePL {
		s = UnPlString(s)
	}
	ret := make([]rune, len(s))
	for i, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			ret[i] = c
		} else {
			ret[i] = '_'
		}
	}
	return string(ret)
}

type MyTime struct {
	hour, minute int
}

func (mt *MyTime) String() string {
	return fmt.Sprintf("%02d:%02d", mt.hour, mt.minute)
}

func (mt *MyTime) Add(mins int) {
	minutes := mt.hour*60 + mt.minute
	minutes += mins
	minutes %= 1440
	if minutes < 0 {
		minutes += 1440
	}
	mt.hour = minutes / 60
	mt.minute = minutes % 60
}

func (a MyTime) Compare(b MyTime) int {
	if a.hour > b.hour {
		return 1
	}
	if a.hour < b.hour {
		return -1
	}
	if a.minute > b.minute {
		return 1
	}
	if a.minute < b.minute {
		return -1
	}
	return 0
}

func MakeMyTimeNow() MyTime {
	t := time.Now()
	return MyTime{t.Hour(), t.Minute()}
}

type TimeRange struct {
	start, end MyTime
}

func (tr *TimeRange) String() string {
	return tr.start.String() + "--" + tr.end.String()
}

func (tr *TimeRange) MoveMinsBack(m int) {
	tr.start.Add(-1 * m)
	tr.end.Add(-1 * m)
}

func (tr *TimeRange) IsCurrent() bool {
	curr := MakeMyTimeNow()

	if tr.start.Compare(tr.end) > 0 {
		// start > end
		return tr.start.Compare(curr) < 1 || tr.end.Compare(curr) > -1 //start >= curr || end <= curr
	} else if tr.start.Compare(tr.end) < 0 {
		// start < end
		//fmt.Println(tr.start.Compare(curr) < 1, tr.end.Compare(curr) > -1)
		return tr.start.Compare(curr) < 1 && tr.end.Compare(curr) > -1 //start <= curr <= end
	} else {
		// start == end
		if curr.Compare(tr.start) == 0 {
			return true
		} else {
			return false
		}
	}
}

func MakeTime(s string) (MyTime, error) {
	var hour, minute int
	_, err := fmt.Sscanf(s, "%2d:%2d", &hour, &minute)
	if err != nil {
		return MyTime{}, err
	}
	return MyTime{hour, minute}, nil
}

func ParseTimeRange(input, separator string) (TimeRange, error) {
	times := strings.Split(input, separator)
	if len(times) != 2 {
		return TimeRange{}, fmt.Errorf("no or too many separators found")
	}
	startTime, err := MakeTime(times[0])
	if err != nil {
		return TimeRange{}, err
	}
	endTime, err := MakeTime(times[1])
	if err != nil {
		return TimeRange{}, err
	}
	return TimeRange{startTime, endTime}, nil
}

func main() {
	doc, err := htmlquery.LoadURL("https://nowyswiat.online/")
	if err != nil {
		log.Fatal(err)
	}

	data := htmlquery.FindOne(doc, "/html/body/main/div[2]/div/div/ul")

	list, err := htmlquery.QueryAll(data, "/li")
	if err != nil {
		log.Fatal(err)
	}
	var selectedNames []string
	for _, x := range list {
		//fmt.Println(htmlquery.OutputHTML(x, true))
		nameNode := htmlquery.FindOne(x, "/div[2]/div[2]/a/h2/span")
		name := strings.TrimSpace(htmlquery.InnerText(nameNode))
		timeNode := htmlquery.FindOne(x, "/div[2]/div[1]/div")
		timeRange := strings.TrimSpace(htmlquery.InnerText(timeNode))
		tr, err := ParseTimeRange(timeRange, " - ")
		if err != nil {
			log.Fatal(err)
		}
		tr.MoveMinsBack(10)
		if tr.IsCurrent() {
			selectedNames = append(selectedNames, FileNameizeString(name, true))
		}
		//fmt.Println(name, FileNameizeString(name, true), timeRange, tr.String(), tr.IsCurrent())
	}
	if len(selectedNames) > 0 {
		fmt.Print(selectedNames[0])
	} else {
		log.Fatal("zero names match")
	}
}
