package puppy

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xxl-job/xxl-job-executor-go"
)

const key_xxl_req = "_xxl_param"

type XxlJob interface {
	Name() string
	Run(context.Context, *xxl.RunReq) error
}

type LocalExec struct {
	XxlExec xxl.Executor

	mux http.Handler
}

func (e *LocalExec) Init() *LocalExec {
	e.XxlExec.Init()
	e.mux = XxlJobMux(e.XxlExec)
	return e
}

func (e *LocalExec) RegJob(job XxlJob) *LocalExec {
	name, jobF := XxlJobFunc(job)
	e.XxlExec.RegTask(name, jobF)
	return e
}

func (e *LocalExec) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mux.ServeHTTP(w, r)
}

func XxlJobMux(exec xxl.Executor) http.Handler {
	mux := chi.NewMux()

	mux.Post("/run", exec.RunTask)
	mux.Post("/kill", exec.KillTask)
	mux.Post("/log", exec.TaskLog)
	mux.Post("/beat", exec.Beat)
	mux.Post("/idleBeat", exec.IdleBeat)

	return mux
}

func WithRunReq(ctx context.Context, req *xxl.RunReq) context.Context {
	return context.WithValue(ctx, key_xxl_req, req)
}

func RunReq(ctx context.Context) *xxl.RunReq {
	return ctx.Value(key_xxl_req).(*xxl.RunReq)
}

func XxlJobFunc(job XxlJob) (string, xxl.TaskFunc) {
	return job.Name(), func(ctx context.Context, req *xxl.RunReq) string {
		if err := job.Run(ctx, req); err != nil {
			return "error: " + err.Error()
		}
		return "ok"
	}
}
