package gbot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func McpHandle() http.Handler {
	s := server.NewMCPServer(
		"Server Demo",
		"1.0.0",
	)

	// 添加工具
	{
		calculatorTool := mcp.NewTool("calculate",
			mcp.WithDescription("执行基本的算术运算"),
			mcp.WithString("operation",
				mcp.Required(),
				mcp.Description("要执行的算术运算类型"),
				mcp.Enum("add", "subtract", "multiply", "divide"), // 保持英文
			),
			mcp.WithNumber("x",
				mcp.Required(),
				mcp.Description("第一个数字"),
			),
			mcp.WithNumber("y",
				mcp.Required(),
				mcp.Description("第二个数字"),
			),
		)

		s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			fmt.Println(request)
			return mcp.FormatNumberResult(0.0), nil
		})
	}

	return server.NewStreamableHTTPServer(s)
}
