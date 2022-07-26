package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

/* panic: 247fc4278353510b <nil>
func init() {
	resp, err := AuthChallenge("582623C0C59E162B", "154051653")
	msg := fmt.Sprintf("%x %v", resp, err)
	panic(msg)
}
*/

var (
	XorKeyAuth = [8]byte{0xFF, 0xCC, 0x73, 0xFE, 0xEC, 0xF3, 0x57, 0xA8}
)

func AuthChallenge(fsid, challenge string) (string, error) {
	if len(fsid) != 16 {
		return "", fmt.Errorf("invalid fsid length")
	}

	iFSID, err := hex.DecodeString(fsid)
	if err != nil {
		return "", err
	}
	iFSID = iFSID[len(iFSID)-6:]
	for i := 0; i < len(iFSID); i++ {
		iFSID[i] ^= XorKeyAuth[i%4]
	}

	iChallenge := []byte(challenge)
	for i := 0; i < 8; i++ {
		iChallenge[i] ^= XorKeyAuth[i]
	}

	token := []byte{
		iFSID[0],
		iChallenge[0],
		iFSID[1],
		iChallenge[1],
		iFSID[2],
		iChallenge[2],
		iFSID[3],
		iChallenge[3],
		iFSID[4],
		0x55,
		0x67,
		iFSID[5],
		iChallenge[4],
		iChallenge[5],
		iChallenge[6],
		iChallenge[7],
	}

	tokenHash := md5.Sum(token)
	octets := make([]byte, len(tokenHash)/2)
	for i := 0; i < len(octets); i++ {
		octets[i] = tokenHash[i*2]
	}

	return fmt.Sprintf("%x", octets), nil
}