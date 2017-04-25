package ini

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"io"
)

type Config struct {
	filepath string
	conflist[]map[string]map[string]string
}

func SetConfig(filepath string) (c *Config) {
	c = new(Config)
	c.filepath = filepath
	c.readList()
	return
}

func (c *Config) GetValue(section, name string) string {
	data := c.conflist
	for _,v := range data {
		m,ok := v[section]
		if !ok {
			return nil
		}
		return m[name]
	}
	return nil
}

func (c *Config) SetValue(section, key, value string) bool {
	data := c.conflist

	var ok bool
	var conf = make(map[string]map[string]string)

	for i,v := range data {
		_,ok = v[section]
		if ok {
			c.conflist[i][section][key] = value
			break
		}
	}
	if !ok {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		c.conflist = append(c.conflist, conf)
	}

	/*for i,v := range data {
		_,ok = v[section]
		index[i] = ok
	}
	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v = range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		c.conflist[i][section][key] = value
		return true
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		c.conflist = append(c.conflist, conf)
		return true
	}*/

	return true
}

func (c *Config) readList()[]map[string]map[string]string {
	file, err := os.Open(c.filepath)
	if err != nil {
		fmt.Sprintf("Error -> %s", err.Error())
		return nil
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for  {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				fmt.Sprintf("Error -> %s", err.Error())
				return nil
			}
			if len(line) == 0 {
				break
			}
		}
		switch  {
		case len(line) == 0:
		case line[0] == '#':
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1:len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		case strings.Contains(line, "="):
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1:len(line)-1])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniqueappend(section) == true {
				c.conflist = append(c.conflist, data)
			}
		}
	}
	return c.conflist
}

func (c *Config) uniqueappend(conf string) bool {
	for _,v:=range c.conflist {
		for k,_:=range v {
			if k==conf {
				return false
			}
		}
	}
	return true
}
