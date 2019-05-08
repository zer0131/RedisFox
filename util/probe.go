package util

import "sync"

type Probe struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func NewProbe(wg *sync.WaitGroup, ch chan struct{}) *Probe {
	return &Probe{wg: wg, ch:ch}
}

func (this *Probe) Done()  {
	this.wg.Done()
}

func (this *Probe) Chan() <-chan struct{} {
	return this.ch
}
