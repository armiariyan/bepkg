package dh

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

// P Big Prime
var P *big.Int

// G Generator
var G *big.Int

func init() {
	//Default Prime and Generator
	P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3DC2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F83655D23DCA3AD961C62F356208552BB9ED529077096966D670C354E4ABC9804F1746C08CA18217C32905E462E36CE3BE39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9DE2BCBF6955817183995497CEA956AE515D2261898FA051015728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6BF12FFA06D98A0864D87602733EC86A64521F2B18177B200CBBE117577A615D6C770988C0BAD946E208E24FA074E5AB3143DB5BFCE0FD108E4B82D120A93AD2CAFFFFFFFFFFFFFFFF", 16)
	G = big.NewInt(2)
}

// IsSafePrime returns true, if the prime of the group is
// a so called safe-prime. For a group with a safe-prime prime
// number the Decisional-Diffie-Hellman-Problem (DDH) is a
// 'hard' problem. The n argument is the number of iterations
// for the probabilistic prime test.
// It's recommend to use DDH-safe groups for DH-exchanges.
func IsSafePrime(n int) bool {
	q := new(big.Int).Sub(P, big.NewInt(1))
	q = q.Div(q, big.NewInt(2))
	return q.ProbablyPrime(n)
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(random io.Reader) (private *big.Int, public *big.Int, err error) {
	if random == nil {
		random = rand.Reader
	}

	// Ensure, that p.G ^ privateKey > than g.P
	// (only modulo calculations are safe)
	// The minimal (and common) value for p.G is 2
	// So 2 ^ (1 + 'bitsize of p.G') > than g.P
	min := big.NewInt(int64(P.BitLen() + 1))
	bytes := make([]byte, (P.BitLen()+7)/8)

	for private == nil {
		_, err = io.ReadFull(random, bytes)
		if err != nil {
			private = nil
			return
		}
		// Clear bits in the first byte to increase
		// the probability that the candidate is < g.P.
		bytes[0] = 0
		if private == nil {
			private = new(big.Int)
		}
		(*private).SetBytes(bytes)
		if (*private).Cmp(min) < 0 {
			private = nil
		}
	}

	public = new(big.Int).Exp(G, private, P)
	return
}

// PublicKey returns the public key corresponding to the given private one.
func PublicKey(private *big.Int) (public *big.Int) {
	public = new(big.Int).Exp(G, private, P)
	return
}

// Check returns a non-nil error if the given public key is
// not a possible element of the group. This means, that the
// public key is < 0 or > g.P.
func Check(peersPublic *big.Int) (err error) {
	if !((*peersPublic).Cmp(big.NewInt(0)) >= 0 && (*peersPublic).Cmp(P) == -1) {
		err = errors.New("peer's public is not a possible group element")
	}
	return
}

// ComputeSecret returns the secret computed from
// the own private and the peer's public key.
func ComputeSecret(private *big.Int, peersPublic *big.Int) (secret *big.Int) {
	secret = new(big.Int).Exp(peersPublic, private, P)
	return
}
