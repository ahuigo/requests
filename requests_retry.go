package requests

import "time"

func (session *Session) SetRetryCount(count int) *Session {
	session.retryCount = count
	return session
}
func (session *Session) SetRetryWaitTime(waitTime time.Duration) *Session {
	session.retryWaitTime = waitTime
	return session
}
func (session *Session) SetRetryCondition(fn func(*Response, error) bool) *Session {
	session.retryConditionFunc = fn
	return session
}