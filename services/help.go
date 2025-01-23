package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/deepfried/consts"

	// "math"
	"os"
	"time"
)

func UuidString(userId pgtype.UUID) string {
	return uuid.UUID(userId.Bytes).String()
}

func getProfileSrc(userId pgtype.UUID) string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/profile_photos/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId))
	}
	return fmt.Sprintf("%sprofile_photos/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId))
}

func getArtSrc(userId pgtype.UUID, filename string) string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/artwork/%s/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId), filename)
	}
	return fmt.Sprintf("%sartwork/%s/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId), filename)
}

func getAudioSrc(userId pgtype.UUID, filename string) string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/audio/%s/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId), filename)
	}
	return fmt.Sprintf("%saudio/%s", os.Getenv("STATIC_DOMAIN"), filename)
}

func getStemsSrc(userId pgtype.UUID, mainFilename string, stemFilename string) string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/stems/%s/%s/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId), mainFilename, stemFilename)
	}
	return fmt.Sprintf("%sstems/%s/%s/%s", os.Getenv("STATIC_DOMAIN"), UuidString(userId), mainFilename, stemFilename)
}

// func getValidEltId(s string) string {
// 	invalids := []string{"'", " ", "!"}
// 	for _, invalid := range invalids {
// 		s = strings.ReplaceAll(s, invalid, "-")
// 	}
// 	return s
// }

func getTimeForUi(timestamp pgtype.Timestamp) string {
	elapsed := time.Since(timestamp.Time)
	totalMinutes := int(elapsed.Minutes())
	totalHours := int(elapsed.Hours())
	totalDays := int((totalHours / 24))

	if totalDays > 1 {
		return fmt.Sprintf("%d days ago", int(totalDays))
	} else if totalDays == 1 {
		return fmt.Sprintf("%d day ago", int(totalDays))
	} else if totalHours > 1 {
		return fmt.Sprintf("%d hours ago", int(totalHours))
	} else if totalHours == 1 {
		return fmt.Sprintf("1 hour ago")
	} else if totalMinutes > 1 {
		return fmt.Sprintf("%d minutes ago", int(totalMinutes))
	} else if totalMinutes == 1 {
		return fmt.Sprintf("1 minute ago")
	} else {
		return fmt.Sprintf("just now")
	}
}

// replaces any characters that are not valid for html ids with '_', as well as any hyphens ('-'),
// so that the resulting string can use a hyphen to separate the username from the title
func modifyInvalidChars(s string) string {
	validIdChars := [65]rune{}
	n := 0
	for i := 'A'; i <= 'Z'; i++ {
		validIdChars[n] = i
		n++
	}

	for i := 'a'; i <= 'z'; i++ {
		validIdChars[n] = i
		n++
	}

	nonLetterValids := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_', ':', '.'}

	for i := 0; i < len(nonLetterValids); i++ {
		validIdChars[n] = nonLetterValids[i]
		n++
	}

	sRunes := []rune(s)

	for i := 0; i < len(sRunes); i++ {
		valid := false
		for j := 0; j < len(validIdChars); j++ {
			if rune(validIdChars[j]) == rune(sRunes[i]) {
				valid = true
				break
			}
		}
		if !valid {
			sRunes[i] = '_'
		}
	}
	return string(sRunes)
}

func getDisplayMusicalKey(key consts.MusicalKey, keySig consts.MusicalKeySignature, majorMinor consts.MajorMinor) string {
	var sigStr string = ""

	if keySig == "flat" {
		sigStr = "&#9837;"
	} else if keySig == "sharp" {
		sigStr = "#"
	}

	isMinorStr := ""

	if majorMinor == "Minor" {
		isMinorStr = "m"
	}

	return fmt.Sprintf("%s%s%s", key, sigStr, isMinorStr)
}
