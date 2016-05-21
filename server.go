package monitor_listener

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"simonf.net/monitor_db"
	"strconv"
)

func StartServer(port int) {
	ph := http.HandlerFunc(handler)
	http.Handle("/", ph)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		b := GetSiteStatusAsJSON()
		w.Write(b)
	} else {
		s := GetSiteStatusAsHTML()
		sendHTMLResponse(w, s)
	}
}

func GetSiteStatusAsJSON() []byte {
	var buffer bytes.Buffer
	is_subsequent := false
	buffer.WriteString("{[")
	for _, c := range db.ListComputers() {
		if is_subsequent {
			buffer.WriteString(",")
		} else {
			is_subsequent = true
		}
		buffer.WriteString(string(c.JSON()))
	}
	buffer.WriteString("]}")
	return buffer.Bytes()
}

func GetSiteStatusAsHTML() string {
	var buffer bytes.Buffer
	for _, c := range db.ListComputers() {
		buffer.WriteString(computerAsDiv(c))
	}
	return buffer.String()
}

func sendHTMLResponse(w http.ResponseWriter, content string) {
	var buffer bytes.Buffer
	buffer.WriteString("<html><head><title>Monitor</title></head><body><h1>Site status</h1>\n")
	buffer.WriteString("<div class=\"computer_list\"")
	buffer.WriteString(content)
	buffer.WriteString("</div>\n</body>\n</html>")
	w.Write(buffer.Bytes())
}

func computerAsDiv(c *monitor_db.Computer) string {
	var buffer bytes.Buffer
	buffer.WriteString("<div class=\"computer\">\n")
	buffer.WriteString("\t<div class=\"name\">")
	buffer.WriteString(c.Name)
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("\t<div class=\"status\">")
	buffer.WriteString(c.Status)
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("\t<div class=\"updated\">")
	buffer.WriteString(c.Updated.String())
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("\t<div class=\"services\">")
	for _, svc := range c.Services {
		buffer.WriteString(serviceAsDiv(svc))
	}
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("</div>\n")
	return buffer.String()
}

func serviceAsDiv(s *monitor_db.Service) string {
	var buffer bytes.Buffer
	buffer.WriteString("<div class=\"service\">\n")
	buffer.WriteString("\t<div class=\"name\">")
	buffer.WriteString(s.Name)
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("\t<div class=\"status\">")
	buffer.WriteString(s.Status)
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("\t<div class=\"updated\">")
	buffer.WriteString(s.Updated.String())
	buffer.WriteString("\t</div>\n")
	buffer.WriteString("</div>\n")
	return buffer.String()
}

func makeStatus(url string, status string) string {
	return fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", url, status)
}
