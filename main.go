package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	var result string
	if len(os.Args) == 2 {
		result = os.Args[1]
	} else {
		// list all files in current folder
		files, err := filepath.Glob("*.ini.dat")
		if err != nil {
			log.Fatalln(err)
		}

		// prompt for a file
		prompt := promptui.Select{
			Label: "Select a file",
			Items: files,
		}

		_, result, err = prompt.Run()
		if err != nil {
			log.Fatalln("Prompt failed", err)
		}
	}

	// open file for reading
	fin, err := os.Open(result)
	if err != nil {
		log.Fatalln(err)
	}
	defer fin.Close()

	// create file for writing
	name := strings.TrimSuffix(result, filepath.Ext(result))
	fout, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer fout.Close()

	// encryption mask
	var mask = []byte{0xA1, 0xB2, 0xAA, 0x12, 0x23, 0xF1, 0xF3, 0xD3, 0x78, 0x02}
	buf := make([]byte, 1024)

	br := bufio.NewReader(fin)
	bw := bufio.NewWriter(fout)
	for {
		n, err := br.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		// apply mask
		for i := 0; i < n; i++ {
			buf[i] = buf[i] ^ mask[i%len(mask)]
		}
		bw.Write(buf[:n])
	}
	bw.Flush()

}
