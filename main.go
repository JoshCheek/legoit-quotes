package main

import (
  "flag"
  "fmt"
  "html/template"
  "log"
  "net/http"
  "os"
  "strings"
  "time"
)

type quotation struct {
  Title            string
  GooglePropertyID string
  Author           string
  Quote            string
}

func splitQuoteUrl(url string, q *quotation) bool {
  normalizedUrl := strings.Replace(url, "//", "/", -1)
  urlParts := strings.Split(normalizedUrl, "/")

  if len(urlParts) == 3 && urlParts[0] == "" {
    q.Author, q.Quote = urlParts[1], urlParts[2]
    return true
  }

  return false
}

func legitServer(w http.ResponseWriter, r *http.Request) {
  q := quotation{
    Title:            "Totally for real quotes",
    GooglePropertyID: os.Getenv("GOOGLE_PROPERTY_ID"),
  }

  whichTemplate := "mainPage"

  if splitQuoteUrl(r.URL.Path, &q) {
    q.Title = fmt.Sprintf("A quote by %s", q.Author)
    whichTemplate = "quotePage"
  }

  legitTemplate.ExecuteTemplate(w, whichTemplate, &q)
}

func main() {
  port := flag.Int("port", 8080, "whatever port number legoit-quotes should run on")
  flag.Parse()

  server := &http.Server{
    Addr:         fmt.Sprintf(":%d", *port),
    Handler:      http.HandlerFunc(legitServer),
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
  }
  log.Fatal(server.ListenAndServe())
}

var legitTemplate = template.Must(template.New("boilerplate").Parse(`
{{ define "boilerplateHeader" }}<!DOCTYPE html>
<html>
  <head>
    <title>{{ .Title }}</title>
    <style>
      a {
        color:               white;
        text-decoration:     none;
      }
      body {
        background-color:    #659EC7;
        font-family:         sans-serif;
        color:               white;
        margin:              0px;
      }
      #content {
        font-size:           200%;
        background-color:    #2B547E;
        margin:              10%;
        padding:             5%;
        -moz-border-radius:  25px;
        border-radius:       25px;
        border-width:        5px;
        border-color:        #A0CFEC;
        border-style:        solid;
      }
      #footer {
        background-color:    #2B547E;
        width:               100%;
        padding:             0px;
        padding-left:        20px;
        padding-top:         20px;
        border-style:        solid;
        border-width:        0px;
        border-top-width:    5px;
        border-color:        #A0CFEC;
        margin:              0%;
        height:              50px;
        position:            fixed;
        bottom:              0px;
      }
      .quote {
        padding-bottom:      5%;
      }
      .author {
        margin-left:         60%;
        color:               #659EC7;
        font-style:          italic;
      }
    </style>

    <script type="text/javascript">

      var _gaq = _gaq || [];
      _gaq.push(['_setAccount', '{{ .GooglePropertyID }}']);
      _gaq.push(['_trackPageview']);

      (function() {
        var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
        ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
        var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
      })();

    </script>
  </head>
  <body>
    <div id="content">
{{ end }}

{{ define "boilerplateFooter" }}
    </div>
    <div id="footer">
      <a href="https://github.com/JoshCheek/legit-quotes">Get the source</a>
    </div>
  </body>
</html>
{{ end }}

{{ define "mainPage" }}{{ template "boilerplateHeader" . }}
<p>Totally for real actual quotes.</p>
{{ template "boilerplateFooter" }}
{{ end }}

{{ define "quotePage" }}{{ template "boilerplateHeader" . }}
<p class="quote">"{{ .Quote }}"</p>
<p class="author">{{ .Author }}</p>
{{ template "boilerplateFooter" }}
{{ end }}
`))
