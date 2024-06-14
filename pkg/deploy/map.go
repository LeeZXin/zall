package deploy

import (
	"os/exec"
	"sync"
)

type cmdMap struct {
	sync.Mutex
	container map[string]*exec.Cmd
}

func newCmdMap() *cmdMap {
	return &cmdMap{
		container: make(map[string]*exec.Cmd),
	}
}

func (m *cmdMap) PutIfAbsent(id string, cmd *exec.Cmd) bool {
	m.Lock()
	defer m.Unlock()
	_, b := m.container[id]
	if b {
		return false
	}
	m.container[id] = cmd
	return true
}

func (m *cmdMap) GetById(id string) *exec.Cmd {
	m.Lock()
	defer m.Unlock()
	return m.container[id]
}

func (m *cmdMap) Remove(id string) {
	m.Lock()
	defer m.Unlock()
	delete(m.container, id)
}

func (m *cmdMap) GetAll() map[string]*exec.Cmd {
	m.Lock()
	defer m.Unlock()
	ret := make(map[string]*exec.Cmd, len(m.container))
	for k, v := range m.container {
		ret[k] = v
	}
	return ret
}
