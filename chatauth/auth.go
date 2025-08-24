// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package chatauth

type SecretSalter interface {
	SaltSecret(secret []byte) []byte
}

func NopSalter() SecretSalter { return new(nopSalter) }

type Signer interface {
	Sign(claims *Claims, salter SecretSalter) (string, error)
}

type Parser interface {
	Parse(token string, salter SecretSalter) (Claims, error)
}

type SignerParser interface {
	Signer
	Parser
}

type nopSalter struct{}

func (nopSalter) SaltSecret(s []byte) []byte { return s }
