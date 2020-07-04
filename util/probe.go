package util

import "sync"

type Probe struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func NewProbe(wg *sync.WaitGroup, ch chan struct{}) *Probe {
	return &Probe{wg: wg, ch: ch}
}

func (p *Probe) Done() {
	p.wg.Done()
}

func (p *Probe) Chan() <-chan struct{} {
	return p.ch
}
