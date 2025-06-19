package req

import (
	"github.com/golang-jwt/jwt"
)

// Custom claims structure
type CustomClaimsModel struct {
	SeatID int64 `json:"uid" comment:"座席ID"`
	//AuthorityId string
	//BufferTime  int64
	jwt.StandardClaims
}
