package main

import (
	"strconv"
)

const (
	CONFIG_PATH = "http://localhost:8080/resources/"
)

func UCA_URL(nbWeeks string, resources []Resource) string {
	// Join multiple resource IDs with ","
	IDs := ""
	for _, resource := range resources {
		if IDs != "" {
			IDs += ","
		}
		IDs += strconv.Itoa(resource.UcaID)
	}

	return "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=" + IDs +
		"&projectId=2&calType=ical&" + nbWeeks + "=8&displayConfigId=128"
}
