package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	_ "github.com/wasilibs/nottinygc"
)

const plugin_module = "promql_extract"

var _prometheus_query_urls = []string{
	"/api/v1/query",
	"/api/v1/query_rnage",
	"/api/v1/series",
}

func main() {
	sdk.Warn(fmt.Sprintf("%s loaded", plugin_module))
	sdk.SetParser(&PromQLParser{})
}

type PromQLParser struct {
}

var _ sdk.Parser = (*PromQLParser)(nil)

// HookIn implements sdk.Parser.
func (*PromQLParser) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
	}
}

// OnCheckPayload implements sdk.Parser.
func (*PromQLParser) OnCheckPayload(*sdk.ParseCtx) (protoNum uint8, protoStr string) {
	return 1, "HTTP"
}

// OnHttpReq implements sdk.Parser.
func (*PromQLParser) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := ctx.BaseCtx
	var matchPath bool
	for _, u := range _prometheus_query_urls {
		if strings.Index(ctx.Path, u) >= 0 {
			matchPath = true
			break
		}
	}

	if !matchPath {
		return sdk.ActionNext()
	}

	payload, err := baseCtx.GetPayload()
	if err != nil {
		sdk.Error("%s parse payload error: %s", plugin_module, err)
		return sdk.ActionNext()
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		sdk.Error("%s read request error: %s", plugin_module, err)
		return sdk.ActionNext()
	}

	if req.Method != "POST" {
		return sdk.ActionNext()
	}

	header := req.Header.Clone()
	if header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		sdk.Error("promql get content-type: %s \n", header.Get("Content-Type"))
		return sdk.ActionNext()
	}

	err = req.ParseForm()
	if err != nil {
		// multi-payloads
		// sdk.Error("%s parser form error: %s", plugin_module, err)
		return sdk.ActionNext()
	}

	attrs := make([]sdk.KeyVal, 0, len(req.Form))
	for k, v := range req.Form {
		attrs = append(attrs, sdk.KeyVal{Key: k, Val: strings.Join(v, ",")})
	}
	if len(attrs) > 0 {
		return sdk.HttpReqActionAbortWithResult(nil, nil, attrs)
	} else {
		return sdk.ActionNext()
	}
}

// OnHttpResp implements sdk.Parser.
func (*PromQLParser) OnHttpResp(*sdk.HttpRespCtx) sdk.Action {
	return sdk.ActionNext()
}

// OnParsePayload implements sdk.Parser.v
func (*PromQLParser) OnParsePayload(*sdk.ParseCtx) sdk.Action {
	return sdk.ActionNext()
}
