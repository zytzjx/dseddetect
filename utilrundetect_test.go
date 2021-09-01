package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestParserdata(t *testing.T) {
	file, err := os.Open("saslist.txt")

	if err != nil {
		log.Fatalf("failed to open")

	}
	// The method os.File.Close() is called
	// on the os.File object to close the file
	defer file.Close()
	// The bufio.NewScanner() function is called in which the
	// object os.File passed as its parameter and this returns a
	// object bufio.Scanner which is further used on the
	// bufio.Scanner.Split() method.
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	alldisks := parseLsiData(lines)
	for _, disk := range alldisks {
		for k, v := range disk {
			fmt.Printf("k=%s, v=%s\n", k, v)
		}
	}

}
