package dvpl_logic

import (
	"errors"
	"hash/crc32"

	"github.com/pierrec/lz4/v4"
)

const (
	RedColor    = "\033[31m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	ResetColor  = "\033[0m"
)

const (
	dvplFooterSize = 20
	dvplTypeNone   = 0
	dvplTypeLZ4    = 2
	dvplExtension  = ".dvpl"
)

type DVPLFooter struct {
	OriginalSize   uint32
	CompressedSize uint32
	CRC32          uint32
	Type           uint32
}

// CompressDVPL compresses a buffer and returns the processed DVPL file buffer.
func CompressDVPL(buffer []byte) ([]byte, error) {
	compressedBlockSize := lz4.CompressBlockBound(len(buffer))
	compressedBlock := make([]byte, compressedBlockSize)

	n, err := lz4.CompressBlock(buffer, compressedBlock, nil)
	if err != nil {
		return nil, err
	}

	compressedBlock = compressedBlock[:n]
	footerBuffer := createDVPLFooter(uint32(len(buffer)), uint32(n), crc32.ChecksumIEEE(compressedBlock), 2)
	return append(compressedBlock, footerBuffer...), nil
}

// DecompressDVPL decompresses a DVPL buffer and returns the uncompressed file buffer.
func DecompressDVPL(buffer []byte) ([]byte, error) {
	footerData, err := readDVPLFooter(buffer)
	if err != nil {
		return nil, err
	}

	targetBlock := buffer[:len(buffer)-20]

	if uint32(len(targetBlock)) != footerData.CompressedSize {
		return nil, errors.New("DVPLSizeMismatch")
	}

	if crc32.ChecksumIEEE(targetBlock) != footerData.CRC32 {
		return nil, errors.New("DVPLCRC32Mismatch")
	}

	if footerData.Type == 0 {
		if !(footerData.OriginalSize == footerData.CompressedSize && footerData.Type == 0) {
			return nil, errors.New("DVPLTypeSizeMismatch")
		}
		return targetBlock, nil
	} else if footerData.Type == 1 || footerData.Type == 2 {
		deDVPLBlock := make([]byte, footerData.OriginalSize)
		n, err := lz4.UncompressBlock(targetBlock, deDVPLBlock)
		if err != nil {
			return nil, err
		}

		if uint32(n) != footerData.OriginalSize {
			return nil, errors.New("DVPLDecodeSizeMismatch")
		}

		return deDVPLBlock, nil
	}

	return nil, errors.New("UNKNOWN DVPL FORMAT")
}

// createDVPLFooter creates a DVPL footer from the provided data.
func createDVPLFooter(inputSize, compressedSize, crc32, typeVal uint32) []byte {
	result := make([]byte, 20)
	writeLittleEndianUint32(result, inputSize, 0)
	writeLittleEndianUint32(result, compressedSize, 4)
	writeLittleEndianUint32(result, crc32, 8)
	writeLittleEndianUint32(result, typeVal, 12)
	copy(result[16:], "DVPL")
	return result
}

func readLittleEndianUint32(b []byte, offset int) uint32 {
	return uint32(b[offset]) | uint32(b[offset+1])<<8 | uint32(b[offset+2])<<16 | uint32(b[offset+3])<<24
}

// readDVPLFooter reads the DVPL footer data from a DVPL buffer.
func readDVPLFooter(buffer []byte) (*DVPLFooter, error) {
	footerBuffer := buffer[len(buffer)-20:]
	if string(footerBuffer[16:]) != "DVPL" || len(footerBuffer) != 20 {
		return nil, errors.New(RedColor + "InvalidDVPLFooter" + ResetColor)
	}

	footerData := &DVPLFooter{}
	footerData.OriginalSize = readLittleEndianUint32(footerBuffer, 0)
	footerData.CompressedSize = readLittleEndianUint32(footerBuffer, 4)
	footerData.CRC32 = readLittleEndianUint32(footerBuffer, 8)
	footerData.Type = readLittleEndianUint32(footerBuffer, 12)
	return footerData, nil
}

func writeLittleEndianUint32(b []byte, v uint32, offset int) {
	b[offset+0] = byte(v)
	b[offset+1] = byte(v >> 8)
	b[offset+2] = byte(v >> 16)
	b[offset+3] = byte(v >> 24)
}
