package qemu

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type QMPClient struct {
	socketPath string
}

func NewQMPClient(socketPath string) *QMPClient {
	return &QMPClient{socketPath: socketPath}
}

func (c *QMPClient) execute(command string, args interface{}) (map[string]interface{}, error) {
	conn, err := net.DialTimeout("unix", c.socketPath, 2*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// QMP Greeting
	decoder := json.NewDecoder(conn)
	var greeting map[string]interface{}
	if err := decoder.Decode(&greeting); err != nil {
		return nil, err
	}

	// qmp_capabilities
	cmd := map[string]interface{}{"execute": "qmp_capabilities"}
	if err := json.NewEncoder(conn).Encode(cmd); err != nil {
		return nil, err
	}
	var res map[string]interface{}
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}

	// Actual command
	cmd = map[string]interface{}{"execute": command}
	if args != nil {
		cmd["arguments"] = args
	}
	if err := json.NewEncoder(conn).Encode(cmd); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&res); err != nil {
		return nil, err
	}

	if errRes, ok := res["error"]; ok {
		return nil, fmt.Errorf("QMP error: %v", errRes)
	}

	return res, nil
}

func (c *QMPClient) Pause() error {
	_, err := c.execute("stop", nil)
	return err
}

func (c *QMPClient) Resume() error {
	_, err := c.execute("cont", nil)
	return err
}

func (c *QMPClient) PowerDown() error {
	_, err := c.execute("system_powerdown", nil)
	return err
}

func (c *QMPClient) SaveSnapshot(name string) error {
	_, err := c.execute("savevm", map[string]string{"name": name})
	return err
}

func (c *QMPClient) LoadSnapshot(name string) error {
	_, err := c.execute("loadvm", map[string]string{"name": name})
	return err
}

func (c *QMPClient) DeleteSnapshot(name string) error {
	_, err := c.execute("delvm", map[string]string{"name": name})
	return err
}
