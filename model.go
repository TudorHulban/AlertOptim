package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type RawAlert []string

type Alert struct {
	RawAlert

	Name        string
	Type        string
	WarningLev  int
	CriticalLev int
	Sustain     int
}

type AlertInfo struct {
	Header       []string
	AlertsZUpper []Alert
	AlertsZLower []Alert
	AlertsLUpper []Alert
	AlertsLLower []Alert
	AlertsOther  []Alert
}

func NewSimple(source string) (*AlertInfo, error) {
	data, errRead := readFile(source)
	if errRead != nil {
		return nil, errRead
	}

	header, posAlerts := isolateHeader(data)
	alertData := isolateAlertData(data, posAlerts)

	res := mapAlerts(extractAlerts(alertData))
	res.Header = header

	return &res, nil
}

func (a *AlertInfo) Spool(w io.Writer) {
	s := func(data []Alert, w io.Writer) {
		for _, a := range data {
			w.Write([]byte(strings.Join(a.RawAlert, "")))
		}
	}

	w.Write([]byte(strings.Join(a.Header, "")))
	s(a.AlertsZUpper, w)
	s(a.AlertsZLower, w)
	s(a.AlertsLUpper, w)
	s(a.AlertsLLower, w)
	s(a.AlertsOther, w)
}

func mapAlerts(alerts []RawAlert) AlertInfo {
	var res AlertInfo

	for _, a := range alerts {
		alert := extractAlert(a)

		if alert.Type == "ZONE" && (alert.CriticalLev > alert.WarningLev) {
			res.AlertsZUpper = append(res.AlertsZUpper, alert)
			continue
		}

		if alert.Type == "ZONE" && (alert.CriticalLev <= alert.WarningLev) {
			res.AlertsZLower = append(res.AlertsZLower, alert)
			continue
		}

		if alert.Type == "LEGACY" && (alert.CriticalLev > alert.WarningLev) {
			res.AlertsLUpper = append(res.AlertsLUpper, alert)
			continue
		}

		if alert.Type == "LEGACY" && (alert.CriticalLev <= alert.WarningLev) {
			res.AlertsLLower = append(res.AlertsLLower, alert)
			continue
		}

		res.AlertsOther = append(res.AlertsOther, alert)
	}

	return res
}

func readFile(path string) ([]string, error) {
	f, errOpen := os.Open(path)
	if errOpen != nil {
		return nil, errOpen
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var res []string

	for scanner.Scan() {
		res = append(res, scanner.Text()+"\n")
	}
	if errScan := scanner.Err(); errScan != nil {
		return nil, errScan
	}

	return res, nil
}

func isolateHeader(data []string) ([]string, int) {
	var res []string
	var posAlertInfo int

	for i, line := range data {
		res = append(res, line)

		if strings.Contains(line, "alerts:") {
			posAlertInfo = i
			break
		}
	}

	return res, posAlertInfo
}

func isolateAlertData(data []string, posAlert int) []string {
	var res []string

	for i := posAlert + 1; i < len(data); i++ {
		res = append(res, data[i])
	}

	return res
}

func extractAlerts(data []string) []RawAlert {
	if len(data) == 0 {
		return nil
	}

	var res []RawAlert

	i := 0
	alert := []string{}

	for i < len(data) {
		if strings.Contains(data[i], "- alert:") && len(alert) > 0 {
			res = append(res, alert)
			alert = []string{}
		}

		alert = append(alert, data[i])

		if i == len(data)-1 {
			res = append(res, alert)

			break
		}

		i++
	}

	return res
}

func extractAlert(r RawAlert) Alert {
	var a Alert
	raw := []string{}

	for _, line := range r {
		vals := strings.Split(line, ":")

		switch strings.Trim(vals[0], " ") {
		case "name":
			{
				a.Name = strings.Title(vals[1])
				continue
			}

		case "type":
			{
				a.Type = strings.Trim(vals[1], " \n")
				continue
			}

		case "warn":
			{
				var errWarn error
				a.WarningLev, errWarn = strconv.Atoi(strings.Trim(vals[1], "  \n"))
				if errWarn != nil {
					log.Println(errWarn)
				}
				continue
			}

		case "critical":
			{
				var errCri error
				a.CriticalLev, errCri = strconv.Atoi(strings.Trim(vals[1], "  \n"))
				if errCri != nil {
					log.Println(errCri)
				}
				continue
			}

		case "sustainPeriod":
			{
				var errSus error
				a.Sustain, errSus = strconv.Atoi(strings.Trim(vals[1], "  \n"))
				if errSus != nil {
					log.Println(errSus)
				}
				continue
			}

		case "- alert":
			{
				continue
			}
		}

		raw = append(raw, line)
	}

	tabs := "    "

	alertPrime := []string{
		tabs + "- alert:" + "\n",
		tabs + "    name:" + a.Name,
		tabs + "    type:" + " " + a.Type + "\n",
		tabs + "    warn:" + " " + strconv.Itoa(a.WarningLev) + "\n",
		tabs + "    critical:" + " " + strconv.Itoa(a.CriticalLev) + "\n",
		tabs + "    sustainPeriod:" + " " + strconv.Itoa(a.Sustain) + "\n",
	}

	a.RawAlert = append(a.RawAlert, alertPrime...)
	a.RawAlert = append(a.RawAlert, raw...)
	a.RawAlert = append(a.RawAlert, "\n")

	return a
}
