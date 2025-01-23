package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jmarren/katana-mp3"
	"io"
	"math"
	"net/http"

	"slices"
	"strconv"
)

func readInt(data []byte) (uint64, error) {
	ioreader := bytes.NewReader(data)
	res, err := binary.ReadUvarint(io.ByteReader(ioreader))
	if err != nil {
		return 0, err
	}
	return res, nil
}

func readDataPoint(data []byte, blockAlign int, numChannels int, bytesPerSample int) {
	chanLen := blockAlign / numChannels
	left := data[0 : bytesPerSample/numChannels]
	right := data[bytesPerSample:]
	fmt.Printf("chanLen: %d\nleft: %b\nright: %b\n", chanLen, left, right)
}

func ParseWav(data *bytes.Buffer) []int32 {
	fmt.Println("*** buffer info ****")
	dataLength := data.Len()
	available := data.Available()
	usedBytes := dataLength - available

	fmt.Printf("buffer length: %d\navailable: %d\nUsedBytes: %d\n", dataLength, available, usedBytes)

	fmt.Println("------- File Parse --------")
	chunkId := string(data.Next(4))
	fmt.Printf("\tchunkId: %s\n", chunkId)

	chunkOneSize := binary.LittleEndian.Uint32(data.Next(4))
	fmt.Printf("\tchunkOneSize: %d\n", chunkOneSize)

	format := data.Next(4)
	fmt.Printf("\tformat: %s\n", format)

	subChunkOneId := data.Next(4)
	fmt.Printf("\tsubChunkOneId: %s\n", subChunkOneId)

	subChunkOneSize := binary.LittleEndian.Uint32(data.Next(4))
	fmt.Printf("\tsubChunkOneSize: %d\n", subChunkOneSize)

	for string(subChunkOneId) != "fmt " {
		data.Next(int(subChunkOneSize))
		subChunkOneId = data.Next(4)
		fmt.Printf("\tsubChunkOneId: %s\n", subChunkOneId)
		subChunkOneSize := binary.LittleEndian.Uint32(data.Next(4))
		fmt.Printf("\tsubChunkOneSize: %d\n", subChunkOneSize)
	}

	audioFormat := binary.LittleEndian.Uint16(data.Next(2))
	fmt.Printf("\taudioFormat: %d\n", audioFormat)

	numChannels := binary.LittleEndian.Uint16(data.Next(2))
	fmt.Printf("\tnumChannels: %d\n", numChannels)

	rawSampleRate := data.Next(4)
	sampleRate := binary.LittleEndian.Uint32(rawSampleRate)
	fmt.Printf("\tsampleRate: %d\n", int(sampleRate))

	byteRate := binary.LittleEndian.Uint32(data.Next(4)) // bytes per second
	fmt.Printf("\tbyteRate: %d\n", byteRate)

	blockAlign := binary.LittleEndian.Uint16(data.Next(2))
	fmt.Printf("\tblockAlign: %d\n", blockAlign)

	bitsPerSample := binary.LittleEndian.Uint16(data.Next(2))
	fmt.Printf("\tbitsPerSample: %d\n", bitsPerSample)

	subChunkTwoId := data.Next(4)
	fmt.Printf("\tsubChunkTwoId: %s\n", subChunkTwoId)

	subChunkTwoSize := binary.LittleEndian.Uint32(data.Next(4))
	fmt.Printf("\tsubChunkTwoSize: %d\n", subChunkTwoSize)

	for string(subChunkTwoId) != "data" {
		data.Next(int(subChunkTwoSize))
		subChunkTwoId = data.Next(4)
		fmt.Printf("\tsubChunkTwoId: %s\n", subChunkTwoId)
		subChunkTwoSize := binary.LittleEndian.Uint32(data.Next(4))
		fmt.Printf("\tsubChunkTwoSize: %d\n", subChunkTwoSize)
	}

	fmt.Println("-------------------------")

	// length of sampled data chunk = bytesPerSample * numChannels * total blocks
	bytesPerSample := bitsPerSample / 8
	samplesInChunk := int(subChunkTwoSize) / int(bytesPerSample)
	blockLength := int(bytesPerSample * blockAlign * numChannels)
	chanLen := int(bytesPerSample * blockAlign)
	numBlocks := int(samplesInChunk) / int(blockAlign*numChannels)

	fmt.Printf("bytesPerSample: %d bytes\nsamplesInChunk: %d samples\nblockLength: %d\nchanLen: %d bytes\nnumBlocks: %d blocks\n", bytesPerSample, samplesInChunk, blockLength, chanLen, numBlocks)

	batchSize := numBlocks / 200

	totals := []int32{}

	m := 0
	total := 0

	for i := 0; i < numBlocks; i++ {
		leftChan := data.Next(chanLen)
		rightChan := data.Next(chanLen)
		for j := 0; j < int(blockAlign); j++ {
			start := int(bytesPerSample) * j
			end := int(bytesPerSample) * (j + 1)
			l, err := GetSignedInt(leftChan[start:end])
			r, err := GetSignedInt(rightChan[start:end])
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
			avg := (l + r) / 2
			total += avg
		}
		m++
		if m >= batchSize {
			totals = append(totals, int32(math.Abs(float64(total/(batchSize)))))
			total = 0
			m = 0
		}
	}

	fmt.Printf("totals: %v\nlen(totals): %d\n", totals, len(totals))
	return totals
}

func ParseAudioArr(data *bytes.Buffer) ([]int32, error) {
	dataBytes := data.Bytes()
	fileType := http.DetectContentType(dataBytes)
	if fileType == "audio/wave" {
		return ParseWav(data), nil
	}
	if fileType == "audio/mpeg" {
		return DecodeMp3(data)
	}

	return []int32{}, fmt.Errorf("expected fileType audio/wave or audio/mpeg, got: %s\n", fileType)

}

func ReturnAsInt(data []byte) (int, error) {
	switch len(data) {
	case 1:
		return int(binary.LittleEndian.Uint16(data)), nil
	case 2:
		return int(binary.LittleEndian.Uint16(data)), nil
	case 3:
		newData := []byte{0, data[0], data[1], data[2]}
		return int(binary.LittleEndian.Uint32(newData)), nil
	case 4:
		return int(binary.LittleEndian.Uint32(data)), nil
	default:
		return 0, fmt.Errorf("error: expected bits per sample of 32 or less\n")
	}
}

func DecodeMp3(data *bytes.Buffer) ([]int32, error) {
	var visArr []int32
	dataBytes := data.Bytes()
	reader := bytes.NewReader(dataBytes)
	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		return visArr, err
	}
	for err != io.EOF {
		err = decoder.ReadFrame()
	}
	visArr = splitAndAvg(decoder.Buf)
	return visArr, nil
}

func splitAndAvg(data []byte) []int32 {
	chunkSize := len(data) / 200

	avgs := []int32{}
	for c := range slices.Chunk(data, chunkSize) {
		var total int
		for i := 0; i < len(c)-2; i += 2 {
			val, err := GetSignedInt(c[i : i+2])
			if err != nil {
				fmt.Printf("error:  %s\n", err)
			}
			total += val
		}
		avg := int32(total / len(c))
		avgs = append(avgs, avg)
	}

	return avgs
}

func getVersionId(audioVersionIdBin int) {
	switch audioVersionIdBin {
	case 0:
		fmt.Println("MPEG Version 2.5")
	case 1:
		fmt.Println("reserved?")
	case 2:
		fmt.Println("MPEG Version 2")
	case 3:
		fmt.Println("MPEG Version 1")
	}
}

func getLayerIndex(layerIndexBin int) {
	switch layerIndexBin {
	case 0:
		fmt.Println("reserved? ")
	case 1:
		fmt.Println("layer III (3)")
	case 2:
		fmt.Println("layer II (2)")
	case 3:
		fmt.Println("layer I (1)")
	}
}

func getBitrate(layerIndex int, bitrateIndex int) (int, error) { /// in kilobits per second
	fmt.Printf("layerIndex: %d\nbitrateIndex: %d\n", layerIndex, bitrateIndex)
	// layer I
	if layerIndex == 3 {
		return bitrateIndex * 32, nil
	} else if layerIndex == 2 { // layer II
		if bitrateIndex < 9 {
			return 32 + ((bitrateIndex - 1) * 16), nil
		} else if bitrateIndex < 13 {
			return 160 + ((bitrateIndex - 9) * 32), nil
		} else if bitrateIndex < 15 {
			return 320 + ((bitrateIndex - 13) * 64), nil
		} else {
			return 0, fmt.Errorf("bitrate: reserved?")
		}
	} else if layerIndex == 1 { // layer III
		if bitrateIndex < 6 {
			return 32 + ((bitrateIndex - 1) * 8), nil
		} else if bitrateIndex < 9 {
			return 80 + ((bitrateIndex - 6) * 16), nil
		} else if bitrateIndex < 12 {
			return 128 + ((bitrateIndex - 9) * 32), nil
		} else if bitrateIndex < 15 {
			return 224 + ((bitrateIndex - 12) * 64), nil
		} else {
			return 0, fmt.Errorf("bitrate: reserved?")
		}
	}

	return 0, fmt.Errorf("expected layerIndex greater than 0 and less than 4, got: %d\n", layerIndex)

}

func getSampleRate(audioVersionId int, sampleRateIndex int) (int, error) {

	if sampleRateIndex > 2 {
		return 0, fmt.Errorf("expected sampleRateIndex < 3, got %d\n", sampleRateIndex)
	}

	// Mpeg v1
	if audioVersionId == 3 {
		switch sampleRateIndex {
		case 0:
			return 44100, nil
		case 1:
			return 48000, nil
		case 2:
			return 32000, nil
		}
	}
	// Mpeg v2
	if audioVersionId == 2 {
		switch sampleRateIndex {
		case 0:
			return 22050, nil
		case 1:
			return 24000, nil
		case 2:
			return 16000, nil
		}
	}
	// mpeg v2.5
	if audioVersionId == 0 {
		switch sampleRateIndex {
		case 0:
			return 11025, nil
		case 1:
			return 12000, nil
		case 2:
			return 8000, nil
		}
	}

	return 0, fmt.Errorf("expected audioVersionId of 0, 2, or 3, got %d\n", audioVersionId)
}

func printChannelMode(channelMode int) {
	switch channelMode {
	case 0:
		fmt.Println("stereo")
	case 1:
		fmt.Println("joint stereo")
	case 2:
		fmt.Println("dual channel (2 mono channels)")
	case 3:
		fmt.Println("single channel (mono)")
	}
}

func getSamplesPerFrame(layerIndex int, audioVersionId int) (int, error) {
	switch audioVersionId {
	case 0: // mpeg v2.5
		switch layerIndex {
		case 1: // layer 3
			return 576, nil
		case 2: // layer 2
			return 1152, nil
		case 3: // layer 1
			return 384, nil
		}
	case 2: // mpeg 2
		switch layerIndex {
		case 1: // layer 3
			return 576, nil
		case 2: // layer 2
			return 1152, nil
		case 3: // layer 1
			return 384, nil
		}
	case 3: // mpeg v1
		switch layerIndex {
		case 1: // layer 3
			return 1152, nil
		case 2: // layer 2
			return 1152, nil
		case 3: // layer 1
			return 384, nil
		}
	}

	return 0, fmt.Errorf("expected audioVersionId of 0, 2, or 3 but got %d\n", audioVersionId)

}

// [11111010]

type firstByteData struct {
	audioVersionId int
	layerIndex     int
	protectedBit   int
}

func HandleFirstByteAfterSync(data byte) (*firstByteData, error) {
	withoutSyncBits := (data << 3) >> 3
	fmt.Printf("data: %b\nwithoutSyncBits: %b\n", data, withoutSyncBits)

	audioVersionIdCmp, err := strconv.ParseUint("00011000", 2, 8)
	if err != nil {
		return &firstByteData{}, err
	}

	audioVersionId := (data & byte(audioVersionIdCmp)) >> 3
	fmt.Printf("audioVersionId: %b\n", audioVersionId)
	getVersionId(int(audioVersionId))

	layerIndexCmp, err := strconv.ParseUint("00000110", 2, 8)
	if err != nil {
		return &firstByteData{}, err
	}

	layerIndex := (data & byte(layerIndexCmp)) >> 1
	fmt.Printf("layerIndex: %b\n", layerIndex)

	protectedBit := (data << 7) >> 7
	fmt.Printf("protectedBit: %b\n", protectedBit)

	return &firstByteData{
		audioVersionId: int(audioVersionId),
		layerIndex:     int(layerIndex),
		protectedBit:   int(protectedBit),
	}, nil
}

type secondByteData struct {
	bitrateIndex    int
	samplerateIndex int
	paddingBit      int
}

func HandleSecondByteAfterSync(data byte) (*secondByteData, error) {
	fmt.Printf("byte: %b\n", data)

	bitrateIndex := data >> 4
	fmt.Printf("bitrateIndex: %b\n", bitrateIndex)

	sampleRateIndexCmp, err := strconv.ParseUint("00001100", 2, 8)
	if err != nil {
		return &secondByteData{}, nil
	}

	sampleRateIndex := (data & byte(sampleRateIndexCmp)) >> 2
	fmt.Printf("sampleRateIndex: %b\n", sampleRateIndex)

	paddingBitCmp, err := strconv.ParseUint("00000010", 2, 8)
	if err != nil {
		return &secondByteData{}, nil
	}
	paddingBit := (data & byte(paddingBitCmp)) >> 1
	fmt.Printf("paddingBit: %b\n", paddingBit)

	// final bit is "private bit" (not sure what its purpose is)
	return &secondByteData{
		bitrateIndex:    int(bitrateIndex),
		samplerateIndex: int(sampleRateIndex),
		paddingBit:      int(paddingBit),
	}, nil

}

type thirdByteData struct {
	channelMode   int
	modeExtension int
	copyrightBit  int
	originalBit   int
	emphasis      int
}

func HandleFirstTwoSideInfoBytes(data []byte, isDual bool) int {
	firstByte := data[0]
	fmt.Printf("firstByte: %b\n", firstByte)

	secondByte := data[1]
	fmt.Printf("secondByte: %b\n", secondByte)

	firstByte = firstByte << 1
	extraBit := secondByte >> 7
	fmt.Printf("%b\n", extraBit)

	// emphasisCmp, err := strconv.ParseUint("00000011", 2, 8)

	privateBitMask, err := strconv.ParseUint("01111100", 2, 8)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	if isDual {
		privateBitMask, err = strconv.ParseUint("01110000", 2, 8)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}

	privateBits := secondByte & byte(privateBitMask)
	fmt.Printf("privateBits: %b\n", privateBits)

	return int(firstByte + extraBit)

}

func getShareValue(data []byte, isDual bool) int {
	if isDual {
		// share is bits 13 - 21 of side info
		// -> meaning bits 5 - 13 of data being passed
		// - bits 5 to 8 of data[0] and bits 1 - 5 of data[1]
		fiveToEightMask, err := strconv.ParseUint("0000111", 2, 8)
		if err != nil {
			fmt.Printf("error:  %s\n", err)
		}
		fiveToEight := data[0] & byte(fiveToEightMask)
		extraBits := data[1] >> 3
		share := (fiveToEight << 4) + extraBits // overflow?
		return int(share)
	}

	return 1
	// share is bits 7 to 11 of data passed
	// firstTwoMask, err := strconv.ParseUint("00000011", 2, 8)
	// if err != nil {
	// 	fmt.Printf("err: %s\n", err)
	// }
	// firstTwo := data[0] & byte(firstTwoMask)
	// next

}

func DecodeSideInformation(data []byte) {
	isDual := false
	if len(data) == 32 {
		isDual = true
	}
	fmt.Printf("isDual: %t\n", isDual)

	mainDataBegin := HandleFirstTwoSideInfoBytes(data[0:2], isDual)

	fmt.Printf("mainDataBegin: %d\n", int(mainDataBegin))

	// shareValue := getShareValue(data[1:3], isDual)

}

// [01 10 1 1 00]
func HandleThirdByteAfterSync(data byte) (*thirdByteData, error) {
	fmt.Printf("third byte: %b\n", data)

	channelMode := data >> 6
	fmt.Printf("channelMode: %b\n", channelMode)

	modeExtensionCmp, err := strconv.ParseUint("00110000", 2, 8)
	if err != nil {
		return &thirdByteData{}, err
	}
	modeExtension := (data & byte(modeExtensionCmp)) >> 4
	fmt.Printf("modeExtension: %b\n", modeExtension)

	copyrightBitCmp, err := strconv.ParseUint("00001000", 2, 8)
	if err != nil {
		return &thirdByteData{}, err
	}
	copyrightBit := (data & byte(copyrightBitCmp)) >> 3
	fmt.Printf("copyrightBit: %b\n", copyrightBit)

	originalBitCmp, err := strconv.ParseUint("00000100", 2, 8)
	if err != nil {
		return &thirdByteData{}, err
	}
	originalBit := (data & byte(originalBitCmp)) >> 2
	fmt.Printf("originalBit: %b\n", originalBit)

	emphasisCmp, err := strconv.ParseUint("00000011", 2, 8)
	if err != nil {
		return &thirdByteData{}, err
	}
	emphasis := (data & byte(emphasisCmp))
	fmt.Printf("emphasis (rarely used): %b\n", emphasis)

	return &thirdByteData{
		channelMode:   int(channelMode),
		modeExtension: int(modeExtension),
		copyrightBit:  int(copyrightBit),
		originalBit:   int(originalBit),
		emphasis:      int(emphasis),
	}, nil

}

func FindSyncWord(dataBytes []byte) (int, error) {
	// dataBytes := data.Bytes()
	allOneByte := "11111111"
	allOneBin, err := strconv.ParseUint(allOneByte, 2, 8)
	if err != nil {
		return 0, fmt.Errorf("err parsing allOneByte: %s\n", err)
	}
	threeStart := "11100000"
	threeStartBin, err := strconv.ParseUint(threeStart, 2, 8)
	if err != nil {
		fmt.Printf("err parsing threeStart: %s\n", err)
	}

	j := 0

	for i := 0; i < len(dataBytes); i++ {
		if dataBytes[i] == byte(allOneBin) {
			if dataBytes[i+1]&byte(threeStartBin) == byte(threeStartBin) {
				fmt.Printf("dataBytes[i]:  %b\n", dataBytes[i])
				fmt.Printf("dataBytes[i + 1]: %b\n", dataBytes[i+1])
				return i, nil
			}
		}
	}

	return j - 1, nil
}

func returnAsInt(data []byte) (int, error) {
	switch len(data) {
	case 1:
		return int(binary.LittleEndian.Uint16(data)), nil
	case 2:
		return int(binary.LittleEndian.Uint16(data)), nil
	case 3:
		newData := []byte{0, data[0], data[1], data[2]}
		return int(binary.LittleEndian.Uint32(newData)), nil
	case 4:
		return int(binary.LittleEndian.Uint32(data)), nil
	default:
		return 0, fmt.Errorf("error: expected bits per sample of 32 or less\n")
	}
}

func GetSignedInt(data []byte) (int, error) {
	leftMost := (data[0]) >> 7
	if leftMost == 0 {
		return returnAsInt(data)
	}
	for i := 0; i < len(data); i++ {
		data[i] = ^(data[i])
	}

	oneLeft := 0
	i := 1
	for oneLeft == 255 {
		data[len(data)-i] += 1
		i++
		oneLeft = int(data[len(data)-i])
	}

	val, err := returnAsInt(data)

	return -1 * val, err
}

func AnalyzeAudioBytes(data []byte) int32 {
	// fmt.Printf("dataLength: %d\n", len(data))
	total := 0
	for i := 0; i < (len(data) / 2); i += 2 {
		val, err := GetSignedInt(data[i : i+2])
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		total += val
	}
	dataLen := len(data) / 2
	return int32(total / dataLen)
}

func ParseHeaders(dataBytes []byte) (int, int, []byte, error) {
	startIndex, err := FindSyncWord(dataBytes)
	fmt.Printf("startIndex: %d\n", startIndex)
	fmt.Println("------- File Parse --------")
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error finding sync word: %s\n", err)
	}

	firstByteData, err := HandleFirstByteAfterSync(dataBytes[startIndex+1])
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error getting first byte data: %s\n", err)
	}

	secondByteData, err := HandleSecondByteAfterSync(dataBytes[startIndex+2])
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error getting second byte data: %s\n", err)
	}
	thirdByteData, err := HandleThirdByteAfterSync(dataBytes[startIndex+3])
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error getting third byte data: %s\n", err)
	}

	bitrate, err := getBitrate(int(firstByteData.layerIndex), int(secondByteData.bitrateIndex))
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error calculating bitrate: %s\n", err)
	}
	fmt.Printf("bitrate: %d\n", bitrate)

	sampleRate, err := getSampleRate(firstByteData.audioVersionId, secondByteData.samplerateIndex)
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error calculating sampleRate: %s\n", err)
	}
	fmt.Printf("sampleRate: %d\n", sampleRate)

	var slotSize = 1 // in bytes
	if firstByteData.layerIndex == 3 {
		slotSize = 4
	}

	paddingSize := 0
	if secondByteData.paddingBit == 1 {
		paddingSize = 1 * slotSize
	}

	printChannelMode(thirdByteData.channelMode)

	samplesPerFrame, err := getSamplesPerFrame(firstByteData.layerIndex, firstByteData.audioVersionId)
	if err != nil {
		return 0, 0, []byte{}, fmt.Errorf("error calculating samples per frame: %s\n", err)
	}
	fmt.Printf("samples per frame: %d\n", samplesPerFrame)

	byteRate := bitrate / 8

	frameLength := int((samplesPerFrame*byteRate*1000)/sampleRate) + paddingSize
	fmt.Printf("frameLength: %d bytes per frame\n", frameLength)

	fmt.Printf("bitrate: %d\n", bitrate)
	fmt.Printf("sampleRate: %d\n", sampleRate)
	fmt.Printf("paddingSize: %d\n", paddingSize)

	frameSize := ((samplesPerFrame / 8 * bitrate) / (sampleRate / 8)) + paddingSize
	fmt.Printf("frameSize: %d\n", frameSize)

	nextHeaderStart := startIndex + frameLength
	rawAudioBytes := dataBytes[startIndex+4 : nextHeaderStart]

	sideInfoLength := 17
	if thirdByteData.channelMode == 3 {
		sideInfoLength = 32
	}
	sideInfoStart := startIndex + 4
	sideInfoEnd := sideInfoStart + sideInfoLength
	sideInfo := dataBytes[sideInfoStart:sideInfoEnd]
	DecodeSideInformation(sideInfo)
	// AnalyzeAudioBytes(rawAudioBytes)

	// frameDuration := (nextHeaderStart - startIndex) / (bitrate * 8)
	// nextHeaderEnd := nextHeaderStart + 4
	return nextHeaderStart, bitrate, rawAudioBytes, nil
}

func ParseMpeg(data *bytes.Buffer) ([]int32, error) {
	// fileSize := data.Len() - data.Available()

	dataBytes := data.Bytes()
	fileSizeInBits := len(dataBytes) * 8

	frameIndexes := []int{}
	startIndex := 0
	numFrames := 0
	bitrate := 0
	// currentSize := 0
	// dataSection := []byte{}
	avgs := []int32{}

	rawData := []byte{}
	for startIndex < len(dataBytes) {
		nextHeaderStart, frameBitrate, rawAudio, err := ParseHeaders(dataBytes)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		bitrate = frameBitrate
		rawData = append(rawData, rawAudio...)
		startIndex += nextHeaderStart
		frameIndexes = append(frameIndexes, startIndex)
		dataBytes = dataBytes[nextHeaderStart:]
		numFrames++
	}

	totalRawSize := len(rawData)
	batchSize := totalRawSize / 100

	rawLen := 0

	for c := range slices.Chunk(rawData, batchSize) {
		val := AnalyzeAudioBytes(c)
		avgs = append(avgs, val)
	}

	fmt.Printf("rawLen: %d\n", rawLen)

	// batchSize := fileSizeInBits / 100

	// fmt.Printf("avgs: %v\n", avgs)
	fmt.Printf("len(avgs): %d\n", len(avgs))
	fmt.Printf("numFrames: %d\n", numFrames)
	// fmt.Printf("frameIndexes: %v\n", frameIndexes)
	fmt.Printf("bitrate: %d\n", bitrate)
	durationInSeconds := (fileSizeInBits / bitrate) / 1000
	fmt.Printf("durationInSeconds: %d\n", durationInSeconds)

	return avgs, nil
}
