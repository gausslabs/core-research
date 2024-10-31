package src

import (
	"github.com/Pro7ech/lattigo/ring"
	"github.com/Pro7ech/lattigo/rlwe"
)

// ParametersLiteral the following security: 
// arora-gb             :: rop: ≈2^554.5, m: ≈2^383.3, dreg: 50, t: 10, mem: ≈2^173.6, tag: arora-gb, ↻: ≈2^380.9, ζ: ≈2^14.0, |S|: ≈2^162.4, prop: ≈2^-216.3
// bkw                  :: rop: ≈2^554.9, m: ≈2^540.6, mem: ≈2^541.6, b: 3, t1: 0, t2: 3, ℓ: 2, #cod: 1910, #top: 0, #test: ≈2^13.8, tag: coded-bkw
// usvp                 :: rop: ≈2^343.8, red: ≈2^343.8, δ: 1.001955, β: 1060, d: 31176, tag: usvp
// bdd                  :: rop: ≈2^343.3, red: ≈2^343.3, svp: ≈2^337.0, β: 1058, η: 1098, d: 31330, tag: bdd
// bdd_hybrid           :: rop: ≈2^145.4, red: ≈2^145.3, svp: ≈2^141.9, β: 214, η: 2, ζ: ≈2^13.5, |S|: ≈2^65.8, d: 6551, prob: ≈2^-48.5, ↻: ≈2^50.7, tag: hybrid
// bdd_mitm_hybrid      :: rop: ≈2^127.4, red: ≈2^126.5, svp: ≈2^126.3, β: 227, η: 2, ζ: ≈2^13.5, |S|: ≈2^145.6, d: 6913, prob: ≈2^-25.8, ↻: ≈2^28.0, tag: hybrid
// dual                 :: rop: ≈2^344.8, mem: ≈2^230.6, m: ≈2^13.9, β: 1063, d: 31844, ↻: 1, tag: dual
// dual_hybrid          :: rop: ≈2^337.8, red: ≈2^337.8, guess: ≈2^329.3, β: 1039, p: 2, ζ: 0, t: 310, β': 1039, N: ≈2^201.0, m: ≈2^14.0
var ParametersLiteral = rlwe.ParametersLiteral{
	LogN: 13,
	LogQ: []int{60, 60, 60},
	Xs:   &ring.Ternary{H: 40},
}

const (
	LogWitness = 64
	LogMessage = 8
)
