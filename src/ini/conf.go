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
	confmap map[string]map[string]string
}

func SetConfig(filepath string) (c *Config) {
	c = new(Config)
	c.filepath = filepath
	c.readList()
	return
}

func (c *Config) GetValue(section, name string) string {
	return c.confmap[section][name]
}

func (c *Config) GetSection(section string) map[string]string {
	return c.confmap[section]
}

func (c *Config) SetValue(section, key, value string) bool {

	if _,ok := c.confmap[section];ok {
		c.confmap[section][key] = value
		return true
	} else {
		c.confmap[section] = make(map[string]string)
		c.confmap[section][key] = value
		return false
	}
}

func (c *Config) readList() map[string]map[string]string {
	file, err := os.Open(c.filepath)
	if err != nil {
		fmt.Sprintf("Error -> %s", err.Error())
		return nil
	}
	defer file.Close()
	c.confmap = make(map[string]map[string]string)
	var section string
	buf := bufio.NewReader(file)
	for  {
		strLine, err := buf.ReadString('\n')
		line := strings.TrimSpace(strLine)
		if err != nil {
			if err != io.EOF {
				fmt.Sprintf("Error -> %s", err.Error())
				return nil
			}
			break
		}
		switch  {
		case len(line) == 0:
		case line[0] == '#':
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1:len(line)-1])
			if _,ok := c.confmap[section]; !ok {
				c.confmap[section] = make(map[string]string)
			}
			fmt.Sprintf("Error -> A same section named %s has exist!", section)
		case strings.Contains(line, "="):
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1:])
			c.confmap[section][strings.TrimSpace(line[0:i])] = value
		}
	}
	return c.confmap
}
