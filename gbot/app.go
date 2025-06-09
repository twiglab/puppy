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

type GBotApp struct {
	App *workwx.WorkwxApp

	Dcp    *puppy.DcpServ
	Weater *puppy.AmapWeather
	AI     *puppy.AI

	Tpl *template.Template
}

func (b *GBotApp) Name() string {
	return "GBotApp"
}

func (b *GBotApp) Run(ctx context.Context, req *xxl.RunReq) (fmt.Stringer, error) {
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
	sb.Grow(512)
	if err = b.Tpl.Execute(&sb, root); err != nil {
		return nil, err
	}

	err = b.App.SendMarkdownMessage(&workwx.Recipient{TagIDs: []string{"1"}}, sb.String(), false)

	return xxl.JobRtn(err)
}

func (a *GBotApp) OnIncomingMessage(msg *workwx.RxMessage) error {
	if msg.MsgType == workwx.MessageTypeText {
		text, _ := msg.Text()
		callWord := text.GetContent()
		s, err := a.AI.Ask(context.Background(), callWord)
		if err != nil {
			return err
		}

		return a.App.SendTextMessage(
			&workwx.Recipient{UserIDs: []string{msg.FromUserID}},
			s,
			false,
		)
	}
	return nil
}
