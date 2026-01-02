package domain

import "sync"

type SecureStringContainer struct {
	content string
	mu      sync.RWMutex
}

func NewSecureStrContainer() *SecureStringContainer {
	return &SecureStringContainer{}
}
func (ssc *SecureStringContainer) Clean() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.content = ""
}

func (ssc *SecureStringContainer) AppendValue(src string) {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.content += src
}

func (ssc *SecureStringContainer) Content() string {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	return ssc.content
}
