package src

import (
	"fmt"
	"math/big"

	"github.com/Pro7ech/lattigo/ring"
	"github.com/Pro7ech/lattigo/rlwe"
	"github.com/Pro7ech/lattigo/utils/sampling"
)

type Encryptor struct {
	params rlwe.Parameters
	buff   [2]ring.RNSPoly
}

func NewEncryptor() *Encryptor {
	params, err := rlwe.NewParametersFromLiteral(ParametersLiteral)
	if err != nil {
		panic(err)
	}

	return &Encryptor{
		params: params,
		buff:   [2]ring.RNSPoly{params.RingQ().NewRNSPoly(), params.RingQ().NewRNSPoly()},
	}
}

func (enc *Encryptor) EncryptNew(witness uint64, data []byte) (ct *Ciphertext) {

	params := enc.params

	if len(data) > params.N() {
		panic(fmt.Errorf("len(data)=%d > params.N()=%d", len(data), params.N()))
	}

	Xs, err := ring.NewSampler(sampling.NewSource(sampling.NewSeed()), params.Q(), params.Xs())
	if err != nil {
		panic(err)
	}

	A := ring.NewUniformSampler(sampling.NewSource(sampling.NewSeed()), params.Q()).ReadNew(params.N())
	Xe, err := ring.NewSampler(sampling.NewSource(sampling.NewSeed()), params.Q(), params.Xe())
	if err != nil {
		panic(err)
	}

	seed := sampling.NewSeed()
	source := sampling.NewSource(seed)

	XD, err := ring.NewSampler(nil, params.Q(), params.Xs())
	if err != nil {
		panic(err)
	}

	Rq := params.RingQ()

	C := *A.Clone()

	buf0, buf1 := enc.buff[0], enc.buff[1]

	for _ = range LogWitness {
		trapdoor(C, A, buf0, buf1, Rq, witness&1 == 1, Xs, Xe, XD, source)
		witness >>= 1
	}

	Xe.Read(buf0)
	Rq.NTT(buf0, buf0)
	Rq.Add(A, buf0, A)
	Rq.INTT(A, A)

	Rq.DivRoundByLastModulusMany(A.Level(), A, buf0, A)

	msg := Rq.NewRNSPoly()

	coeffs := msg[0]
	for i := range data {
		coeffs[i] = uint64(data[i])
	}

	Rq = Rq.AtLevel(0)

	Rq.MulScalarBigint(msg, new(big.Int).ModInverse(new(big.Int).SetUint64(1<<LogMessage), Rq.Modulus()), msg)

	Rq.Add(A, msg, A)

	return &Ciphertext{
		Seed: seed,
		C:    C,
		A:    A,
	}
}

func trapdoor(a, A, buf0, buf1 ring.RNSPoly, Rq ring.RNSRing, w bool, Xs, Xe, XD ring.Sampler, source *sampling.Source) {

	Xs.Read(buf0)
	Rq.MForm(buf0, buf0)
	Rq.NTT(buf0, buf0)

	// A * s
	Rq.MulCoeffsMontgomery(A, buf0, A)

	// a * s
	Rq.MulCoeffsMontgomery(a, buf0, a)

	// a * s + e
	Xe.Read(buf0)
	Rq.NTT(buf0, buf0)
	Rq.Add(a, buf0, a)

	var seed [32]byte
	if w {
		_ = source.NewSeed()
		seed = source.NewSeed()
	} else {
		seed = source.NewSeed()
		_ = source.NewSeed()
	}

	genDInv(Rq, buf0, buf1, XD.WithSource(sampling.NewSource(seed)))
	Rq.MulCoeffsMontgomery(a, buf1, a)
}

func genD(Rq ring.RNSRing, D ring.RNSPoly, Xe ring.Sampler) {

	for {

		Xe.Read(D)

		Rq.NTT(D, D)

		var isInvertible bool = true
		for i := range D {
			if !isInvertible {
				break
			}
			for j := range D[i] {
				if D[i][j] == 0 {
					isInvertible = false
					break
				}
			}
		}

		if isInvertible {
			break
		}
	}

	Rq.MForm(D, D)
}

func genDInv(Rq ring.RNSRing, D, DInv ring.RNSPoly, Xe ring.Sampler) {

	genD(Rq, D, Xe)

	DInv.Ones()
	Rq.MForm(DInv, DInv)

	for i, s := range Rq {

		qi := s.Modulus

		for j := qi - 2; j > 0; j >>= 1 {
			if j&1 == 1 {
				s.MulCoeffsMontgomery(DInv[i], D[i], DInv[i])
			}

			s.MulCoeffsMontgomery(D[i], D[i], D[i])
		}
	}

	return
}
