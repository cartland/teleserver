package lib

import "time"

const (
	// Maximum time to try broadcasting to a subscriber.
	broacastWait = time.Second
)

type Broadcaster struct {
	chs         []chan<- Metric
	cast        chan Metric
	join, leave chan chan Metric
}

// NewBroadcaster creates a new broadcaster that broadcasts a message to many
// subscribers.
func NewBroadcaster() *Broadcaster {
	b := &Broadcaster{
		cast:  make(chan Metric),
		join:  make(chan chan Metric),
		leave: make(chan chan Metric),
	}
	go b.run()
	return b
}

func trySend(ch chan<- Metric, m Metric) {
	select {
	case <-time.After(broacastWait):
	case ch <- m:
	}
}

func (b *Broadcaster) run() {
	for {
		select {

		case m := <-b.join:
			b.chs = append(b.chs, m)

		case m := <-b.leave:
			for i, ch := range b.chs {
				if m == ch {
					b.chs = append(b.chs[:i], b.chs[i+1:]...)
				}
			}

		case m := <-b.cast:
			for _, ch := range b.chs {
				go trySend(ch, m)
			}

		}
	}
}

// Subscribe will join the broadcasting of metrics, with an argument that causes
// the subr to leave the broadcast once
func (b *Broadcaster) Subscribe(done <-chan struct{}) <-chan Metric {
	ch := make(chan Metric)
	b.join <- ch
	go func() { <-done; b.leave <- ch }()
	return ch
}

// Cast sends a metric to be broadcast to all subscribers.
func (b *Broadcaster) Cast(m Metric) { b.cast <- m }
