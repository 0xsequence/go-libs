// Package xlog (eXtended log) extends log/slog with additional helper functions and types.
package xlog

import (
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xsequence/go-libs/alert"
)

// slog.Any("error", err)
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

// slog.Any("error", fmt.Errorf(format, args...))
func Errorf(format string, args ...any) slog.Attr {
	return slog.Any("error", fmt.Errorf(format, args...))
}

// slog.Any("error", alert.Error(err))
func Alert(err error) slog.Attr {
	return slog.Any("error", alert.ErrorSkip(2, err))
}

// slog.Any("error", alert.Errorf(format, args...)))
func Alertf(format string, args ...any) slog.Attr {
	return slog.Any("error", alert.ErrorSkip(2, fmt.Errorf(format, args...)))
}

// slog.Uint64("id", ID)
func ID(ID uint64) slog.Attr {
	return slog.Uint64("id", ID)
}

// slog.Uint64("chainId", chainID)
func ChainID(chainID uint64) slog.Attr {
	return slog.Uint64("chainId", chainID)
}

// slog.Uint64("chainId", chainID)
func ChainIDString(chainID string) slog.Attr {
	id, _ := strconv.ParseUint(chainID, 10, 64)
	return ChainID(id)
}

// slog.String("name", name)
func ChainNetworkName(name string) slog.Attr {
	return slog.String("name", name)
}

// slog.String("orderID", orderID)
func OrderID(orderID string) slog.Attr {
	return slog.String("orderID", orderID)
}

// slog.String("tokenId", tokenID)
func TokenIDString(tokenID string) slog.Attr {
	return slog.String("tokenId", tokenID)
}

// slog.String("tokenId", tokenID.String())
func TokenIDBigInt(tokenID big.Int) slog.Attr {
	return TokenIDString(tokenID.String())
}

// slog.String("dataType", "currency")
func DataType(dataType string) slog.Attr {
	return slog.String("dataType", dataType)
}

// slog.String("dataSource", source)
func DataSource(dataSource string) slog.Attr {
	return slog.String("dataSource", dataSource)
}

// slog.Uint64("projectId", projectID)
func ProjectID(projectID uint64) slog.Attr {
	return slog.Uint64("projectId", projectID)
}

// slog.Uint64("ecosystemId", projectID)
func EcosystemID(ecosystemID uint64) slog.Attr {
	return slog.Uint64("ecosystemId", ecosystemID)
}

func Stringer[T fmt.Stringer](k string, v T) slog.Attr {
	return slog.String(k, v.String())
}

func PointerSlice[T any](key string, slice []*T) slog.Attr {
	var b strings.Builder
	b.WriteString("[")

	for i, item := range slice {
		if i > 0 {
			b.WriteString(" ")
		}

		if item == nil {
			b.WriteString("<nil>")
		} else {
			b.WriteString(fmt.Sprintf("&%+v", *item))
		}
	}

	b.WriteString("]")

	return slog.String(key, b.String())
}
