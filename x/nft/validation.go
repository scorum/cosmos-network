// Package nft is created with purpose to change nft id validation (allow a leading digit).
package nft

import (
	"fmt"
	"regexp"
	_ "unsafe"
)

//go:linkname reNFTID github.com/cosmos/cosmos-sdk/x/nft.reNFTID
var reNFTID *regexp.Regexp

func init() {
	reNFTID = regexp.MustCompile(fmt.Sprintf(`^%s$`, `[a-zA-Z0-9][a-zA-Z0-9/:-]{2,100}`))
}
