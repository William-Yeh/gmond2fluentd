// Redirect metrics from Ganglia Monitoring Daemon (gmond) to Fluentd.
//
//
// Copyright 2015 William Yeh <william.pjyeh@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"net"

	"bytes"
	"fmt"
	"strconv"
	"time"

	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
)

const PERIOD_DEFAULT int = 60 // seconds
const PERIOD_MIN int = 30     // seconds

const USAGE string = `Extract metrics from Ganglia Monitoring Daemon (gmond) to Fluentd.

Usage:
  gmond2fluentd file <json_file>  [options]  
  gmond2fluentd tcp  [options]  
  gmond2fluentd --help
  gmond2fluentd --version

Options:
  -s <gmond>, --src <gmond>         gmond source host:port
                                      [default: 127.0.0.1:8649].
  -d <fluentd>, --dest <fluentd>    fluentd in_forward TCP host:port
                                      [default: 127.0.0.1:24224].
  -t <tag>, --tag <tag>             tag sending to Fluentd's in_forward plugin
                                      [default: ganglia].
  -p <seconds>, --period <seconds>  interval of metric query [default: 60]
  --stdout                          also dump to stdout.`

// syntax sugar
type MetricEntry map[string]string

// Ganglia's XML format
var REGEX_METRIC_TIME = regexp.MustCompile(`^\s*<HOST\s+NAME=.+REPORTED="?(\d+)"?`)
var REGEX_METRIC_LINE = regexp.MustCompile(`^\s*<METRIC\s+(.+)>\s*$`)
var REGEX_METRIC_ITEM_PAIR = regexp.MustCompile(`([^=]+)="?([^"]*)"?\s*`)

func main() {
	arguments := process_cmdline()
	//fmt.Println("---", arguments)

	period, _ := strconv.Atoi(arguments["--period"].(string))

	ticker := time.NewTicker(time.Duration(period) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				main_procedure(arguments)
			case <-quit:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				os.Exit(0)
				return
			}
		}

	}()

	// main thread: sleep forever
	select {}
}

// This func parses and validates cmdline args
func process_cmdline() map[string]interface{} {

	arguments, _ := docopt.Parse(USAGE, nil, true, "0.1", false)

	// validate "--period"
	period, err := strconv.Atoi(arguments["--period"].(string))
	if err != nil {
		arguments["--period"] = strconv.Itoa(PERIOD_DEFAULT)
	} else if period < PERIOD_MIN {
		arguments["--period"] = strconv.Itoa(PERIOD_MIN)
	}

	//fmt.Println("---", arguments)
	return arguments
}

func main_procedure(arguments map[string]interface{}) {
	// connect to this socket
	ganglia_endpoint := arguments["--src"].(string)
	conn, err := net.Dial("tcp", ganglia_endpoint)
	if err != nil {
		checkError(err)
		return
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	metric_entries := getGangliaMetrics(conn)
	output(metric_entries, arguments)
}

// query Ganglia monitor daemon (gmond) to obtain metrics
func getGangliaMetrics(conn net.Conn) []MetricEntry {

	var entries = make([]MetricEntry, 50)
	var timestamp string = strconv.FormatInt(time.Now().Unix(), 10)
	//fmt.Println(timestamp)

	scanner := bufio.NewScanner(conn) // receive XML from gmond
	for scanner.Scan() {              // for each line
		line := scanner.Text()
		//fmt.Println(line)

		if result := REGEX_METRIC_TIME.FindStringSubmatch(line); result != nil {
			timestamp = result[1]
			//fmt.Println(timestamp)
		} else if result := REGEX_METRIC_LINE.FindStringSubmatch(line); result != nil {
			entry, err := parseMetricLine(result[1])
			//fmt.Println(entry)

			if err == nil {
				entry["time"] = timestamp
				entries = append(entries, entry)
				//fmt.Println(entries)
			}
		}
	}
	//fmt.Println(entries)

	// compact nil slice items
	metric_entries := entries[:0]
	for _, x := range entries {
		if x != nil {
			metric_entries = append(metric_entries, x)
		}
	}

	//fmt.Println(metric_entries)
	return metric_entries
}

func parseMetricLine(line string) (MetricEntry, error) {

	//fmt.Printf("=== %s \n", line)

	res := REGEX_METRIC_ITEM_PAIR.FindAllStringSubmatch(line, -1)
	if res == nil {
		return nil, errors.New("mismatch metric in xml")
	}
	//fmt.Printf("%v \n", res)

	entry := map[string]string{}

	for _, v := range res {
		//fmt.Println(v)
		key := strings.ToLower(v[1])
		value := v[2]
		entry[key] = value
	}

	//fmt.Println(entry)
	return entry, nil
}

func output(metric_entries []MetricEntry, arguments map[string]interface{}) {

	// keep ref count to avoid gc
	var out_file *os.File
	var fluentd_conn net.Conn

	will_send_to_tcp := arguments["tcp"].(bool)
	will_append_to_file := arguments["file"].(bool)
	will_send_to_stdout := arguments["--stdout"].(bool)

	if will_send_to_tcp {
		fluentd_endpoint := arguments["--dest"].(string)
		conn, err := net.Dial("tcp", fluentd_endpoint)
		if err != nil {
			checkError(err)
			return
		}
		fluentd_conn = conn // incr ref count
		defer func() {
			if fluentd_conn != nil {
				fluentd_conn.Close()
			}
		}()
	}

	if will_append_to_file {
		f, err := os.OpenFile(arguments["<json_file>"].(string), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		out_file = f // incr ref count
		defer func() {
			if out_file != nil {
				out_file.Close()
			}
		}()
	}

	fluentd_tag := arguments["--tag"].(string)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if value, ok := metric_entries[0]["time"]; ok {
		timestamp = value
	}

	for _, item := range metric_entries {
		s, err := json.Marshal(item)
		if err != nil {
			continue
		}

		json_string := string(s)

		if will_send_to_stdout {
			fmt.Println(json_string)
		}
		if will_append_to_file {
			out_file.WriteString(json_string)
			out_file.WriteString("\n")
		}
		if will_send_to_tcp {
			// format: [tag, time, record]
			// @see http://docs.fluentd.org/articles/in_forward
			var buffer bytes.Buffer
			buffer.WriteString("[\"")
			buffer.WriteString(fluentd_tag)
			buffer.WriteString("\",")
			buffer.WriteString(timestamp)
			buffer.WriteString(",")
			buffer.WriteString(json_string)
			buffer.WriteString("]\n")

			fluentd_conn.Write(buffer.Bytes())
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//os.Exit(1)
	}
}
