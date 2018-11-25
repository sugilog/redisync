package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
)

var (
	LB = func() []byte {
		b := make([]byte, 2)
		b[0] = '\r'
		b[1] = '\n'
		return b
	}()
)

type Operation struct {
	Args       []string
	LenOpe     int
	LenNextArg int
	Buf        []byte
}

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(conn, "SYNC\r\n")

	reader := bufio.NewReader(conn)
	var ope Operation

	for {
		if line, err := reader.ReadBytes('\n'); err != nil {
			panic(err)
		} else {
			err = ope.Parse(line)

			if err != nil {
				panic(err)
			}

			if ope.LenOpe == 0 {
				fmt.Println(ope)
				ope = Operation{}
			}
		}
	}
}

func (ope *Operation) Parse(line []byte) error {
	if ope.LenOpe > 0 {
		if n, err := ParseStartOfArg(line); err != nil {
			return err
		} else {
			if n > 0 {
				ope.LenNextArg = n
			} else {
				arg := bytes.TrimSuffix(line, LB)
				n := len(arg)
				ope.LenNextArg -= n

				if ope.LenNextArg == 0 {
					arg := append(ope.Buf, arg...)
					ope.Args = append(ope.Args, string(arg))
					ope.Buf = ope.Buf[:0]
					ope.LenOpe -= 1
				} else {
					ope.Buf = append(ope.Buf, arg...)
				}
			}
		}
	} else {
		if n, err := ParseStartOfOperation(line); err != nil {
			return err
		} else {
			if n > 0 {
				ope.LenOpe = n
			}
		}
	}

	return nil
}

// RESP
// https://redis.io/topics/protocol

func ParseStartOfOperation(line []byte) (int, error) {
	// ex:
	// [42 49 13 10] "*1\r\n"
	// [42 50 13 10] "*2\r\n"
	if line[0] == 42 {
		return strconv.Atoi(string(bytes.TrimSuffix(line, LB)[1:]))
	} else {
		return -1, nil
	}
}

func ParseStartOfArg(line []byte) (int, error) {
	// ex:
	// [36 51 13 10] "$3\r\n"
	if line[0] == 36 {
		return strconv.Atoi(string(bytes.TrimSuffix(line, LB)[1:]))
	} else {
		return -1, nil
	}
}
