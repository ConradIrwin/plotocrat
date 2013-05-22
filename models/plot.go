package models

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

type plotValidator struct {
	*bufio.Reader
	values int
	errors int
}

func newPlotValidator(input io.Reader) *plotValidator {
	return &plotValidator{bufio.NewReader(input), 0, 0}
}

func NewPlot(name string, file io.Reader) (*Plot, error) {
	plot := &Plot{Name: name, UploadedAt: time.Now()}
	hash := sha256.New()
	data := new(bytes.Buffer)

	pipe := io.TeeReader(file, data)
	pipe = io.TeeReader(pipe, hash)
	validator := newPlotValidator(pipe)

	for {
		err := validator.checkNextLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if validator.values == 0 {
		return nil, fmt.Errorf("Must provide at least one number")
	}

	if !utf8.Valid(data.Bytes()) {
		return nil, fmt.Errorf("Invalid UTF-8 in upload")
	}

	plot.Uid = hashUid(hash)
	plot.Data = data.String()

	return plot, nil
}

func (plot *Plot) Values() []float64 {
	data := make([]float64, 0)
	reader := bufio.NewReader(strings.NewReader(plot.Data))

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return data;
		}
		value, err := parseFloat(line)
		if err == nil {
			data = append(data, value)
		}
	}
}

func (self *plotValidator) checkNextLine() error {
	line, err := self.ReadString('\n')
	if err != nil {
		return err
	}

	if isBlankOrComment(line) {
		return nil
	}

	_, err = parseFloat(line)

	if err == nil {
		self.values += 1
		if self.values > 1000*1000 {
			return fmt.Errorf("Too many values in file (max is 1,000,000)")
		}
	} else {
		self.errors += 1
		if self.errors > 1000 {
			return fmt.Errorf("File too messy to read. Prefix lines with # to make them comments.")
		}
	}
	return nil
}

func isBlankOrComment(line string) bool {
	line = strings.TrimSpace(line)
	return len(line) == 0 || line[0] == '#'
}

func parseFloat(line string) (float64, error) {
	line = strings.TrimSpace(line)
	return strconv.ParseFloat(line, 64)
}

func hashUid(h hash.Hash) string {
	b := make([]byte, 0, h.Size())
	h.Sum(b)
	return fmt.Sprintf("%x", b[0:10])
}
