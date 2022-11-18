package trace

import (
	"context"
	"net/http"
	"time"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/gin"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type FireTraceService struct {
	idService        contract.IDService
	traceIDGenerator contract.IDService
	spanIDGenerator  contract.IDService
}

func NewFireTraceService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	idService := container.MustMake(contract.IDKey).(contract.IDService)
	return &FireTraceService{idService: idService}, nil
}

// WithTrace register new trace to context
func (t *FireTraceService) WithTrace(ctx context.Context, trace *contract.TraceContext) context.Context {
	if ginC, ok := ctx.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(ctx, ContextKey, trace)
		return newC
	}
}

// GetTrace From trace context
func (t *FireTraceService) GetTrace(c context.Context) *contract.TraceContext {
	if ginC, ok := c.(*gin.Context); ok {
		if val, ok2 := ginC.Get(string(ContextKey)); ok2 {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := c.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}
	return nil
}

// NewTrace generate a new trace
func (t *FireTraceService) NewTrace() *contract.TraceContext {
	var traceID, spanID string
	if t.traceIDGenerator != nil {
		traceID = t.traceIDGenerator.NewID()
	} else {
		traceID = t.idService.NewID()
	}

	if t.spanIDGenerator != nil {
		spanID = t.spanIDGenerator.NewID()
	} else {
		spanID = t.idService.NewID()
	}
	tc := &contract.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CspanID:    "",
		Annotation: map[string]string{},
	}
	return tc
}

// ChildSpan instance a sub trace with new span id
func (t *FireTraceService) ChildSpan(tc *contract.TraceContext) *contract.TraceContext {
	var childSpanID string
	if t.spanIDGenerator != nil {
		childSpanID = t.spanIDGenerator.NewID()
	} else {
		childSpanID = t.idService.NewID()
	}
	childSpan := &contract.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CspanID:  childSpanID,
		Annotation: map[string]string{
			contract.TraceKeyTime: time.Now().String(),
		},
	}
	return childSpan
}

// ExtractHTTP GetTrace By Http
func (t *FireTraceService) ExtractHTTP(req *http.Request) *contract.TraceContext {
	tc := &contract.TraceContext{}
	tc.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract.TraceKeySpanID)
	tc.SpanID = req.Header.Get(contract.TraceKeyCspanID)
	tc.CspanID = ""
	if tc.TraceID == "" {
		tc.TraceID = t.idService.NewID()
	}
	if tc.SpanID == "" {
		tc.SpanID = t.idService.NewID()
	}
	return tc
}

// InjectHTTP Set Trace to Http
func (t *FireTraceService) InjectHTTP(req *http.Request, tc *contract.TraceContext) *http.Request {
	req.Header.Add(contract.TraceKeyTraceID, tc.TraceID)
	req.Header.Add(contract.TraceKeySpanID, tc.SpanID)
	req.Header.Add(contract.TraceKeyCspanID, tc.CspanID)
	req.Header.Add(contract.TraceKeyParentID, tc.ParentID)
	return req
}

func (t *FireTraceService) ToMap(tc *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}
	m[contract.TraceKeyTraceID] = tc.TraceID
	m[contract.TraceKeySpanID] = tc.SpanID
	m[contract.TraceKeyCspanID] = tc.CspanID
	m[contract.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}
	return m
}
