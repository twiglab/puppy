package puppy

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/twiglab/doggy/pkg/oc"
	"github.com/ybbus/jsonrpc/v3"
)

type DcpServ struct {
	URL string
	cli jsonrpc.RPCClient
}

func NewDcpServ(url string, client *req.Client) *DcpServ {
	fmt.Println(url)
	jr := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{HTTPClient: client})
	return &DcpServ{
		URL: url,
		cli: jr,
	}
}

func (s *DcpServ) Sum(ctx context.Context, start, end time.Time, ids []string) (int64, error) {
	var out oc.Reply
	err := s.cli.CallFor(ctx, &out, oc.CALL_SUM, &oc.AreaArg{
		Start: start.UnixMilli(),
		End:   end.UnixMilli(),
		IDs:   ids,
	})
	if err != nil {
		return 0, err
	}
	return out.ValueA, nil
}
