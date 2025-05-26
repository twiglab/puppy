package puppy

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xxl-job/xxl-job-executor-go"
)

func x() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://127.0.0.1/xxl-job-admin"),
		xxl.AccessToken(""),            //请求令牌(默认为空)
		xxl.ExecutorIp("127.0.0.1"),    //可自动获取
		xxl.ExecutorPort("9999"),       //默认9999（非必填）
		xxl.RegistryKey("golang-jobs"), //执行器名称
		//xxl.SetLogger(&logger{}),       //自定义日志
	)

	exec.Init()
}

const key_xxl_req = "_xxl_param"

type XxlJob interface {
	Run(context.Context, *xxl.RunReq) string
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

func XxlJobFunc(job XxlJob) xxl.TaskFunc {
	return func(ctx context.Context, req *xxl.RunReq) string {
		return job.Run(ctx, req)
	}
}
