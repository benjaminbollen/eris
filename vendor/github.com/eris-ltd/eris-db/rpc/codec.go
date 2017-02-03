// Copyright 2015-2017 Monax Industries Limited.
// This file is part of the Monax platform (Monax)

// Monax is free software: you can use, redistribute it and/or modify
// it only under the terms of the GNU General Public License, version
// 3, as published by the Free Software Foundation.

// Monax is distributed WITHOUT ANY WARRANTY pursuant to
// the terms of the Gnu General Public Licence, version 3, including
// (but not limited to) Clause 15 thereof. See the text of the
// GNU General Public License, version 3 for full terms.

// You should have received a copy of the GNU General Public License,
// version 3, with Monax.  If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"io"
)

// Used for rpc request and response data.
type Codec interface {
	EncodeBytes(interface{}) ([]byte, error)
	Encode(interface{}, io.Writer) error
	DecodeBytes(interface{}, []byte) error
	Decode(interface{}, io.Reader) error
}
