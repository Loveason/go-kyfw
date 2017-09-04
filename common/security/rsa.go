package security

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func RSAEncrypt(pwd string, card string, modulus string, exponent string) (encryptStr string, err error) {
	var msg []byte
	msg, err = formatToHex(pwd, card)
	bigint := new(big.Int)
	bigint.SetString(modulus, 16)
	publicKey := new(rsa.PublicKey)
	var E int64
	E, err = strconv.ParseInt(exponent, 16, 32)
	if err != nil {
		return
	}
	publicKey.E = int(E)
	publicKey.N = bigint
	var encryptBytes []byte
	encryptBytes, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
	if err != nil {
		return
	}
	encryptStr = bytesToHexString(encryptBytes)
	return
}

func formatToHex(pwd string, card string) (byteArray []byte, err error) {
	var c [8]int64
	var d [8]int64
	c[0] = 0x06
	d[0] = 0x00
	d[1] = 0x00
	e := 1
	h := 2

	for g := 0; g < len(pwd); g += 2 {
		c[e], err = strconv.ParseInt(string(pwd[g])+string(pwd[g+1]), 16, 32)
		if err != nil {
			return
		}
		e++
	}
	c[4] = 0xff
	c[5] = 0xff
	c[6] = 0xff
	c[7] = 0xff
	for g := len(card) - 13; g < len(card)-1; g += 2 {
		d[h], err = strconv.ParseInt(string(card[g])+string(card[g+1]), 16, 32)
		if err != nil {
			return
		}
		h++
	}

	for g := 0; g < 8; g++ {
		t := byte(d[g] ^ c[g])
		byteArray = append(byteArray, t)
	}
	return
}

func bytesToHexString(bts []byte) string {
	var sa = make([]string, 0)
	for _, v := range bts {
		sa = append(sa, fmt.Sprintf("%02x", v))
	}
	ss := strings.Join(sa, "")
	return ss
}
