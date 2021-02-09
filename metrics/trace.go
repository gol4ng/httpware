package metrics

import (
	"context"
	"time"
)

type InstrumentTrace struct {
	DNSStart             func(context.Context, time.Duration)
	DNSDone              func(context.Context, time.Duration)
	ConnectStart         func(context.Context, time.Duration)
	ConnectDone          func(context.Context, time.Duration)
	TLSHandshakeStart    func(context.Context, time.Duration)
	TLSHandshakeDone     func(context.Context, time.Duration)
	GotConn              func(context.Context, time.Duration)
	PutIdleConn          func(context.Context, time.Duration)
	WroteHeaders         func(context.Context, time.Duration)
	WroteRequest         func(context.Context, time.Duration)
	Got100Continue       func(context.Context, time.Duration)
	Wait100Continue      func(context.Context, time.Duration)
	GotFirstResponseByte func(context.Context, time.Duration)
}
