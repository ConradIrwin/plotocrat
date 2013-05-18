package data

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"strconv"
	"strings"
	"bytes"
	"time"
	"unicode/utf8"
)

type Plot struct {
	Uid        string
	Name       string
	Data       string
	UploadedAt time.Time
}

type plotFileReader struct {
	*bufio.Reader
	values int
	errors int
}

func newPlotFileReader(input io.Reader) *plotFileReader {
	return &plotFileReader{bufio.NewReader(input), 0, 0}
}

func Parse(name string, file io.Reader) (*Plot, error) {
	plot := &Plot{Name: name, UploadedAt: time.Now()}
	hash := sha256.New()
	data := new(bytes.Buffer)

	pipe := io.TeeReader(file, data)
	pipe = io.TeeReader(pipe, hash)
	parser := newPlotFileReader(pipe)

	for {
		_, err := parser.readFloat()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if !utf8.Valid(data.Bytes()) {
		return nil, fmt.Errorf("invalid utf8 in upload")
	}

	plot.Uid = hashUid(hash)
	plot.Data = data.String()

	return plot, nil
}

func (plot *Plot) Values() []float64 {
	data := make([]float64, 0)
	parser := newPlotFileReader(strings.NewReader(plot.Data))

	for {
		value, err := parser.readFloat()
		if err == io.EOF {
			return data;
		}
		if err != nil {
			continue
		}
		data = append(data, value)
	}
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
}

func hashUid(h hash.Hash) string {
	b := make([]byte, 0, h.Size())
	h.Sum(b)
	return fmt.Sprintf("%x", b[0:10])
}
