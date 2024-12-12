package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

//go:generate go run github.com/gqgs/argsgen@latest

type options struct {
	start  int  `arg:"start year,required"`
	end    int  `arg:"end year,required"`
	days   int  `arg:"number of days to compute,required"`
	top    int  `arg:"number of top times to consider,required"`
	header bool `arg:"generate CSV header"`
}

func main() {
	opts := options{
		start:  2015,
		end:    time.Now().Year(),
		top:    10,
		header: true,
	}
	opts.MustParse()

	if err := generateStats(opts, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func generateStats(opts options, writer io.StringWriter) error {
	if opts.header {
		generateHeader(opts.start, opts.end, writer)
	}
	var statsByYear [][]int
	for year := opts.start; year <= opts.end; year++ {
		stats, err := yearStats(year, opts.days, opts.top)
		if err != nil {
			return err
		}
		statsByYear = append(statsByYear, stats)
	}

	for dayIndex := range opts.days {
		writer.WriteString(fmt.Sprint(dayIndex + 1))
		for yearIndex := range statsByYear {
			writer.WriteString(",")
			writer.WriteString(fmt.Sprint(statsByYear[yearIndex][dayIndex]))
		}
		writer.WriteString("\n")
	}
	return nil
}

func generateHeader(start, end int, writer io.StringWriter) {
	writer.WriteString("day")
	for year := start; year <= end; year++ {
		writer.WriteString(",")
		writer.WriteString(fmt.Sprint(year))
	}
	writer.WriteString("\n")
}

func yearStats(year, days, top int) ([]int, error) {
	var averages []int
	for i := range days {
		times, err := yearDayStats(year, i+1, top)
		if err != nil {
			return nil, err
		}
		average, err := timeAverage(times)
		if err != nil {
			return nil, err
		}
		averages = append(averages, average)
	}
	return averages, nil
}

func timeAverage(times []string) (int, error) {
	var result int
	for _, t := range times {
		parsed, err := time.Parse(time.TimeOnly, t)
		if err != nil {
			return 0, err
		}
		seconds := parsed.Second() + parsed.Minute()*60
		result += seconds
	}

	return result / len(times), nil
}

func yearDayStats(year, day, top int) ([]string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/leaderboard/day/%d", year, day)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var times []string
	for desc := range node.Descendants() {
		if !isTimeSpan(desc) {
			continue
		}
		var month, day, time string
		fmt.Sscanf(desc.FirstChild.Data, "%s %s  %s", &month, &day, &time)
		times = append(times, time)
		if len(times) >= top {
			break
		}
	}

	return times, nil
}

func isTimeSpan(node *html.Node) bool {
	if node.Data != "span" {
		return false
	}

	for _, attr := range node.Attr {
		if attr.Key == "class" && attr.Val == "leaderboard-time" {
			return true
		}
	}
	return false
}
