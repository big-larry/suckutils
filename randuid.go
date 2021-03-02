package suckutils

import (
	"encoding/binary"
	"math/rand"
	"time"
)

var real_rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandUID(salt uint32) []byte {
	v := real_rnd.Uint32()
	r := rand.New(rand.NewSource(int64(v + salt)))
	buf := make([]byte, 8)
	r.Read(buf[4:])
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func CheckRandUID(buf []byte, salt uint32) bool {
	r := rand.New(rand.NewSource(int64(binary.BigEndian.Uint32(buf) + salt)))
	buf1 := make([]byte, 4)
	r.Read(buf1)
	return buf[4] == buf1[0] && buf[5] == buf1[1] && buf[6] == buf1[2] && buf[7] == buf1[3]
}
