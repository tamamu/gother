// Gother 御座る

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var Order = binary.LittleEndian

func ReadInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, Order, &ret)
	return
}

func ReadInt64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, Order, &ret)
	return
}

type GotherDataRange struct {
	Initial int64
	Size    int64
}

type GotherFile struct {
	File     *os.File
	RangeMap map[string]GotherDataRange
}

func (gf GotherFile) GetData(name string) (data []byte, err error) {
	rg, ok := gf.RangeMap[name]
	if ok == false {
		err = errors.New(name + " not found.")
		return
	}
	gf.File.Seek(rg.Initial, 0)
	data = make([]byte, rg.Size)
	_, err = gf.File.Read(data)
	return
}

func GetHeader(file *os.File) (head []byte, err error) {
	// 1. Get header size
	file.Seek(0, 0)
	size := make([]byte, 4)
	_, err = file.Read(size)
	if err != nil {
		return
	}

	// 2. Read header
	file.Seek(0, 0)
	head = make([]byte, ReadInt32(size))
	_, err = file.Read(head)
	return
}

func Parse(bytes []byte) (datas map[string]GotherDataRange) {
	// 1. Read amount of files
	amount := int(ReadInt32(bytes[4:8]))

	// Next byte position
	idx := 8

	// 2. Read [filename \0 start end]
	datas = make(map[string]GotherDataRange)
	for i := 0; i < amount; i++ {
		// Read a file name for get \0
		var name []byte
		for bytes[idx] != 0 {
			name = append(name, bytes[idx])
			idx++
		}

		idx += 1 // length of \0
		initial := ReadInt64(bytes[idx : idx+8])
		idx += 8 // length of initial 64bit
		last := ReadInt64(bytes[idx : idx+8])
		idx += 8 // length of size 64bit

		datas[string(name)] = GotherDataRange{initial, last}
	}

	return datas
}

func Open(path string) (gf GotherFile, err error) {
	gf.File, err = os.Open(path)
	if err != nil {
		return
	}
	head, err := GetHeader(gf.File)
	if err != nil {
		return
	}
	gf.RangeMap = Parse(head)
	return
}

func ShowUsage() {
	fmt.Printf("Usage: decode gother_file target_file [result_file]")
}

func main() {
	var target string
	var inputfile string
	var resultfile string

	switch len(os.Args) {
	default:
		ShowUsage()
		os.Exit(1)
	case 3:
		inputfile = os.Args[1]
		target = os.Args[2]
		resultfile = "result.out"
	case 4:
		inputfile = os.Args[1]
		target = os.Args[2]
		resultfile = os.Args[3]
	}

	gf, err := Open(inputfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data, err := gf.GetData(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create(resultfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = out.Write(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
