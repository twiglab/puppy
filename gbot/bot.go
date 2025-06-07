package gbot

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/it512/xxl-job-exec"
	"github.com/twiglab/puppy"
)

type GBot struct {
	Dcp    *puppy.DcpServ
	Weater *puppy.AmapWeather

	Tpl *template.Template
}

func (b *GBot) Name() string {
	return "gbot"
}

func (b *GBot) Run(ctx context.Context, req *xxl.RunReq) (fmt.Stringer, error) {
	var (
		err error
		jp  JobParam

		start time.Time
		end   time.Time
	)

	if err = json.Unmarshal([]byte(req.ExecutorParams), &jp); err != nil {
		return nil, err
	}

	now := time.Now()
	br := BotResult{
		ProjName: jp.Proj,
		Date:     now,
	}

	start, end = OpeningTime(now)
	br.Total, err = b.Dcp.Sum(ctx, start, end, jp.Entry)
	if err != nil {
		return nil, err
	}

	start, end = NightTime(now)
	br.Night, err = b.Dcp.Sum(ctx, start, end, jp.Entry)
	if err != nil {
		return nil, err
	}

	start, end = OpeningTime(BeforWeekDay(now))
	br.BeforWeekDay, err = b.Dcp.Sum(ctx, start, end, jp.Entry)

	wi, _ := b.Weater.GetWeather(ctx, "320100")

	root := map[string]any{
		"W": &wi,
		"R": &br,
	}

	var sb strings.Builder
	sb.Grow(256)
	if err = b.Tpl.Execute(&sb, root); err != nil {
		return nil, err
	}

	return xxl.JobRtn(err)
}
