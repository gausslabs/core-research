package main

import (
	"app/src"
	"fmt"
	"time"
)

func main() {

	var witness uint64 = 0xffffffffffffffff

	data := make([]byte, 1<<13)

	now := time.Now()
	ct := src.NewEncryptor().EncryptNew(witness, data)
	fmt.Println("Gen", time.Since(now))

	data, err := ct.MarshalBinary()
	if err != nil {
		panic(err)
	}

	ct = new(src.Ciphertext)
	if err := ct.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	now = time.Now()
	fmt.Println(src.NewDecryptor().DecryptNew(witness, ct))
	fmt.Println("Verify", time.Since(now))
}
