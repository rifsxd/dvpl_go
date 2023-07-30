package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	lz4 "github.com/pierrec/lz4/v4"
)

func main() {
	fmt.Println()
	color.Cyan("######################################################################")
	color.Cyan("############# RXD DVPL CONVERTER GOLANG EDITION V2.1.0 ###############")
	color.Cyan("######################################################################")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("No mode selected. Try 'dvplgo --help or -h' for advice.")
		fmt.Println()
		return
	}

	realArgs := os.Args[1:]
	optionalArgs := realArgs[1:]

	keepOriginals := false

	for _, arg := range optionalArgs {
		if arg == "--keep-originals" || arg == "-ko" || arg == "--keep-original" {
			keepOriginals = true
			break
		}
	}

	processDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting the current working directory: %v", err)
	}

	switch strings.ToLower(realArgs[0]) {
	case "compress", "comp", "cp", "c":
		// compress
		color.Cyan("Compressing...")
		color.Cyan("              ")
		loadingDone := make(chan struct{})
		go loadingCircle(loadingDone)
		count, err := Recursion(filepath.Clean(processDir), keepOriginals, true)
		close(loadingDone)
		if err != nil {
			log.Fatalf("Compression failed: %v", err)
		}
		printSuccess(fmt.Sprintf("Compression completed. %s compressed.", formatFileCount(count)))
		fmt.Println()
	case "decompress", "decomp", "dcp", "d":
		// decompress
		color.Cyan("Decompressing...")
		color.Cyan("                ")
		loadingDone := make(chan struct{})
		go loadingCircle(loadingDone)
		count, err := Recursion(filepath.Clean(processDir), keepOriginals, false)
		close(loadingDone)
		if err != nil {
			log.Fatalf("Decompression failed: %v", err)
		}
		printSuccess(fmt.Sprintf("Decompression completed. %s decompressed.", formatFileCount(count)))
		fmt.Println()
	case "--help", "-h":
		color.Cyan(`dvplgo [mode] [--keep-originals]
	mode can be the following:
		compress (comp, cp, c): compresses files into dvpl
		decompress (decomp, dcp, d): decompresses dvpl files into standard files
		--help (-h): show this help message
		--keep-originals (--keep-original, -ko): flag keeps the original files after compression/ decompression`)
		fmt.Println()
	default:
		log.Fatalf("Incorrect mode selected. Use Help for information")
	}
}

// Recursion is the main code that does all the heavy lifting
func Recursion(originalsDir string, keepOrignals bool, compression bool) (int, error) {
	count := 0

	dirList, err := ioutil.ReadDir(originalsDir)
	if err != nil {
		return count, err
	}

	for _, dirItem := range dirList {
		if dirItem.IsDir() {
			subdirCount, err := Recursion(filepath.Join(originalsDir, dirItem.Name()), keepOrignals, compression)
			if err != nil {
				return count, err
			}
			count += subdirCount
		} else if (compression && !strings.HasSuffix(dirItem.Name(), ".dvpl")) ||
			(!compression && strings.HasSuffix(dirItem.Name(), ".dvpl")) {
			// Process files based on the compression flag
			count++

			filePath := filepath.Join(originalsDir, dirItem.Name())
			fileData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v\n", filePath, err)
				continue
			}

			var processedBlock []byte
			if compression {
				processedBlock = compressDVPL(fileData)
				filePath += ".dvpl"
			} else {
				processedBlock, err = decompressDVPL(fileData)
				if err != nil {
					log.Printf("Failed to decompress file %s: %v\n", filePath, err)
					continue
				}
				filePath = strings.TrimSuffix(filePath, ".dvpl")
			}

			err = ioutil.WriteFile(filePath, processedBlock, 0644)
			if err != nil {
				log.Printf("Failed to write file %s: %v\n", filePath, err)
				continue
			}

			if !keepOrignals {
				err = os.Remove(filepath.Join(originalsDir, dirItem.Name()))
				if err != nil {
					log.Printf("Failed to remove file %s: %v\n", filePath, err)
				}
			}
		}
	}

	return count, nil
}

// CompressDVPL is equivalent to the compressDVPL JavaScript function
func compressDVPL(buffer []byte) []byte {
	// Try compressing the data
	compressedData := make([]byte, lz4.CompressBlockBound(len(buffer)))
	compressedSize, err := lz4.CompressBlock(buffer, compressedData, nil)
	if err != nil || compressedSize >= len(buffer) {
		// If compression fails or data becomes bigger, store the uncompressed data
		footerBuffer := toDVPLFooter(len(buffer), len(buffer), crc32.ChecksumIEEE(buffer), 0)
		return append(buffer, footerBuffer...)
	}

	// Compression successful
	compressedData = compressedData[:compressedSize]
	footerBuffer := toDVPLFooter(len(buffer), len(compressedData), crc32.ChecksumIEEE(compressedData), 2)
	return append(compressedData, footerBuffer...)
}

// DecompressDVPL is equivalent to the decompressDVPL JavaScript function
func decompressDVPL(buffer []byte) ([]byte, error) {
	footerData, err := readDVPLFooter(buffer)
	if err != nil {
		return nil, err
	}

	targetBlock := buffer[:len(buffer)-20]

	if len(targetBlock) != int(footerData.cSize) {
		return nil, fmt.Errorf("DVPLSizeMismatch")
	}

	if crc32.ChecksumIEEE(targetBlock) != footerData.crc32 {
		return nil, fmt.Errorf("DVPLCRC32Mismatch")
	}

	if footerData.typ == 0 {
		// Data is uncompressed
		if !(footerData.oSize == footerData.cSize && footerData.typ == 0) {
			return nil, fmt.Errorf("DVPLTypeSizeMismatch")
		}
		return targetBlock, nil
	} else if footerData.typ == 1 || footerData.typ == 2 {
		// Data is compressed, decompress it
		deDVPLBlock := make([]byte, footerData.oSize)
		_, err := lz4.UncompressBlock(targetBlock, deDVPLBlock)
		if err != nil {
			return nil, err
		}
		return deDVPLBlock, nil
	}

	return nil, fmt.Errorf("UNKNOWN DVPL FORMAT")
}

// ToDVPLFooter is equivalent to the toDVPLFooter JavaScript function
func toDVPLFooter(inputSize, compressedSize int, crc32 uint32, typ int) []byte {
	result := make([]byte, 20)
	binary.LittleEndian.PutUint32(result[0:4], uint32(inputSize))
	binary.LittleEndian.PutUint32(result[4:8], uint32(compressedSize))
	binary.LittleEndian.PutUint32(result[8:12], crc32)
	binary.LittleEndian.PutUint32(result[12:16], uint32(typ))
	copy(result[16:], []byte("DVPL"))
	return result
}

// ReadDVPLFooter is equivalent to the readDVPLFooter JavaScript function
func readDVPLFooter(buffer []byte) (dvplFooter, error) {
	footerBuffer := buffer[len(buffer)-20:]
	if string(footerBuffer[16:]) != "DVPL" || len(footerBuffer) != 20 {
		return dvplFooter{}, fmt.Errorf("InvalidDVPLFooter")
	}

	footerObject := dvplFooter{}
	footerObject.oSize = int(binary.LittleEndian.Uint32(footerBuffer[0:4]))
	footerObject.cSize = int(binary.LittleEndian.Uint32(footerBuffer[4:8]))
	footerObject.crc32 = binary.LittleEndian.Uint32(footerBuffer[8:12])
	footerObject.typ = int(binary.LittleEndian.Uint32(footerBuffer[12:16]))
	return footerObject, nil
}

type dvplFooter struct {
	oSize int
	cSize int
	crc32 uint32
	typ   int
}

func printFileCount(count int, action string) {
	if count == 1 {
		fmt.Printf("1 file %s.\n", action)
	} else {
		fmt.Printf("%d files %s.\n", count, action)
	}
}

func formatFileCount(count int) string {
	if count == 1 {
		return "1 file"
	}
	return fmt.Sprintf("%d files", count)
}

// Add this function to print messages in green color
func printSuccess(message string) {
	color.Green("[✔] " + message)
}

// Add this function to print error messages in red color
func printError(message string) {
	color.Red("[✗] " + message)
}

func loadingCircle(done <-chan struct{}) {
	c := color.New(color.FgCyan)
	defer fmt.Println() // New line after loading animation
	for {
		select {
		case <-done:
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r")
				c.Printf("%c ", r) // Add a space after the loading circle character
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
