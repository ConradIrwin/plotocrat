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

type plotFileReader struct {
	*bufio.Reader
	values int
	errors int
}

func Parse(name string, file io.Reader) (*Plot, error) {
	plot := &Plot{"", name, []float64{}, time.Now()}
	hash := sha256.New()
	pipe := bufio.NewReader(io.TeeReader(file, hash))
	parser := &plotFileReader{pipe, 0, 0}

	for {
		value, err := parser.readFloat()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		plot.Data = append(plot.Data, value)
	}

	plot.Uid = hashUid(hash)

	return plot, nil
}

func (self plotFileReader) readFloat() (float64, error) {
	for {
		line, err := self.ReadString('\n')
		if err != nil {
			return 0, err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		value, parseErr := strconv.ParseFloat(line, 64)

		if parseErr == nil {
			self.values += 1
			if self.values > 1000*1000 {
				return 0, fmt.Errorf("Too many values in file (max is 1,000,000)")
			}
			return value, nil
		} else {
			self.errors += 1
			if self.errors > 1000 {
				return 0, fmt.Errorf("File too messy to read. Prefix lines with # to make them comments.")
			}
		}
	}
	panic("remove this in Go 1.1")
}

func hashUid(h hash.Hash) string {
	b := make([]byte, 0, h.Size())
	h.Sum(b)
	return fmt.Sprintf("%x", b[0:10])
}
