package io

type IOKey struct {	
	key string
}

func NewIOKey(key string) *IOKey {
	return &IOKey{
		key: key,
	}
}