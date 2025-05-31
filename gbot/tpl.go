package gbot

import (
	"text/template"
	"time"
)

const msgTpl = `
{{ .R.ProjName }} {{ .R.Date.Format "2006年01月02日" }} {{ .R.Date | weekday}}
{{.W.DayWeather}} - {{ .W.NightWeather }} {{.W.NightTemp}}~{{.W.DayTemp}}度
>客流总数为：<font color="warning"> {{ .R.Total }} </font>，晚间（20~22点）客流为：<font color="warning"> {{ .R.Night }} </font>
`

func weekday(t time.Time) string {
	w := t.Weekday()
	switch w {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	}
	return ""
}

func GBotTemplate() *template.Template {
	tpl, _ := template.New("msgTpl").
		Funcs(template.FuncMap{"weekday": weekday}).
		Parse(msgTpl)
	return tpl
}
