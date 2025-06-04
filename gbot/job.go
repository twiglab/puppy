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

	"github.com/xen0n/go-workwx/v2"
)

type GBotJob struct {
	JobName string
	AdCode  string

	Dcp    *puppy.DcpServ
	MsgBot *workwx.WebhookClient
	Weater *puppy.AmapWeather

	Tpl *template.Template
}

func (b *GBotJob) Name() string {
	return b.JobName
}

func (b *GBotJob) Run(ctx context.Context, req *xxl.RunReq) (fmt.Stringer, error) {
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

	wi, _ := b.Weater.GetWeather(ctx, b.AdCode)

	root := map[string]any{
		"W": &wi,
		"R": &br,
	}

	var sb strings.Builder
	sb.Grow(256)
	if err = b.Tpl.Execute(&sb, root); err != nil {
		return nil, err
	}

	err = b.MsgBot.SendMarkdownMessage(sb.String())

	return xxl.JobRtn(err)
}
