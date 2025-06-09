package puppy

import (
	"context"

	"github.com/imroc/req/v3"
)

type AIReq struct {
	Message string `json:"message"`
}

type AIResp struct {
	Code   string `json:"code"`
	Result string `json:"result"`
}

type AI struct {
	client *req.Client
	base   string
}

func NewAI(base string) *AI {
	c := req.C().SetBaseURL(base)
	return &AI{
		client: c,
		base:   base,
	}

}

func (a *AI) Ask(ctx context.Context, q string) (string, error) {
	var resp AIResp
	_, err := a.client.R().
		SetContext(ctx).
		SetBody(AIReq{Message: q}).
		SetSuccessResult(&resp).
		SetErrorResult(&resp).
		Post("/message")
	if err != nil {
		return "", err
	}

	return resp.Result, nil
}
