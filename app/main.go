package main
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)
const FILE_HEADER_SIZE = 100
// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]
	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}
	//	header := make([]byte, 100)
		header := make([]byte, FILE_HEADER_SIZE)
		_, err = databaseFile.Read(header)
		if err != nil {
			log.Fatal(err)
		}
		var pageSize uint16
		if err := binary.Read(bytes.NewReader(header[16:18]), binary.BigEndian, &pageSize); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}
		schemaBuffer := make([]byte, pageSize-FILE_HEADER_SIZE)
		_, err = databaseFile.ReadAt(schemaBuffer, FILE_HEADER_SIZE)
		if err != nil {
			log.Fatal(err)
		}
		var tableCount uint16
		if err := binary.Read(bytes.NewReader(schemaBuffer[3:5]), binary.BigEndian, &tableCount); err != nil {
			fmt.Println("Failed to read integer:", err)
			return
		}
		// You can use print statements as follows for debugging, they'll be visible when running tests.
		fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")
		// Uncomment this to pass the first stage
		fmt.Printf("database page size: %v", pageSize)
		fmt.Printf("number of tables: %v", tableCount)
	
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
