package server

import (
	"errors"
	"github.com/id_gen/config"
	"sync"
	"time"
)

//type Code int64
func getCurrenMills() int64 {
	return time.Now().UnixNano() / 1000
}

var (
	twepoch int64 = getCurrenMills()
)

const (
	workerIdBits       int64 = 5 /** machine id */
	datacenterIdBits   int64 = 5 /**datacenter id*/
	maxWorkerId        int64 = -1 ^ (-1 << uint64(workerIdBits))
	maxDatacenterId    int64 = -1 ^ (-1 << uint64(datacenterIdBits))
	sequenceBits       int64 = 12
	workerIdShift      int64 = sequenceBits
	datacenterIdShift  int64 = sequenceBits + workerIdBits
	timestampLeftShift int64 = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       int64 = -1 ^ (-1 << uint64(sequenceBits))
)

type IdGenerator struct {
	inc           int64
	machineId     int64
	centorId      int64
	lastTimeStamp int64
	sync.Mutex
}

func (this *IdGenerator) Init(machineId, centorId int64) {
	this.inc = 0
	this.machineId = machineId
	this.centorId = centorId
	this.lastTimeStamp = getCurrenMills()
}
func (this *IdGenerator) setTimeStamp() {
	t := getCurrenMills()
	for t < this.lastTimeStamp {
		t = getCurrenMills()
	}
	this.lastTimeStamp = t
}
func (this *IdGenerator) GenId() (int64, error) {
	this.Lock()
	t := getCurrenMills()
	if t < this.lastTimeStamp {
		this.Unlock()
		return -1, errors.New("error generator")
	}
	if t == this.lastTimeStamp {
		this.inc = (this.inc + 1) & int64(sequenceMask)
		if 0 == this.inc {
			this.setTimeStamp()
		}
	} else {
		this.inc = 0
	}
	ret := ((this.lastTimeStamp - int64(twepoch)) << uint64(timestampLeftShift)) | this.centorId<<uint64(datacenterIdShift) | this.machineId<<uint64(workerIdShift) | this.inc
	this.Unlock()
	return ret, nil
}

type Server struct {
	cfg  *config.Config
	keys map[string]*IdGenerator
	sync.RWMutex
}

func NewServer(c *config.Config) (*Server, error) {
	s := new(Server)
	s.cfg = c
	s.keys = make(map[string]*IdGenerator)
	return s, nil
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) Next(k string) (int64, error) {
	var t *IdGenerator = nil
	var ok bool = false

	s.Lock()
	if t, ok = s.keys[k]; !ok {
		t = new(IdGenerator)
		t.Init(s.cfg.MachineId, s.cfg.CentorId)
		s.keys[k] = t
	}
	s.Unlock()

	if id, err := t.GenId(); err == nil {
		return id, err
	} else {
		return -1, err
	}
}
