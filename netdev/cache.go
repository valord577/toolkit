package netdev

import "net"

type set[T comparable] map[T]struct{}

func (s set[T]) set(t T) {
	s[t] = struct{}{}
}

func (s set[T]) del(t T) {
	delete(s, t)
}

func (s set[T]) len() int {
	return len(s)
}

func (s set[T]) has(t T) bool {
	_, ok := s[t]
	return ok
}

type cachedMapNetAddr map[string]set[string]

func (c cachedMapNetAddr) add(net string, addrs []net.Addr) {
	value, ok := c[net]
	if !ok || value == nil {
		c[net] = make(set[string], len(addrs))
	}
	for _, netAddr := range addrs {
		c[net].set(netAddr.String())
	}
}

func (c cachedMapNetAddr) eq(other cachedMapNetAddr) bool {
	if len(c) != len(other) {
		return false
	}
	for k, addrs1 := range c {
		addrs2, ok := other[k]
		if !ok {
			return false
		}
		if addrs1.len() != addrs2.len() {
			return false
		}

		for addr2 := range addrs2 {
			ok = addrs1.has(addr2)
			if !ok {
				return false
			}
		}
	}
	return true
}
