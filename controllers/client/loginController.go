package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/tools"
	"forum/tools/session"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

type LoginPage struct {
	Path string
	Err  string
}

func (p *LoginPage) ServeHTTP(w http.ResponseWriter, r *http.Request, m map[string]string) {
	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.Err = err.Error()
			fmt.Print(err)
		} else {
			urlapi := os.Getenv("url_api") + "login"
			params := tools.PlaintTextToMap(reqBody)
			params["username"], _ = url.QueryUnescape(params["username"])
			params["password"], _ = url.QueryUnescape(params["password"])
			jsonData, err := json.Marshal(params)
			if err != nil {
				p.Err = "{\"err\":\"500\",\"msg\":\"" + err.Error() + "\"}"
				tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
				tmpl.Execute(w, p)
				return
			}
			res, err := http.Post(urlapi, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				p.Err = "{\"err\":\"500\",\"msg\":\"" + err.Error() + "\"}"
				tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
				tmpl.Execute(w, p)
				return
			}
			resp, err := ioutil.ReadAll(res.Body)
			if err != nil {
				p.Err = "{\"err\":\"500\",\"msg\":\"" + err.Error() + "\"}"
				tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
				tmpl.Execute(w, p)
				return
			}
			if res.StatusCode != 200 {
				p.Err = "{\"err\":\"500\",\"msg\":\"" + string(resp) + "\"}"
				tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
				tmpl.Execute(w, p)
				return
			} else {
				data := map[string]string{}
				err := json.Unmarshal(resp, &data)
				if err != nil {
					p.Err = "{\"err\":\"500\",\"msg\":\"" + string(resp) + "\"}"
					tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
					tmpl.Execute(w, p)
					return
				}
				cookie := http.Cookie{Name: session.GlobalSessions.CookieName, Value: data["SID"], Path: "/", HttpOnly: true, MaxAge: int(session.GlobalSessions.Maxlifetime)}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			}
		}
	}
	tmpl := template.Must(template.ParseFiles(CurrentFolder + p.Path))
	tmpl.Execute(w, p)
}