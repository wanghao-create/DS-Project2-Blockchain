package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const maxNonce = math.MaxInt64
const difficulty = 24
const printInterval = 10

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// when verifying compare the hash result with target by the former 'difficulty' bits
	target = target.Lsh(target, uint(256-difficulty))
	temp := &ProofOfWork{b, target}
	return temp
}

// the process of mining
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block")

	for nonce < maxNonce {
		// merge the data to do hash on it
		data := bytes.Join(
			[][]byte{
				pow.block.PrevBlockHash,
				pow.block.Data,
				[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
				[]byte(strconv.FormatInt(int64(difficulty), 16)),
				[]byte(strconv.FormatInt(int64(nonce), 16)),
			},
			[]byte{},
		)

		// compute the hash result and check
		hash = sha256.Sum256(data)
		if math.Remainder(float64(nonce), printInterval) == 0 {
			fmt.Printf("\rCurrent trying: %x", hash)
		}
		hashInt.SetBytes(hash[:])
		// transform the comparasion of bits to comparision of big ints
		// may be much slower than byte comparasion
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\rCurrent trying: %x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			[]byte(strconv.FormatInt(pow.block.Timestamp, 16)),
			[]byte(strconv.FormatInt(int64(difficulty), 16)),
			[]byte(strconv.FormatInt(int64(pow.block.Nonce), 16)),
		},
		[]byte{},
	)

	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	temp := hashInt.Cmp(pow.target) == -1

	return temp
}
