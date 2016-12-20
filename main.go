package pyradebug

import (
	"errors"
	"runtime"
	"strings"
	"sync"
)

type PyraDebug struct {
	Enable  bool
	lck     sync.Mutex
	history map[string]*GoroutineInfo
}

type GoroutineInfo struct {
	id     string
	name   string
	status string
}

func InitPyraDebug() *PyraDebug {
	pd := &PyraDebug{}
	return pd
}

func GetGoroutineId() string {
	b1 := make([]byte, 50)
	runtime.Stack(b1, false)
	s1 := string(b1)
	as1 := strings.Split(s1, "\n")
	s2 := as1[0]
	s3 := s2[:len(s2)-11][10:]
	return s3
}

func (pd *PyraDebug) SetGoroutineName(name string) {
	if !pd.Enable {
		return
	}
	id := GetGoroutineId()
	pd.lck.Lock()
	defer pd.lck.Unlock()
	// pd.history = append(pd.history, &GoroutineInfo{id: id, name: name})
	_, ok := pd.history[id]
	if ok {
		panic(errors.New("Goroutine name is set twice"))
	}
	pd.history[id] = &GoroutineInfo{id: id, name: name}
}

func parseFirstLine(s2 string) (id string, status string) {
	i1 := strings.Index(s2, "[")
	id = s2[:i1-1]
	i2 := strings.Index(s2, "]")
	status = s2[i1+1 : i2]
	return id, status
}

func (pd *PyraDebug) ListGoroutines(bufferSize int) (result []*GoroutineInfo) {
	b1 := make([]byte, bufferSize)
	runtime.Stack(b1, true)
	s1 := string(b1)
	as1 := strings.Split(s1, "\n")
	gs := [][]string{nil}
	for _, s := range as1 {
		if s == "" {
			gs = append(gs, nil)
		} else {
			gs[len(gs)-1] = append(gs[len(gs)-1], s)
		}
	}

	for _, as2 := range gs {
		if len(as2) == 0 {
			continue
		}
		id, status := parseFirstLine(as2[0])
		result = append(result, &GoroutineInfo{id: id, status: status})
	}

	pd.GiveNames(result)

	return result
}
func (pd *PyraDebug) GiveNames(a1 []*GoroutineInfo) {
	pd.lck.Lock()
	defer pd.lck.Unlock()

	for _, gi := range a1 {
		gi2 := pd.history[gi.id]
		if gi2 != nil {
			gi.name = gi2.name
		}
	}
}
