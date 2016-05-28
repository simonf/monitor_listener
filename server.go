package monitor_listener

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func StartServer(port int) {
	ph := http.HandlerFunc(handler)
	sh := http.FileServer(http.Dir("."))
	http.Handle("/index.html", ph)
	http.Handle("/*", sh)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		b := GetSiteStatusAsJSON()
		w.Write(b)
	} else {
		sendHTMLResponse(w)
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

func sendHTMLResponse(w http.ResponseWriter) {
	const tpl = `
  <!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8"/>
      <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
      <meta name="viewport" content="width=device-width, initial-scale=1"/>
    //meta(http-equiv="refresh" content="30")
      <title>Monitor</title>
      <link rel='stylesheet', href='/bootstrap-3.3.6-dist/css/bootstrap.min.css'/>
      <link rel='stylesheet', href='/font-awesome-4.5.0/css/font-awesome.min.css'/>
      <link rel='stylesheet', href='/stylesheets/style.css'/>
      <script src='/javascripts/jquery-1.12.1.min.js'/>
      <script src='/bootstrap-3.3.6-dist/js/bootstrap.min.js'/>
    </head>
    <body>
      <div class='content'>
        <div class='container-fluid'>
          <h1>Monitor</h1>
          <div class='row'>
          {{range .}}
          <div class='col-xs12 col-sm-12 col-md-12 col-lg-6'>
            <div class='row'>
              <div class='col-xs-12 col-sm-12 col-md-12' id='agent'>
                <div class='{{ .Status }}'>
                  <h3><a data-target='#{{ .Name }}' data-toggle='collapse'  href='#{{ .Name }}'>{{ .Name }} &nbsp; {{ .Status }}</a></h3>
                  <p>{{ .IP }}<br/>{{ .Updated }}</p>
                  {{if .Status }}
                    <div class='collapse' id='{{ .Name }}'>
                    {{range .Services}}
                      <div class='row'>
                        <div class='col-xs-12 col-sm-12 col-md-12'>
                          <p>
                            <a class='{{if .Status}}btn btn-lg btn-block btn-success{{else}}btn btn-lg btn-block btn-danger{{end}}' title='{{ .Name }}'>
                              <span class='{{ .Status }}' title='dummy'>{{ .Name}}</span>
                            </a>
                          </p>
                        </div>
                      </div>
                    {{end}}
                    </div>
                  {{end}}
                </div>
              </div>
            </div>
          </div>
          {{end}}
         </div>
        </div>
      </div>
    </body>
  </html>
  `
	t := template.Must(template.New("webpage").Parse(tpl))
	err := t.Execute(w, db.ListComputers())
	if err != nil {
		log.Println("executing template:", err)
	}
}
