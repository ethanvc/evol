package httpkit

func SendRequest(sa *SingleAttempt, req, resp any) error {
	return DefaultClient.SendRequest(sa, req, resp)
}

var DefaultClient = &HttpClient{}
