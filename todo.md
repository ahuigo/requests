# Optimize client 
    t := http.DefaultTransport.(*http.Transport).Clone()
    t.MaxIdleConns = 100
    t.MaxConnsPerHost = 100
    t.MaxIdleConnsPerHost = 100
        
    httpClient = &http.Client{
      Timeout:   10 * time.Second,
      Transport: t,
    }
