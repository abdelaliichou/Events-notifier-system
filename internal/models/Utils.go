package models

import "strings"

// Constant static values
const (
	AppName            = "ICHOU_GoApp"
	Version            = "1.0.0"
	M1_GROUPE_1_lANGUE = "13295"
	M1_GROUPE_2_lANGUE = "13345"
	M1_GROUPE_3_lANGUE = "13397"
	M1_GROUPE_1_OPTION = "7224"
	M1_GROUPE_2_OPTION = "7225"
	M1_GROUPE_3_OPTION = "62962"
	M1_GROUPE_OPTION   = "62090"
	M1_TUTORAT_L2      = "56529"
)

// Function to generate calendar URL with multiple resource IDs

func CalendarURL(nbWeeks string, RESOURCE_ID ...string) string {
	// Join multiple resource IDs with ","
	joinedResources := strings.Join(RESOURCE_ID, ",")

	return "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=" + joinedResources +
		"&projectId=2&calType=ical&" + nbWeeks + "=8&displayConfigId=128"
}

// ICalendar format
/*
BEGIN:VCALENDAR                 # début du calendrier
METHOD:PUBLISH                  # attributs du calendrier
PRODID:-//ADE/version 6.0
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VEVENT                    # début d'un event
DTSTAMP:20250111T152020Z        # attributs d'un event
DTSTART:20250220T070000Z
DTEND:20250220T090000Z
SUMMARY:*Option CM Méthodes approchées
LOCATION:IS_A104
DESCRIPTION:\n\nMASTER 1 INFO\nNGUYEN MINH HIEU\n\n(Updated :20/11/2024 1
8:01)                          # ligne multiple pour plus de données commençant par un espace
UID:ADE60323032342d323032352d5543412d33343338392d302d34
CREATED:19700101T000000Z
LAST-MODIFIED:20241120T170100Z
SEQUENCE:2141518841
END:VEVENT                      # fin d'un event
BEGIN:VEVENT                    # début d'un event
DTSTAMP:20250113T153810Z
DTSTART:20250114T070000Z
DTEND:20250114T090000Z
SUMMARY:CM Big Data Infrastructure
LOCATION:IS_E005_Amphi Garcia
DESCRIPTION:\n\nMASTER 1 INFO\nTOUMANI FAROUK\n\n(Updated :20/11/2024 17:
UID:ADE60323032342d323032352d5543412d34313837392d302d30
CREATED:19700101T000000Z
LAST-MODIFIED:20241120T160400Z
SEQUENCE:2141518784
END:VEVENT                      # fin d'un event
END:VCALENDAR                   # fin du calendrier
*/
