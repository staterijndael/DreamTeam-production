package common

import (
	"context"
	"dt/logwrap"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"encoding/json"
	"github.com/semrush/zenrpc"
	"time"
)

type ILogger interface {
	Debug(string, ...interface{})
	Error(string, ...interface{})
	Info(string, ...interface{})
}

func Logger(handler zenrpc.InvokeFunc) zenrpc.InvokeFunc {
	return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
		req, ok := zenrpc.RequestFromContext(ctx)
		ip := "<nil>"
		if ok {
			ip = req.RemoteAddr
		}

		start, response := time.Now(), handler(ctx, method, params)
		since := time.Since(start)
		var id []byte
		var err error = nil
		id, err = zenrpc.IDFromContext(ctx).MarshalJSON()
		if err != nil {
			id = []byte("<error marshalling id: " + err.Error() + ">")
		}

		logwrap.Info("uid=%d ip=%s requestId=%s duration=%v method=%s.%s params=%s err=%s",
			requestContext.CurrentUser(ctx).ID,
			ip,
			string(id),
			since,
			zenrpc.NamespaceFromContext(ctx),
			method,
			params,
			response.Error,
		)

		logwrap.Debug("uid=%d ip=%s requestId=%s duration=%v method=%s.%s params=%s err=%s response=%s",
			requestContext.CurrentUser(ctx).ID,
			ip,
			string(id),
			since,
			zenrpc.NamespaceFromContext(ctx),
			method,
			params,
			response.Error,
			response.JSON(),
		)

		return response
	}
}

func RequestBuilder(handler zenrpc.InvokeFunc) zenrpc.InvokeFunc {
	return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
		namespace := zenrpc.NamespaceFromContext(ctx)
		m := &zenrpc.Request{
			Version:   "2.0",
			Method:    namespace + "." + method,
			ID:        nil,
			Params:    params,
			Namespace: namespace,
		}

		return handler(context.WithValue(ctx, "rpcRequest", m), method, params)
	}
}

func NicknameChecker(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
	return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
		currentUser := requestContext.CurrentUser(ctx)
		rpcRequest := ctx.Value("rpcRequest").(*zenrpc.Request)
		fullMethod := zenrpc.NamespaceFromContext(ctx) + "." + method
		if currentUser.NicknameID == nil && fullMethod != "user.setnickname" {
			return zenrpc.Response{
				Version: "2.0",
				ID:      rpcRequest.ID,
				Result:  nil,
				Error:   errors.New(errors.EnterNickname, nil, nil),
			}
		}

		return h(ctx, method, params)
	}
}
