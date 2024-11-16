package main
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"math"
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

func debugPrint(arguments ...interface{}) {
		debug:=0
		if debug==1{
		 fmt.Print("debug: ")
    for _, arg := range arguments {
        fmt.Print(arg," ")
    }
    fmt.Println()
		
		}
   
}
func uint16ArrayToString(array []uint16) string { 
byteSlice := make([]byte, len(array)*2) 
for i, value := range array {
 byteSlice[i*2] = byte(value) 
 byteSlice[i*2+1] = byte(value >> 8) } 
 return string(byteSlice)
 }
 func isHighBit(inte uint64)bool{
//   var b byte = byte(inte)
		
  // fmt.Println(b)
    result := inte ^ 128

    // Check if the result matches your condition
    if result < 128 {
        return true
        }
  	return false
 }
 func highBitParse(inte uint64) uint64{
 		return uint64(byte(inte) & 0b01111111 )
 }
 func byteGrabber(array []byte,startIndex uint64,nextIndex *uint32) uint64{
 		debugPrint("hai",*nextIndex)
 		var byteCount uint32=0
 		debugPrint(isHighBit(uint64(array[int(startIndex+uint64(*nextIndex+byteCount))])))
 		for isHighBit(uint64(array[int(startIndex+uint64(*nextIndex+byteCount))]))==true{
 		
	 		//*nextIndex+=1
	 		byteCount+=1
 		}
 	//	fmt.Println("hail",*nextIndex)
 	//	*nextIndex+=1
   	byteCount+=1
 		bufferHighBit:=make([]byte,byteCount)
 		debugPrint("byteCOunt",byteCount)
 		for i:=0;i<int(byteCount);i++{
 			debugPrint("ori byte",uint64(array[int(startIndex)+i+int(*nextIndex)]))
 			bufferHighBit[i]=byte(highBitParse(uint64(array[int(startIndex)+i+int(*nextIndex)])))
 		}
 		debugPrint("hitler",bufferHighBit,byteCount)
 			var size uint64
 			length:=int(byteCount)
 			for i:=1;i<=int(byteCount);i++{
	 		//	fmt.Println(uint64(array[int(startIndex)+i+int(*nextIndex)]))
	 			size+= uint64(int(bufferHighBit[length-i])*(int(math.Pow(2,float64(7*(i-1))))))
	 			debugPrint("size",size)
 			}
 			/*
				if err := binary.Read(bytes.NewReader(bufferHighBit[0:byteCount]), binary.BigEndian, &size); err != nil {
					fmt.Println("Failed to read integer:", err)
					return 0
				}
				*/
		*nextIndex+=byteCount
		debugPrint("padahal",*nextIndex)
 		return size
 		
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
				
				// fmt.Println(offsetStartContent);
				var contentSize uint16=pageSize-offsetStartContent
				contents:=make([]byte, contentSize)
				_, err = databaseFile.ReadAt(contents, int64(offsetStartContent))
				if err != nil {
					log.Fatal(err)
				}
			//	fmt.Println(contents)

				var tableCount uint16
			if err := binary.Read(bytes.NewReader(schemaBuffer[3:5]), binary.BigEndian, &tableCount); err != nil {
				fmt.Println("Failed to read integer:", err)
				return
			}
			var tableIndex uint64=0
			var nextIndex uint32=0
			for i:=0; i < int(tableCount); i++{
			
				var tableSize uint64=uint64(byteGrabber(contents,tableIndex,&nextIndex))
				debugPrint("haish",nextIndex)
				nextIndex+=1
				var gap_size_Rowid uint32=nextIndex
				var tableHeaderSize uint64=uint64(byteGrabber(contents,tableIndex,&nextIndex))
				var tableTypeSize uint64 = uint64(0.5 * float64(byteGrabber(contents,tableIndex,&nextIndex)) - 0.5*13)
				var tableNameSize uint64=uint64(0.5*float64(byteGrabber(contents,tableIndex,&nextIndex))-0.5*13)
				var tableTblNameSize uint64=uint64(0.5*float64(byteGrabber(contents,tableIndex,&nextIndex))-0.5*13)
				tableName := make([]uint16, tableTblNameSize)
				var indexTableName uint64=tableIndex+uint64(gap_size_Rowid)+tableHeaderSize+tableTypeSize +tableNameSize
				for ii:=0; ii < int(tableTblNameSize); ii++{
							tableName[ii]=uint16(contents[int(indexTableName)+ii])
							
							
							
				}
				debugPrint("=========",uint16ArrayToString(tableName))
				fmt.Println(uint16ArrayToString(tableName))
				tableIndex+=tableSize+uint64(gap_size_Rowid)
				nextIndex =0
				
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
