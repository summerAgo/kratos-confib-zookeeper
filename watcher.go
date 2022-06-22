package zookeeper

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
	"log"
)

type watcher struct {
	source *source
	evt    chan zk.Event
}

func newWatcher(s *source) *watcher {
	w := &watcher{
		source: s,
		evt:    make(chan zk.Event, 1),
	}

	go func() {
		for {
			_, _, event, err := w.source.client.ExistsW(w.source.options.path)
			if err != nil {
				log.Printf("watch children path %s failed, err: %s", w.source.options.path, err.Error())
			}
			log.Printf("watch children path %s successful!", w.source.options.path)

			select {
			case e := <-event:
				w.evt <- e
			}
		}
	}()
	return w
}

func (s *watcher) Next() ([]*config.KeyValue, error) {
	for {
		select {
		case event := <-s.evt:
			log.Printf("path is %s", event.Path)
			return s.source.Load()
		}
	}
}

func (s *watcher) Stop() error {
	s.source.client.Close()
	return nil
}
