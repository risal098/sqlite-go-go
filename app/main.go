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
//regex create table and .(
//43 52 45  41 54 45 20 54 41 42 4c 45 20 ....0a /28

/*
100 db header 
then btree header, total size 12
	offset 3,size 2 is number of tables or cell
	offset 5, size 2 is pointing to first content
then content, , first byte is size of record, then 9 size  each table
	1 size of the whole schema and table header info each
	


*/
const FILE_HEADER_SIZE = 100
// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]
	//	databaseFile, err := os.Open(databaseFilePath)
			
	databaseFile, err := os.Open(databaseFilePath)
	if err != nil {
				log.Fatal(err)
			}
	switch command {
		case ".tables":
			 //create table regex
		//	 var rgxarr_create_table = [13]int{67, 82, 69 , 65, 84, 69 ,32 ,84 ,65, 66, 76 ,69 ,32}
			//dot and parentheses regex
	//		var rgxdot int=10
	//		var rgxparent int =40
				offsetStartContentPointer:=5
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
				var offsetStartContent uint16
				if err := binary.Read(bytes.NewReader(schemaBuffer[offsetStartContentPointer:7]), binary.BigEndian, &offsetStartContent); err != nil {
					fmt.Println("Failed to read integer:", err)
					return
				}
				
				 fmt.Println(offsetStartContent);
				var contentSize uint16=pageSize-offsetStartContent
				contents:=make([]byte, contentSize)
				_, err = databaseFile.ReadAt(contents, int64(offsetStartContent))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(contents)
		
		case ".dbinfo":
		
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
					// You can use print statements as follows for debugging, they'll be visible when running tests.
			fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")
			// Uncomment this to pass the first stage
			fmt.Printf("database page size: %v", pageSize)
			
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

			fmt.Printf("number of tables: %v", tableCount)
			fmt.Println(schemaBuffer)
		
		default:
			fmt.Println("Unknown command", command)
			os.Exit(1)
	}
}
