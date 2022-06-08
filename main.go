/*
	Copyright (C) 2022  Michael Ablassmeier <abi@grinser.de>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License along
    with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/buger/goterm"
	"github.com/jessevdk/go-flags"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

type Flights struct {
	Data    []FlightInfo `json:"data"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

type FlightInfo struct {
	LastName        string `json:"lastname"`
	FirstName       string `json:"firstname"`
	BestTaskPoints  string `json:"besttaskpoints"`
	TakeoffLocation string `json:"takeofflocation"`
	LandingLocation string `json:"landinglocation"`
}

type Options struct {
	Day      string `short:"d" long:"day" description:"date selection: 08.06.2022" required:"false"`
	Interval int    `short:"i" long:"interval" description:"Refresh interval in seconds" default:"0"`
	Limit    int    `short:"l" long:"limit" description:"Limit to X results" default:"0"`
}

func json_loads(data []byte) Flights {
	var resp Flights
	err := json.Unmarshal([]byte(data), &resp)
	if err != nil {
		logrus.Error("Cant load json response:")
		logrus.Fatal(err)
	}
	return resp
}

func httpReq(url string) Flights {
	logrus.Debug(url)

	response, error := http.Get(url)
	if error != nil {
		logrus.Fatal(error)
	}
	body, _ := ioutil.ReadAll(response.Body)
	logrus.Debugf("Response: [%s]", string(body))
	return json_loads(body)
}

func success(resp Flights) bool {
	if resp.Success {
		return true
	}
	return false
}

func clearconsole(options Options) {
	if options.Interval > 0 {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Flush()
	}
}

func main() {
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	currentTime := time.Now()
	var day string
	if options.Day == "" {
		day = currentTime.Format("02.01.2006")
	} else {
		day = options.Day
	}
	Api := struct {
		url string
	}{
		url: "https://en.dhv-xc.de/api/fli/flights?d=",
	}

	Api.url = Api.url + day

	for {
		clearconsole(options)

		f := httpReq(Api.url)
		if !success(f) {
			logrus.Fatalf("Request failed: [%s]", f.Message)
		}

		sort.SliceStable(f.Data, func(i, j int) bool {
			floatNumA, _ := strconv.ParseFloat(f.Data[i].BestTaskPoints, 32)
			floatNumB, _ := strconv.ParseFloat(f.Data[j].BestTaskPoints, 32)
			return floatNumA > floatNumB
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"Surname",
			"XC-Points",
			"Takeoff",
			"Landing",
		})
		table.SetBorder(false)
		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
			tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor},
		)
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		)

		for i := 0; i < len(f.Data); i++ {
			fp, _ := strconv.ParseFloat(f.Data[i].BestTaskPoints, 32)
			points := fmt.Sprintf("%.2f", fp)
			table.Append([]string{
				f.Data[i].FirstName,
				f.Data[i].LastName,
				points,
				f.Data[i].TakeoffLocation,
				f.Data[i].LandingLocation,
			})

			if options.Limit > 0 && i+1 >= options.Limit {
				break
			}
		}
		table.Render()
		sleep := time.Duration(1.15*float64(options.Interval)) * time.Second
		if options.Interval == 0 {
			break
		}
		time.Sleep(sleep)
	}
}
