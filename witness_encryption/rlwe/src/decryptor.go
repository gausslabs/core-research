package src

import (
	"fmt"

	"github.com/Pro7ech/lattigo/ring"
	"github.com/Pro7ech/lattigo/rlwe"
)

type Decryptor struct {
	params rlwe.Parameters
	buf    ring.RNSPoly
}

func NewDecryptor() *Decryptor {

	params, err := rlwe.NewParametersFromLiteral(ParametersLiteral)
	if err != nil {
		panic(err)
	}

	return &Decryptor{
		params: params,
		buf:    params.RingQ().NewRNSPoly(),
	}
}

func (dec *Decryptor) DecryptNew(witness uint64, ct *Ciphertext) (data []byte) {

	if ct.D == nil {
		ct.Expand(dec.params)
	}

	Rq := dec.params.RingQ()
	buff := Rq.NewRNSPoly()

	for i := range LogWitness {

		d := ct.D[i][witness&1]

		if i == 0 {
			Rq.MulCoeffsMontgomery(ct.C, d, buff)
		} else {
			Rq.MulCoeffsMontgomery(buff, d, buff)
		}

		witness >>= 1
	}

	Rq.INTT(buff, buff)

	Rq.DivRoundByLastModulusMany(buff.Level(), buff, Rq.NewRNSPoly(), buff)
	buff.Resize(0)

	Rq = Rq.AtLevel(0)

	Rq.Sub(buff, ct.A, buff)

	fmt.Println(Rq.Stats(buff))

	t := uint64(1 << LogMessage)

	Rq.MulScalar(buff, t, buff)
	Rq.AddScalar(buff, Rq[0].Modulus>>1, buff)
	data = make([]byte, Rq.N())
	for i := range buff[0] {
		data[i] = -uint8(buff[0][i])
	}

	return
}
