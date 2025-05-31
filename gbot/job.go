package gbot

import (
	"context"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/twiglab/puppy"
	"github.com/xxl-job/xxl-job-executor-go"

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

func (b *GBotJob) Run(ctx context.Context, req *xxl.RunReq) error {
	var (
		err error
		jp  JobParam
	)

	if err = json.Unmarshal([]byte(req.ExecutorParams), &jp); err != nil {
		return err
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())

	br := BotResult{
		ProjName: jp.Proj,
		Date:     now,
	}

	br.Total, err = b.Dcp.Sum(ctx, start, end, jp.Entry)
	if err != nil {
		return err
	}

	start = time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location())
	br.Night, err = b.Dcp.Sum(ctx, start, end, jp.Entry)
	if err != nil {
		return err
	}

	wi, err := b.Weater.GetWeather(ctx, b.AdCode)
	if err != nil {
		return err
	}

	root := map[string]any{
		"W": &wi,
		"R": &br,
	}

	var sb strings.Builder
	sb.Grow(256)
	if err = b.Tpl.Execute(&sb, root); err != nil {
		return err
	}

	return b.MsgBot.SendMarkdownMessage(sb.String())
}
