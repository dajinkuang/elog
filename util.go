package elog

import (
	"context"
	"fmt"
)

const __CtxElogOrderMapKey = "ctx_elog_order_map_key"

func SetTraceInfo(ctx context.Context, traceId, parentId, spanId string) context.Context {
	om := NewOrderMap()
	om.Set(TraceId, traceId)
	om.Set(ParentId, parentId)
	om.Set(SpanId, spanId)
	src := FromContext(ctx)
	if src == nil {
		src = NewOrderMap()
	}
	src.AddVals(om)
	return setContext(ctx, src)
}

// copy ctx 中的 日志的trace信息 GogoLanguage,GogoPlatform,GogoDeviceId,GogoVersion,GogoUserToken
// 其它的全部丢弃，比如超时设置等
func CopyTraceInfo(ctx context.Context) context.Context {
	src := FromContext(ctx)
	if src == nil {
		src = NewOrderMap()
	}
	return setContext(context.Background(), src)
}

func GetTraceInfo(ctx context.Context) (traceId, parentId, spanId string) {
	om := FromContext(ctx)
	if tmp, ok := om.Get(TraceId); ok {
		traceId = tmp.(string)
	}
	if tmp, ok := om.Get(ParentId); ok {
		parentId = tmp.(string)
	}
	if tmp, ok := om.Get(SpanId); ok {
		spanId = tmp.(string)
	}
	return
}

func FromContext(ctx context.Context) *OrderedMap {
	ret := ctx.Value(__CtxElogOrderMapKey)
	if ret == nil {
		return nil
	}
	return ret.(*OrderedMap)
}

// 别人不需要用
func setContext(ctx context.Context, dt *OrderedMap) context.Context {
	ctx = context.WithValue(ctx, __CtxElogOrderMapKey, dt)
	return ctx
}

func ValueFromOM(ctx context.Context, key interface{}) interface{} {
	src := FromContext(ctx)
	if src == nil {
		return nil
	}
	val, ok := src.Get(fmt.Sprintf("%v", key))
	if !ok {
		return nil
	}
	return val
}
