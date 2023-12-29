package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	_ "github.com/wasilibs/nottinygc"
)

const plugin_module = "header_extraction"

func main() {
	sdk.Warn(fmt.Sprintf("%s loaded", plugin_module))
	cfg := &Config{}
	cfg.init("/etc/deepflow-agent/header.yaml")
	sdk.SetParser(&HeaderParser{cfg: cfg})
}

type HeaderParser struct {
	cfg *Config
}

var _ sdk.Parser = (*HeaderParser)(nil)

func (p *HeaderParser) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
	}
}

func (p *HeaderParser) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := ctx.BaseCtx
	if !p.cfg.allowCapturePort(baseCtx.DstPort) {
		return sdk.ActionNext()
	}

	for _, v := range p.cfg.ProcName {
		if m, err := regexp.Match(v, []byte(baseCtx.ProcName)); err != nil || !m {
			return sdk.ActionNext()
		}
	}

	payload, err := baseCtx.GetPayload()
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		return sdk.ActionAbortWithErr(err)
	}
	headers := req.Header.Clone()
	attrs := make([]sdk.KeyVal, 0, len(headers))
	for k, v := range headers {
		attrs = append(attrs, sdk.KeyVal{
			Key: k, Val: strings.Join(v, ", "),
		})
	}
	if len(attrs) > 0 {
		return sdk.HttpReqActionAbortWithResult(nil, nil, attrs)
	} else {
		return sdk.ActionNext()
	}
}

func (p *HeaderParser) OnHttpResp(ctx *sdk.HttpRespCtx) sdk.Action {
	return sdk.ActionNext()
}

func (p *HeaderParser) OnCheckPayload(ctx *sdk.ParseCtx) (uint8, string) {
	// 这里是协议判断的逻辑， 返回 0 表示失败
	// return 0, ""
	return 1, "some protocol"
}

func (p *HeaderParser) OnParsePayload(ctx *sdk.ParseCtx) sdk.Action {
	return sdk.ActionNext()
}
