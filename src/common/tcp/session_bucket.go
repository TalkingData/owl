package tcp

import (
	"sync"
)

type SessionBucket struct {
	m    map[string]*Session
	lock sync.RWMutex
}

func (sb *SessionBucket) Add(sess *Session) {
	sb.lock.Lock()
	sb.m[sess.RemoteAddr()] = sess
	sb.lock.Unlock()
}

func (sb *SessionBucket) Get(ip string) *Session {
	sb.lock.RLock()
	defer sb.lock.RUnlock()
	if sess, ok := sb.m[ip]; ok {
		return sess
	}
	return nil
}

func (sb *SessionBucket) Delete(ip string) {
	sb.lock.Lock()
	delete(sb.m, ip)
	sb.lock.Unlock()
}

func (sb *SessionBucket) All() map[string]*Session {
	sb.lock.RLock()
	defer sb.lock.RUnlock()
	ss := make(map[string]*Session, len(sb.m))
	for k, s := range sb.m {
		ss[k] = s
	}
	return ss
}
