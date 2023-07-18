package window

import (
	"sync"
	"time"

	"github.com/dop251/goja"
)

type stop interface {
	stop()
}

type Timer struct {
	sync.Mutex
	vm     *goja.Runtime
	jobId  int
	stopCh chan struct{}
	jobCh  chan func()

	stopMap map[int]stop
}

func NewTimer(vm *goja.Runtime) *Timer {
	return &Timer{
		vm:      vm,
		stopCh:  make(chan struct{}),
		jobCh:   make(chan func()),
		stopMap: make(map[int]stop),
	}
}

func (t *Timer) Start() error {
	if err := t.vm.Set("setTimeout", t.setTimeout); err != nil {
		return err
	}
	if err := t.vm.Set("clearTimeout", t.clearTimeout); err != nil {
		return err
	}
	if err := t.vm.Set("setInterval", t.setInterval); err != nil {
		return err
	}
	if err := t.vm.Set("clearInterval", t.clearInterval); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case f := <-t.jobCh:
				f()
			case <-t.stopCh:
				return
			}
		}
	}()
	return nil
}

func (t *Timer) Stop() {
	close(t.stopCh)
	for _, v := range t.stopMap {
		v.stop()
	}
}

func (t *Timer) setTimeout(call goja.FunctionCall) goja.Value {
	return t.schedule(call, false)
}

func (t *Timer) setInterval(call goja.FunctionCall) goja.Value {
	return t.schedule(call, true)
}

func (t *Timer) clearTimeout(tm *timeout) {
	if tm != nil && !tm.canceled {
		tm.stop()
	}
}

func (t *Timer) clearInterval(i *interval) {
	if i != nil && !i.canceled {
		i.stop()
	}
}

func (t *Timer) schedule(call goja.FunctionCall, repeating bool) goja.Value {
	if fn, ok := goja.AssertFunction(call.Argument(0)); ok {
		delay := call.Argument(1).ToInteger()
		var args []goja.Value
		if len(call.Arguments) > 2 {
			args = call.Arguments[2:]
		}
		f := func() { _, _ = fn(nil, args...) }
		t.jobId++
		if repeating {
			return t.vm.ToValue(t.addInterval(f, time.Duration(delay)*time.Millisecond))
		} else {
			return t.vm.ToValue(t.addTimeout(f, time.Duration(delay)*time.Millisecond))
		}
	}
	return nil
}

type job struct {
	canceled bool
	fn       func()
}

type timeout struct {
	job
	timer *time.Timer
}

func (t *timeout) stop() {
	t.canceled = true
	t.timer.Stop()
}

func (t *Timer) addTimeout(f func(), duration time.Duration) *timeout {
	tm := &timeout{
		job: job{fn: f},
	}
	tm.timer = time.AfterFunc(duration, func() {
		t.jobCh <- func() {
			tm.fn()
		}
	})
	t.stopMap[t.jobId] = tm
	return tm
}

type interval struct {
	job
	ticker   *time.Ticker
	stopChan chan struct{}
}

func (i *interval) stop() {
	i.canceled = true
	close(i.stopChan)
}

func (t *Timer) addInterval(f func(), timeout time.Duration) *interval {
	// https://nodejs.org/api/timers.html#timers_setinterval_callback_delay_args
	if timeout <= 0 {
		timeout = time.Millisecond
	}

	i := &interval{
		job:      job{fn: f},
		ticker:   time.NewTicker(timeout),
		stopChan: make(chan struct{}),
	}

	t.stopMap[t.jobId] = i

	go func() {
		for {
			select {
			case <-i.stopChan:
				i.ticker.Stop()
				break
			case <-i.ticker.C:
				t.jobCh <- i.fn
			}
		}
	}()
	return i
}
