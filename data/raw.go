package data

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
    "hash"
)

type Raw struct {
    Uid         string
	Name        string
	Data        string
	Uploaded_at time.Time
}

func parse(input io.Reader, on_value func(float64) error, on_error func(string) error) error {
	line_input := bufio.NewReader(input)
	for {
		line, err := line_input.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		value, parse_err := strconv.ParseFloat(line, 64)

		if parse_err != nil {
			if on_error != nil {
				err = on_error(line)
			}
		} else {
			err = on_value(value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func hash_uid(h hash.Hash) string {
    b := make([]byte, 0, h.Size());
    h.Sum(b);
    return fmt.Sprintf("%x", b[0:10]);
}

func Upload(name string, file io.Reader) (*Raw, error) {
	raw := &Raw{"", name, "", time.Now()}
	data := new(bytes.Buffer)
	hash := sha256.New()

	values := 0
	errors := 0

	pipe := io.TeeReader(file, data)
	pipe = io.TeeReader(pipe, hash)

	err := parse(pipe, func(value float64) error {
		values += 1
		if values > 1000*1000 {
			return fmt.Errorf("Too many values in file (max is 1,000,000)")
		}
		return nil

	}, func(line string) error {
		errors += 1
		if errors > 1000 {
			return fmt.Errorf("File too messy to read. Prefix lines with # to make them comments.")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

    raw.Uid = hash_uid(hash)
	raw.Data = data.String()

	return raw, nil
}

func (r Raw) Series() []float64 {

	series := []float64{}

	err := parse(strings.NewReader(r.Data), func(value float64) error {
		series = append(series, value)
		return nil
	}, nil)

	if err != nil {
		panic(err)
	}

	return series
}
