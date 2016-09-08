// Gother 御座る

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func ShowUsage() {
	fmt.Printf("Usage: gother target_directory [output_file]")
}

func GetFileNames(dirname string) (names []string) {
	dir, err := os.Open(dirname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	defer dir.Close()

	fi, err := dir.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	dir.Seek(0, 0)
	fn, _ := dir.Readdirnames(-1)

	for idx := 0; idx < len(fn); idx++ {
		fn[idx] = dirname + "/" + fn[idx]
	}
	for idx := 0; idx < len(fn); idx++ {
		if fi[idx].IsDir() {
			names = append(names, GetFileNames(fn[idx])...)
		} else {
			names = append(names, fn[idx])
		}
	}

	return
}

func WriteBytesFor(f *os.File, val interface{}) int64 {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, val)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return -1
	}

	_, err = f.Write(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return -2
	}

	return int64(binary.Size(val))
}

func Reverse(data []byte) []byte {
	for idx := 0; idx < len(data)/2; idx++ {
		opp := len(data) - 1 - idx
		data[idx], data[opp] = data[opp], data[idx]
	}
	return data
}

func main() {
	var target string
	var outfile string

	switch len(os.Args) {
	case 1:
		ShowUsage()
		os.Exit(1)
	case 2:
		target = os.Args[1]
		outfile = target + ".gzr"
	case 3:
		target = os.Args[1]
		outfile = os.Args[2]
	default:
		ShowUsage()
		os.Exit(1)
	}
	names := GetFileNames(target)
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	defer out.Close()

	var data_sizes []int64
	data_initial := int64(0)

	// 1. Amount of files
	data_initial += WriteBytesFor(out, int32(len(names)))

	// Calculate data initial position
	var bs [][]byte
	for idx := 0; idx < len(names); idx++ {
		bs = append(bs, append([]byte(names[idx]), '|'))
		data_initial += int64(len(bs[idx]) + 1 + 8 + 1 + 8 + 1)
	}

	// 2. [Filename | initial64bit size64bit]
	for idx := 0; idx < len(names); idx++ {
		WriteBytesFor(out, bs[idx])
		f, err := os.Open(names[idx])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		stat, _ := f.Stat()
		data_sizes = append(data_sizes, stat.Size())
		WriteBytesFor(out, int64(data_initial))
		WriteBytesFor(out, int64(data_sizes[idx]))
		data_initial += stat.Size()
		f.Close()
	}

	// 3. [Data]
	for idx := 0; idx < len(names); idx++ {
		f, err := os.Open(names[idx])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		data := make([]byte, data_sizes[idx])
		_, err = f.Read(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		data = Reverse(data)
		WriteBytesFor(out, data)
		f.Close()
	}

}
