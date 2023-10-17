package appenders

import (
	"math/big"
	"strconv"
)

type IntText int64

func (a IntText) Append(out []byte) []byte {
	i64 := int64(a)
	return strconv.AppendInt(out, i64, 10)
}

var _ Appender = IntText(0)

type UintText uint64

func (a UintText) Append(out []byte) []byte {
	u64 := uint64(a)
	return strconv.AppendUint(out, u64, 10)
}

var _ Appender = UintText(0)

type BigIntText struct {
	Pointer *big.Int
}

func (a BigIntText) Append(out []byte) []byte {
	return a.Pointer.Append(out, 10)
}

var _ Appender = BigIntText{}

type FloatText float64

func (a FloatText) Append(out []byte) []byte {
	f64 := float64(a)
	return strconv.AppendFloat(out, f64, 'g', -1, 64)
}

var _ Appender = FloatText(0)

type BigFloatText struct {
	Pointer *big.Float

	Format           byte
	Prec             int
	HasFormatAndPrec bool
}

func (a BigFloatText) Append(out []byte) []byte {
	if a.HasFormatAndPrec {
		return a.Pointer.Append(out, a.Format, a.Prec)
	}
	return a.Pointer.Append(out, 'g', -1)
}

var _ Appender = BigFloatText{}
