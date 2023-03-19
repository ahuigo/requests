# http context

## set context
    
    // import net/http/httptrace
    trace := &httptrace.ClientTrace{
        ConnectStart: func(network, addr string) {
            fmt.Println(time.Now(),"ConnectStart:", "network=",network, "addr", addr)
        },
        ConnectDone: func(network, addr string, err error) {
            fmt.Println(time.Now(),"ConnectDone:", "network=",network, "addr", addr)
        },
    }
    ctx := httptrace.WithClientTrace(ctx=context.Background(), trace)
        // net/http/httptrace/trace.go:41
        1. oldTrace := ContextClientTrace(ctx)
        2. trace.compose(oldTrace)
        3. return ctx = context.WithValue(ctx, clientEventContextKey{}, trace)
    req = new(Request).WithContext(ctx)


What is `trace.compose(oldTrace)` doing?
    
    // net/http/httptrace/trace.go:179
    tv := reflect.ValueOf(trace *ClientTrace).Elem()
	ov := reflect.ValueOf(old *ClientTrace).Elem()
    for i := 0; i < tv.Type().NumField(); i++ {
        tf := tv.Field(i)
        of := ov.Field(i)

        //filter function
        if tf.Type().Kind!=reflect.Func{continue} 

        //filter nil
        if of.IsNil(){continue} 

        // makeCopy: (Otherwise it creates a recursive call cycle )
        tfCopy := reflect.ValueOf(tf.Interface())

        // wrap
        newFunc := reflect.MakeFunc(hookType, func(args []reflect.Value) []reflect.Value {
			tfCopy.Call(args)
			return of.Call(args)
		})
		tv.Field(i).Set(newFunc)
    }


    

## client.Do(req)
client.Do(req)

    // http/client.go:715
    1. c.send(req, deadline= client.Timeout)
    // http/client.go:175
    2. resp, didTimeout, err = send(req, rt=c.transport(), deadline)
    // http/client.go:251
    3. resp, err = rt.RoundTrip(req)
    1. resp, err:= transport.RoundTrip(req)
        // http/transport.go:507
        1. trace := httptrace.ContextClientTrace(req.Context())
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

