package component

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/friedelschoen/st8/config"
	"github.com/friedelschoen/st8/notify"
	"github.com/friedelschoen/st8/proto"
)

var attrPattern = regexp.MustCompile(`^!(\w+)=(\S*|"[^"]*")\s*`)

func HookComponent(path string) ComponentBuilder {
	return func(args map[string]string, events *proto.EventHandlers) (Component, error) {
		cmd := exec.Command(path)
		cmd.Env = os.Environ()
		for k, v := range args {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Stdin = nil
		cmd.Stderr = os.Stderr
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		if err := cmd.Start(); err != nil {
			return nil, err
		}

		return func(block *proto.Block, not *notify.Notification) error {
			if buffer.Len() > 0 {
				text := buffer.Bytes()
				if idx := bytes.LastIndexByte(text, '\n'); idx != -1 {
					text = text[idx+1:]
				}

				attrs := make(map[string]string)
				for len(text) > 0 {
					m := attrPattern.FindSubmatch(text)
					if m == nil {
						break
					}
					key := string(m[1])
					value := string(m[2])
					attrs[key] = value
					text = text[len(m[0]):]
				}
				attrs["full_text"] = string(text)
				config.UnmarshalConf(attrs, "", block)
				buffer.Reset()
			}

			return nil
		}, nil
	}
}
