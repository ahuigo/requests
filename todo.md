# Split client and request
client:

    Request:
        traceInfo
        rawRequest
        ...
# Optimize client 
    t := http.DefaultTransport.(*http.Transport).Clone()
    t.MaxIdleConns = 100
    t.MaxConnsPerHost = 100
    t.MaxIdleConnsPerHost = 100
        
    httpClient = &http.Client{
      // set requests timout
      Timeout:   10 * time.Second,
      Transport: t,
    }

## Support custom Transport

   http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1
    http.DefaultTransport.(*http.Transport).MaxIdleConns= 1
    r = requests.R()
    for i := 1; i < 10000; i++ {
                go  r.Get(url)
    }
    
