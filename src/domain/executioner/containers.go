package executioner

import "sync"

type SecureErrContainer struct {
	content          string
	comingFromStdout bool
	mu               sync.RWMutex
	howMany          int
	maxLines         int
}

func NewSecureErrContainer(maxLines int) *SecureErrContainer {
	return &SecureErrContainer{maxLines: maxLines}
}
func (ssc *SecureErrContainer) Clean() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.UnsafeClean()
}
func (ssc *SecureErrContainer) UnsafeClean() {
	ssc.content = ""
	ssc.howMany = 0
}
func (ssc *SecureErrContainer) FromStdout() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.comingFromStdout = true
}
func (ssc *SecureErrContainer) FromStderr() {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.comingFromStdout = false
}
func (ssc *SecureErrContainer) ComingFromStdOut() bool {
	ssc.mu.RLock()
	defer ssc.mu.RUnlock()
	return ssc.comingFromStdout
}
func (ssc *SecureErrContainer) AppendValue(src string) {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	ssc.comingFromStdout = false
	ssc.content += src
	ssc.howMany++
	if ssc.howMany > ssc.maxLines {
		ssc.UnsafeClean()
	}
}

func (ssc *SecureErrContainer) Content() string {
	ssc.mu.Lock()
	defer ssc.mu.Unlock()
	return ssc.content
}
