package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// main speaks a JSON-lines protocol on stdin/stdout.
func main() {
	lines := make(chan []byte, 1024)
	go func() {
		r := bufio.NewReaderSize(os.Stdin, 1<<20)
		for {
			line, err := r.ReadBytes('\n')
			if len(line) > 1 {
				lines <- line
			}
			if err != nil {
				close(lines)
				return
			}
		}
	}()

	out := bufio.NewWriterSize(os.Stdout, 1<<16)
	send := func(v any) {
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		out.Write(b)
		out.WriteByte('\n')
		out.Flush()
	}

	var bot *Bot
	var lastAim [2]float64
	for {
		line, ok := <-lines
		if !ok {
			return
		}
		// Drain everything queued so we only act on the newest observation.
		batch := [][]byte{line}
	drain:
		for {
			select {
			case l, ok := <-lines:
				if !ok {
					break drain
				}
				batch = append(batch, l)
			default:
				break drain
			}
		}

		var latest []byte
		for _, raw := range batch {
			var head struct {
				Type string `json:"type"`
			}
			if json.Unmarshal(raw, &head) != nil {
				continue
			}
			switch head.Type {
			case "hello":
				var h Hello
				if json.Unmarshal(raw, &h) == nil {
					bot = NewBot(&h)
				}
				send(map[string]string{"type": "ready", "name": "dilksh"})
			case "obs":
				latest = raw
			case "end":
				return
			}
		}
		if latest == nil || bot == nil {
			continue
		}
		var o Obs
		if err := json.Unmarshal(latest, &o); err != nil {
			fmt.Fprintln(os.Stderr, "bad obs:", err)
			continue
		}
		act := bot.Act(&o, lastAim)
		lastAim = act.Aim
		send(act)
	}
}
