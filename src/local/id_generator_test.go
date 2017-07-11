package server

import (
	"github.com/id_gen/config"
	"sync"
	"testing"
	"time"
)

func Test_t(t *testing.T) {

	var wg sync.WaitGroup
	cfg, err := config.ParseConfigFile("../../config/test.json")
	if err != nil {
		t.Error("parse config error")
		return
	}
	ch := make(chan int64, 10)
	t.Log(cfg.CentorId, "\t", cfg.MachineId, "\t", "\t", "time Now :=", "\t", time.Now().UnixNano()/int64(1000))
	server, _ := NewServer(cfg)
	for i := 0; i < 10; i++ {
		go func(s *Server) {
			wg.Add(1)
			id, _ := server.Next("test")
			ch <- id
			wg.Done()
		}(server)
	}
	go func() {
		<-ch
	}()
	wg.Wait()
}
