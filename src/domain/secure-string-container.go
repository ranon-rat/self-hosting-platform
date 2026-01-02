package domain

import "sync"

type SecureStringContainer struct {
	content          string
	comingFromStdout bool
	mu               sync.RWMutex
}

func NewSecureStrContainer() *SecureStringContainer {
	return &SecureStringContainer{}
}
func (ssc *SecureStringContainer) Clean() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()

	ssc.content = ""
}
func (ssc *SecureStringContainer) FromStdout() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.comingFromStdout = true
}
func (ssc *SecureStringContainer) FromStderr() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.comingFromStdout = false
}
func (ssc *SecureStringContainer) ComingFromStdOut() bool {
	ssc.mu.RLock()
	defer ssc.mu.RUnlock()
	return ssc.comingFromStdout
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
