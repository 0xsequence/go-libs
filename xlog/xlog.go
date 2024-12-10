// Package xlog (eXtended log) extends log/slog with additional helper functions and types.
package xlog

import (
	"fmt"
	"log/slog"
	"math/big"
	"strconv"

	"github.com/0xsequence/go-libs/networks"
	"github.com/0xsequence/go-sequence/lib/prototyp"
)

// slog.Any("error", err)
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

// slog.Any("error", fmt.Errorf(format, args...))
func Errorf(format string, args ...any) slog.Attr {
	return slog.Any("error", fmt.Errorf(format, args...))
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

// slog.Group("network", ChainID(network.ChainID), ChainNetworkName(network.Name))
func ChainIDNetwork(network *networks.Network) slog.Attr {
	if network == nil {
		return slog.Group("network",
			ChainNetworkName("unknown"),
		)
	}

	return slog.Group("network",
		ChainID(network.ChainID),
		ChainNetworkName(network.Name),
	)
}

// slog.String("contractAddress", contractAddress.String())
func ContractAddress(contractAddress prototyp.Hash) slog.Attr {
	return slog.String("contractAddress", contractAddress.String())
}

// slog.String("collectionAddress", collectionAddress.String())
func CollectionAddress(collectionAddress prototyp.Hash) slog.Attr {
	return slog.String("collectionAddress", collectionAddress.String())
}

// slog.String("currencyAddress", currencyAddress.String())
func CurrencyAddress(currencyAddress prototyp.Hash) slog.Attr {
	return slog.String("currencyAddress", currencyAddress.String())
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
func TokenID(tokenID prototyp.BigInt) slog.Attr {
	return TokenIDString(tokenID.String())
}

// slog.String("tokenId", tokenID.String())
func TokenIDBigInt(tokenID big.Int) slog.Attr {
	return TokenIDString(tokenID.String())
}

// slog.String("entity", tokenID.String())
func Entity(entity string) slog.Attr {
	return slog.String("entity", entity)
}

// slog.String("source", source)
func Source(source string) slog.Attr {
	return slog.String("source", source)
}

// slog.Uint64("projectId", projectID)
func ProjectID(projectID uint64) slog.Attr {
	return slog.Uint64("projectId", projectID)
}

func Stringer[T fmt.Stringer](k string, v T) slog.Attr {
	return slog.String(k, v.String())
}
