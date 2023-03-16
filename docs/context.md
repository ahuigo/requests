# http context

## set context
    req = req.WithContext(ctx)

## client.Do(req)

    1. resp, err:= transport.RoundTrip(req)
        // http/transport.go
        1. trace := httptrace.ContextClientTrace(req.ctx)
            // net/http/httptrace/trace.go
            1. 	trace, _ := ctx.Value(clientEventContextKey{}).(*ClientTrace)
        2. treq := &transportRequest{Request: req, trace: trace, cancelKey: {req:req}}
        2.1 cm = connectMethod{targetAddr: treq.URL, targetScheme: treq.URL.Scheme}
        3. pconn, err := t.getConn(treq, cm)
        4. trace:
            //1. is called before a connection is created
            trace.GetConn(cm.targetAddr)

            // is called after a successful connection. conn: httptrace.GotConnInfo
            trace.GotConn(pc.gotIdleConnTrace(pc.idleAt))

            //ohter(..)
            // PutIdleConn is called when the connection is returned to idel pool
            PutIdleConn(err)
            ConnectStart func(network, addr string)
            ConnectDone func(network, addr string, err error)

            // WroteHeaderField is called after the Transport has written
            WroteHeaderField func(key string, value []string)
            // WroteRequest is called with the result of writing the request and any body. 
            WroteRequest func(WroteRequestInfo)

