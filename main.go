package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

//go:generate go run github.com/gqgs/argsgen@latest

type options struct {
	startYear int  `arg:"start year,required"`
	endYear   int  `arg:"end year,required"`
	startDay  int  `arg:"start day,required"`
	endDay    int  `arg:"end day,required"`
	top       int  `arg:"number of top times to consider,required"`
	header    bool `arg:"generate CSV header"`
	latest    bool `arg:"compute stats of latest day"`
}

func main() {
	opts := options{
		startYear: 2015,
		endYear:   time.Now().Year(),
		startDay:  1,
		endDay:    time.Now().Day(),
		top:       10,
		header:    true,
	}
	opts.MustParse()

	if opts.latest {
		opts.header = false
		opts.startDay = opts.endDay
	}

	if err := generateStats(opts, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func generateStats(opts options, writer io.StringWriter) error {
	if opts.header {
		generateHeader(opts.startYear, opts.endYear, writer)
	}
	var statsByYear [][]int
	for year := opts.startYear; year <= opts.endYear; year++ {
		stats, err := yearStats(year, opts.startDay, opts.endDay, opts.top)
		if err != nil {
			return err
		}
		statsByYear = append(statsByYear, stats)
	}

	for dayIndex := range opts.endDay - opts.startDay + 1 {
		writer.WriteString(fmt.Sprint(opts.startDay + dayIndex))
		for yearIndex := range statsByYear {
			writer.WriteString(",")
			writer.WriteString(fmt.Sprint(statsByYear[yearIndex][dayIndex]))
		}
		writer.WriteString("\n")
	}
	return nil
}

func generateHeader(startYear, endYear int, writer io.StringWriter) {
	writer.WriteString("day")
	for year := startYear; year <= endYear; year++ {
		writer.WriteString(",")
		writer.WriteString(fmt.Sprint(year))
	}
	writer.WriteString("\n")
}

func yearStats(year, startDay, endDay, top int) ([]int, error) {
	averages := make([]int, endDay-startDay+1)
	g := new(errgroup.Group)
	for day := startDay; day <= endDay; day++ {
		g.Go(func() error {
			times, err := yearDayStats(year, day, top)
			if err != nil {
				return err
			}
			average, err := timeAverage(times)
			if err != nil {
				return err
			}
			averages[day-startDay] = average
			return nil
		})
	}
	return averages, g.Wait()
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
