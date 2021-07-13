package main

import (
	"io"
	"log"
	"strconv"
	"strings"
)

type RawAlert []string

type Alert struct {
	RawAlert

	Name        string
	Author      string
	Description string
	Type        string
	WarningLev  float64
	CriticalLev float64
	Sustain     string
	Action      string
}

type AlertInfo struct {
	Header       []string
	AlertsZUpper []Alert
	AlertsZLower []Alert
	AlertsLUpper []Alert
	AlertsLLower []Alert
	AlertsOther  []Alert
}

// NewOptim Constructor would create an optimised version of passed alert.
func NewOptim(source string) (*AlertInfo, error) {
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

// TODO: assess refactoring of logic in a generic node extractor.
func extractAlert(r RawAlert) Alert {
	var a Alert
	posToken := startPos(r[1])

	raw := []string{}

	for i := 1; i < len(r); i++ {
		pos := strings.Index(r[i], ":")

		if pos != -1 {
			token := strings.Trim(r[i][:pos], " ")
			tokenVal := r[i][pos+1:]

			switch token {
			case "author":
				{
					if startPos(r[i]) == posToken {
						a.Author = tokenVal
						continue
					}
				}

			case "name":
				{
					if startPos(r[i]) == posToken {
						a.Name = tokenVal
						continue
					}
				}

			case "type":
				{
					if startPos(r[i]) == posToken {
						a.Type = strings.Trim(tokenVal, " \n")
						continue
					}
				}

			case "description":
				{
					if startPos(r[i]) == posToken {
						item := []string{}

						for (startPos(r[i+1]) > posToken) || len(r[i+1]) == 1 {
							item = append(item, r[i])
							i++
						}

						item = append(item, r[i])
						a.Description = strings.Join(item, "")

						continue
					}
				}

			case "query":
				{
					if startPos(r[i]) == posToken {
						item := []string{}

						for (startPos(r[i+1]) > posToken) || len(r[i+1]) == 1 {
							item = append(item, r[i])
							i++
						}

						item = append(item, r[i])
						a.Action = strings.Join(item, "")

						continue
					}
				}

			case "warn":
				{
					if startPos(r[i]) == posToken {
						var errWarn error
						a.WarningLev, errWarn = strconv.ParseFloat(strings.Trim(tokenVal, "  \n"), 64)
						if errWarn != nil {
							log.Println(errWarn)
						}

						continue
					}
				}

			case "critical":
				{
					if startPos(r[i]) == posToken {
						var errCri error
						a.CriticalLev, errCri = strconv.ParseFloat(strings.Trim(tokenVal, "  \n"), 64)
						if errCri != nil {
							log.Println(errCri)
						}

						continue
					}
				}

			case "sustainPeriod":
				{
					if startPos(r[i]) == posToken {
						a.Sustain = tokenVal
						continue
					}
				}
			}
		}

		raw = append(raw, r[i])
	}

	tabs := "  "

	alertPrime := []string{
		tabs + "- alert:" + "\n",
		tabs + "    name:" + a.Name,
	}

	a.RawAlert = append(a.RawAlert, alertPrime...)

	if len(a.Author) != 0 {
		a.RawAlert = append(a.RawAlert, tabs+"    author:"+a.Author)
	}

	if len(a.Description) != 0 {
		a.RawAlert = append(a.RawAlert, a.Description)
	}

	if len(a.Type) != 0 {
		a.RawAlert = append(a.RawAlert, tabs+"    type:"+" "+a.Type+"\n")
	}

	a.RawAlert = append(a.RawAlert, tabs+"    warn:"+" "+strconv.FormatFloat(a.WarningLev, 'f', -1, 64)+"\n")
	a.RawAlert = append(a.RawAlert, tabs+"    critical:"+" "+strconv.FormatFloat(a.CriticalLev, 'f', -1, 64)+"\n")

	if len(a.Sustain) != 0 {
		a.RawAlert = append(a.RawAlert, tabs+"    sustainPeriod:"+a.Sustain)
	}

	if len(a.Action) != 0 {
		a.RawAlert = append(a.RawAlert, a.Action)
	}

	a.RawAlert = append(a.RawAlert, raw...)
	a.RawAlert = append(a.RawAlert, "\n")

	return a
}
