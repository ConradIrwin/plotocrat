package data

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"
	"time"
)

type Plot struct {
	Uid        string
	Name       string
	Data       []float64
	UploadedAt time.Time
}

func Parse(name string, file io.Reader) (*Plot, error) {
	plot := &Plot{"", name, []float64{}, time.Now()}
	hash := sha256.New()
	pipe := bufio.NewReader(io.TeeReader(file, hash))

	values := 0
	errors := 0

	for {
		line, err := pipe.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		value, parseErr := strconv.ParseFloat(line, 64)

		if parseErr == nil {
			plot.Data = append(plot.Data, value)
			values += 1
			if values > 1000*1000 {
				return nil, fmt.Errorf("Too many values in file (max is 1,000,000)")
			}
		} else {
			errors += 1
			if errors > 1000 {
				return nil, fmt.Errorf("File too messy to read. Prefix lines with # to make them comments.")
			}
		}
	}

	plot.Uid = hashUid(hash)

	return plot, nil
}

func hashUid(h hash.Hash) string {
	b := make([]byte, 0, h.Size())
	h.Sum(b)
	return fmt.Sprintf("%x", b[0:10])
}
