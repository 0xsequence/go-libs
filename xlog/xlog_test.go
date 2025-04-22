package xlog_test

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/test-go/testify/assert"

	"github.com/0xsequence/go-libs/networks"
	"github.com/0xsequence/go-libs/xlog"
)

func TestError(t *testing.T) {
	err := errors.New("sample error")
	attr := xlog.Error(err)

	assert.Equal(t, "error", attr.Key)
	assert.Equal(t, err, attr.Value.Any())
}

func TestErrorf(t *testing.T) {
	format := "formatted error %s"
	args := []any{"test"}
	attr := xlog.Errorf(format, args...)

	assert.Equal(t, "error", attr.Key)
	assert.Equal(t, fmt.Sprintf(format, args...), attr.Value.String())
}

func TestID(t *testing.T) {
	id := uint64(12345)
	attr := xlog.ID(id)

	assert.Equal(t, "id", attr.Key)
	assert.Equal(t, id, attr.Value.Uint64())
}

func TestChainID(t *testing.T) {
	chainID := uint64(67890)
	attr := xlog.ChainID(chainID)

	assert.Equal(t, "chainId", attr.Key)
	assert.Equal(t, chainID, attr.Value.Uint64())
}

func TestChainIDString(t *testing.T) {
	chainID := "12345"
	attr := xlog.ChainIDString(chainID)

	assert.Equal(t, "chainId", attr.Key)
	expectedID, _ := strconv.ParseUint(chainID, 10, 64)
	assert.Equal(t, expectedID, attr.Value.Uint64())
}

func TestChainNetworkName(t *testing.T) {
	name := "NetworkName"
	attr := xlog.ChainNetworkName(name)

	assert.Equal(t, "name", attr.Key)
	assert.Equal(t, name, attr.Value.String())
}

func TestChainIDNetwork(t *testing.T) {
	network := &networks.Network{
		ChainID: 12345,
		Name:    "NetworkName",
	}
	attr := xlog.ChainIDNetwork(network)

	assert.Equal(t, "network", attr.Key)
	// assert individual fields (this is a bit dependent on how you handle logging with slog)
}

func TestContractAddress(t *testing.T) {
	address := prototyp.Hash("sampleAddress")
	attr := xlog.ContractAddress(address)

	assert.Equal(t, "contractAddress", attr.Key)
	assert.Equal(t, address.String(), attr.Value.String())
}

func TestCollectionAddress(t *testing.T) {
	address := prototyp.Hash("sampleCollectionAddress")
	attr := xlog.CollectionAddress(address)

	assert.Equal(t, "collectionAddress", attr.Key)
	assert.Equal(t, address.String(), attr.Value.String())
}

func TestCurrencyAddress(t *testing.T) {
	address := prototyp.Hash("sampleCurrencyAddress")
	attr := xlog.CurrencyAddress(address)

	assert.Equal(t, "currencyAddress", attr.Key)
	assert.Equal(t, address.String(), attr.Value.String())
}

func TestOrderID(t *testing.T) {
	orderID := "sampleOrderID"
	attr := xlog.OrderID(orderID)

	assert.Equal(t, "orderID", attr.Key)
	assert.Equal(t, orderID, attr.Value.String())
}

func TestTokenIDString(t *testing.T) {
	tokenID := "sampleTokenID"
	attr := xlog.TokenIDString(tokenID)

	assert.Equal(t, "tokenId", attr.Key)
	assert.Equal(t, tokenID, attr.Value.String())
}

func TestTokenID(t *testing.T) {
	tokenID := prototyp.NewBigIntFromDecimalString("1234567890")
	attr := xlog.TokenID(tokenID)

	assert.Equal(t, "tokenId", attr.Key)
	assert.Equal(t, tokenID.String(), attr.Value.String())
}

func TestTokenIDPtr(t *testing.T) {
	tokenID := prototyp.NewBigIntFromDecimalString("1234567890")
	attr := xlog.TokenIDPtr(&tokenID)

	assert.Equal(t, "tokenId", attr.Key)
	assert.Equal(t, tokenID.String(), attr.Value.String())

	// Test with nil pointer
	var nilTokenID *prototyp.BigInt
	attrNil := xlog.TokenIDPtr(nilTokenID)

	assert.Equal(t, "tokenId", attrNil.Key)
	assert.Equal(t, "empty", attrNil.Value.String())
}

func TestTokenIDBigInt(t *testing.T) {
	tokenID := big.NewInt(1234567890)
	attr := xlog.TokenIDBigInt(*tokenID)

	assert.Equal(t, "tokenId", attr.Key)
	assert.Equal(t, tokenID.String(), attr.Value.String())
}

func TestEntity(t *testing.T) {
	entity := "sampleEntity"
	attr := xlog.Entity(entity)

	assert.Equal(t, "entity", attr.Key)
	assert.Equal(t, entity, attr.Value.String())
}

func TestSource(t *testing.T) {
	source := "sampleSource"
	attr := xlog.Source(source)

	assert.Equal(t, "source", attr.Key)
	assert.Equal(t, source, attr.Value.String())
}

func TestProjectID(t *testing.T) {
	projectID := uint64(100)
	attr := xlog.ProjectID(projectID)

	assert.Equal(t, "projectId", attr.Key)
	assert.Equal(t, projectID, attr.Value.Uint64())
}

type sampleStringer struct{}

func (s sampleStringer) String() string { return "sampleString" }

func TestStringer(t *testing.T) {
	str := sampleStringer{}
	attr := xlog.Stringer("sampleKey", str)

	assert.Equal(t, "sampleKey", attr.Key)
	assert.Equal(t, "sampleString", attr.Value.String())
}

func TestPointerSlice(t *testing.T) {
	type sampleType struct {
		ID int
	}
	slice := []*sampleType{
		{ID: 1},
		nil,
		{ID: 2},
	}
	attr := xlog.PointerSlice("sampleSlice", slice)

	expected := "[&{ID:1} <nil> &{ID:2}]"
	assert.Equal(t, "sampleSlice", attr.Key)
	assert.Equal(t, expected, attr.Value.String())
}
