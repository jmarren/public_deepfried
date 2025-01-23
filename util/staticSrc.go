package util

import (
	"fmt"
	"os"
)

//	func GetJsSrc(filename string) string {
//		if os.Getenv("ENV") == "prod" {
//			return fmt.Sprintf("https://%s/public/js/%s", os.Getenv("STATIC_DOMAIN"), filename)
//		}
//		src := fmt.Sprintf("%sjs/%s", os.Getenv("STATIC_DOMAIN"), filename)
//		fmt.Printf("src: %s\n", src)
//		return src
//	}
func GetStaticSrc(filename string) string {
	if os.Getenv("ENV") == "prod" {
		return fmt.Sprintf("https://%s/public/%s", os.Getenv("STATIC_DOMAIN"), filename)
	}
	src := fmt.Sprintf("%s%s", os.Getenv("STATIC_DOMAIN"), filename)
	// fmt.Printf("src: %s\n", src)
	return src
}
