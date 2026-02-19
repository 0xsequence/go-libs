package xlog_test

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/test-go/testify/assert"

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

//go:noinline
func makeAlertAttr(err error) any {
	return xlog.Alert(err).Value.Any()
}

func TestAlert(t *testing.T) {
	baseErr := errors.New("sample alert error")
	value := makeAlertAttr(baseErr)
	alertErr, ok := value.(error)
	if !ok {
		t.Fatalf("expected error value, got %T", value)
	}

	assert.Equal(t, "sample alert error", alertErr.Error())

	type stackFramer interface {
		StackFrames() []uintptr
	}
	sf, ok := alertErr.(stackFramer)
	if !ok {
		t.Fatalf("expected alert error to expose stack frames, got %T", alertErr)
	}

	frames := sf.StackFrames()
	resolved := runtime.CallersFrames(frames)
	var firstFunc string
	for {
		frame, more := resolved.Next()
		if frame.Function != "" {
			firstFunc = frame.Function
			break
		}
		if !more {
			break
		}
	}
	if firstFunc == "" {
		t.Fatal("expected at least one resolved stack frame")
	}
	if strings.Contains(firstFunc, "xlog.Alert") {
		t.Fatalf("expected caller frame, got wrapper frame: %q", firstFunc)
	}
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

func TestTokenIDBigInt(t *testing.T) {
	tokenID := big.NewInt(1234567890)
	attr := xlog.TokenIDBigInt(*tokenID)

	assert.Equal(t, "tokenId", attr.Key)
	assert.Equal(t, tokenID.String(), attr.Value.String())
}

func TestDataSource(t *testing.T) {
	source := "sampleSource"
	attr := xlog.DataSource(source)

	assert.Equal(t, "dataSource", attr.Key)
	assert.Equal(t, source, attr.Value.String())
}

func TestDataType(t *testing.T) {
	dataType := "currency"
	attr := xlog.DataType(dataType)

	assert.Equal(t, "dataType", attr.Key)
	assert.Equal(t, dataType, attr.Value.String())
}

func TestProjectID(t *testing.T) {
	projectID := uint64(100)
	attr := xlog.ProjectID(projectID)

	assert.Equal(t, "projectId", attr.Key)
	assert.Equal(t, projectID, attr.Value.Uint64())
}

func TestEcosystemID(t *testing.T) {
	ecosystemID := uint64(100)
	attr := xlog.EcosystemID(ecosystemID)

	assert.Equal(t, "ecosystemId", attr.Key)
	assert.Equal(t, ecosystemID, attr.Value.Uint64())
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
