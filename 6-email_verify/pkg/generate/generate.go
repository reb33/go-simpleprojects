package generate

import (
	"encoding/hex"
	"github.com/google/uuid"
)

func Hash() string {
	id := uuid.New()
	hexStr := hex.EncodeToString(id[:])
	return hexStr
}
