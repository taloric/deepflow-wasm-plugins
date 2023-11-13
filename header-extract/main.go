package main

import (
	"bufio"
	"bytes"
	"net/http"
	"strings"

	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
)

func main() {
	sdk.Warn("header extraction plugin loaded")
	sdk.SetParser(&HeaderParser{
		port_whilelist: []uint16{14317, 14318}, // only parse 14317 and 14318 port data now
	})
}

type HeaderParser struct {
	port_whilelist []uint16
}

func (p *HeaderParser) HookIn() []sdk.HookBitmap {
	return []sdk.HookBitmap{
		sdk.HOOK_POINT_HTTP_REQ,
	}
}

func (p *HeaderParser) OnHttpReq(ctx *sdk.HttpReqCtx) sdk.Action {
	baseCtx := ctx.BaseCtx
	for po := range p.port_whilelist {
		if baseCtx.DstPort != p.port_whilelist[po] {
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
