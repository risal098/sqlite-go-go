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
	1 size of the whole schema and table header info each (exclude rowid and itself)
	schema size is from "create table until ) parentheses"
	main header is offset 2 from whole info (after size of record and rowid)
	offset 5 in header is size of schema, so there is 5 fixed byte in header, and dynamic size in byte iinfo
	


*/
func uint16ArrayToString(array []uint16) string { 
byteSlice := make([]byte, len(array)*2) 
for i, value := range array {
 byteSlice[i*2] = byte(value) 
 byteSlice[i*2+1] = byte(value >> 8) } 
 return string(byteSlice)
 }
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
				
			//	 fmt.Println(offsetStartContent);
				var contentSize uint16=pageSize-offsetStartContent
				contents:=make([]byte, contentSize)
				_, err = databaseFile.ReadAt(contents, int64(offsetStartContent))
				if err != nil {
					log.Fatal(err)
				}
		//		fmt.Println(contents,len(contents))
				var tableCount uint16
			if err := binary.Read(bytes.NewReader(schemaBuffer[3:5]), binary.BigEndian, &tableCount); err != nil {
				fmt.Println("Failed to read integer:", err)
				return
			}
			var tableIndex uint64=0
			//allTableName := make([][]uint16, tableCount)
			for i:=0; i < int(tableCount); i++{
				var tableSize uint64=uint64(contents[tableIndex ])
			//	var tableId uint16=uint16(contents[tableIndex +1])
				var tableHeaderSize uint64=uint64(contents[tableIndex +2])
				var tableTypeSize uint64 = uint64(0.5 * float64(contents[tableIndex+3]) - 0.5*13)
				var tableNameSize uint64=uint64(0.5*float64(contents[tableIndex +4])-0.5*13)
				var tableTblNameSize uint64=uint64(0.5*float64(contents[tableIndex +5])-0.5*13)
				tableName := make([]uint16, tableTblNameSize)
				//fmt.Println("current",tableIndex,">",tableSize,tableHeaderSize ,tableIndex+tableSize+2)
			//	fmt.Println("tblnmsz",tableTblNameSize)
				var indexTableName uint64=tableIndex+2+tableHeaderSize+tableTypeSize +tableNameSize
				for ii:=0; ii < int(tableTblNameSize); ii++{
							tableName[ii]=uint16(contents[int(indexTableName)+ii])
							
							
							
				}
			//	fmt.Println("size",tableSize,"next pos",(contents[tableIndex +tableSize+2]),(contents[tableIndex +tableSize+3]),(contents[tableIndex +tableSize+4]))
	//			fmt.Println(tableIndex,tableSize,2,tableHeaderSize,tableTypeSize ,tableNameSize,"{",float64(contents[tableIndex +4]),"}",tableIndex +tableSize+3,"===<>")
		//		fmt.Print(indexTableName)
			//	fmt.Print("  ->")
				fmt.Println(uint16ArrayToString(tableName))
				tableIndex+=tableSize+2
				
			}
		
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
