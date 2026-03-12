// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// OracleVerifierMetaData contains all meta data concerning the OracleVerifier contract.
var OracleVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_oracles\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"_f\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_queueAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"submitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vectorLength\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"JobCompleted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"expectedConfigDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"f\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_jobId\",\"type\":\"uint256\"}],\"name\":\"getResult\",\"outputs\":[{\"internalType\":\"int128[]\",\"name\":\"\",\"type\":\"int128[]\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_jobId\",\"type\":\"uint256\"}],\"name\":\"isCompleted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleQueue\",\"outputs\":[{\"internalType\":\"contractIOracleQueue\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"results\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"submitter\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"saved\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"configDigest\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"seqNr\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"},{\"internalType\":\"bytes32[]\",\"name\":\"rs\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"ss\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32\",\"name\":\"rawVs\",\"type\":\"bytes32\"}],\"name\":\"transmit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"usedSeqNr\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OracleVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleVerifierMetaData.ABI instead.
var OracleVerifierABI = OracleVerifierMetaData.ABI

// OracleVerifier is an auto generated Go binding around an Ethereum contract.
type OracleVerifier struct {
	OracleVerifierCaller     // Read-only binding to the contract
	OracleVerifierTransactor // Write-only binding to the contract
	OracleVerifierFilterer   // Log filterer for contract events
}

// OracleVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleVerifierSession struct {
	Contract     *OracleVerifier   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleVerifierCallerSession struct {
	Contract *OracleVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OracleVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleVerifierTransactorSession struct {
	Contract     *OracleVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OracleVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleVerifierRaw struct {
	Contract *OracleVerifier // Generic contract binding to access the raw methods on
}

// OracleVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleVerifierCallerRaw struct {
	Contract *OracleVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// OracleVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleVerifierTransactorRaw struct {
	Contract *OracleVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleVerifier creates a new instance of OracleVerifier, bound to a specific deployed contract.
func NewOracleVerifier(address common.Address, backend bind.ContractBackend) (*OracleVerifier, error) {
	contract, err := bindOracleVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleVerifier{OracleVerifierCaller: OracleVerifierCaller{contract: contract}, OracleVerifierTransactor: OracleVerifierTransactor{contract: contract}, OracleVerifierFilterer: OracleVerifierFilterer{contract: contract}}, nil
}

// NewOracleVerifierCaller creates a new read-only instance of OracleVerifier, bound to a specific deployed contract.
func NewOracleVerifierCaller(address common.Address, caller bind.ContractCaller) (*OracleVerifierCaller, error) {
	contract, err := bindOracleVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleVerifierCaller{contract: contract}, nil
}

// NewOracleVerifierTransactor creates a new write-only instance of OracleVerifier, bound to a specific deployed contract.
func NewOracleVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleVerifierTransactor, error) {
	contract, err := bindOracleVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleVerifierTransactor{contract: contract}, nil
}

// NewOracleVerifierFilterer creates a new log filterer instance of OracleVerifier, bound to a specific deployed contract.
func NewOracleVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleVerifierFilterer, error) {
	contract, err := bindOracleVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleVerifierFilterer{contract: contract}, nil
}

// bindOracleVerifier binds a generic wrapper to an already deployed contract.
func bindOracleVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleVerifier *OracleVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleVerifier.Contract.OracleVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleVerifier *OracleVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleVerifier.Contract.OracleVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleVerifier *OracleVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleVerifier.Contract.OracleVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleVerifier *OracleVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleVerifier *OracleVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleVerifier *OracleVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleVerifier.Contract.contract.Transact(opts, method, params...)
}

// ExpectedConfigDigest is a free data retrieval call binding the contract method 0xa3a6034f.
//
// Solidity: function expectedConfigDigest() view returns(bytes32)
func (_OracleVerifier *OracleVerifierCaller) ExpectedConfigDigest(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "expectedConfigDigest")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ExpectedConfigDigest is a free data retrieval call binding the contract method 0xa3a6034f.
//
// Solidity: function expectedConfigDigest() view returns(bytes32)
func (_OracleVerifier *OracleVerifierSession) ExpectedConfigDigest() ([32]byte, error) {
	return _OracleVerifier.Contract.ExpectedConfigDigest(&_OracleVerifier.CallOpts)
}

// ExpectedConfigDigest is a free data retrieval call binding the contract method 0xa3a6034f.
//
// Solidity: function expectedConfigDigest() view returns(bytes32)
func (_OracleVerifier *OracleVerifierCallerSession) ExpectedConfigDigest() ([32]byte, error) {
	return _OracleVerifier.Contract.ExpectedConfigDigest(&_OracleVerifier.CallOpts)
}

// F is a free data retrieval call binding the contract method 0x26121ff0.
//
// Solidity: function f() view returns(uint8)
func (_OracleVerifier *OracleVerifierCaller) F(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "f")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// F is a free data retrieval call binding the contract method 0x26121ff0.
//
// Solidity: function f() view returns(uint8)
func (_OracleVerifier *OracleVerifierSession) F() (uint8, error) {
	return _OracleVerifier.Contract.F(&_OracleVerifier.CallOpts)
}

// F is a free data retrieval call binding the contract method 0x26121ff0.
//
// Solidity: function f() view returns(uint8)
func (_OracleVerifier *OracleVerifierCallerSession) F() (uint8, error) {
	return _OracleVerifier.Contract.F(&_OracleVerifier.CallOpts)
}

// GetResult is a free data retrieval call binding the contract method 0x995e4339.
//
// Solidity: function getResult(uint256 _jobId) view returns(int128[], address, uint256)
func (_OracleVerifier *OracleVerifierCaller) GetResult(opts *bind.CallOpts, _jobId *big.Int) ([]*big.Int, common.Address, *big.Int, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "getResult", _jobId)

	if err != nil {
		return *new([]*big.Int), *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetResult is a free data retrieval call binding the contract method 0x995e4339.
//
// Solidity: function getResult(uint256 _jobId) view returns(int128[], address, uint256)
func (_OracleVerifier *OracleVerifierSession) GetResult(_jobId *big.Int) ([]*big.Int, common.Address, *big.Int, error) {
	return _OracleVerifier.Contract.GetResult(&_OracleVerifier.CallOpts, _jobId)
}

// GetResult is a free data retrieval call binding the contract method 0x995e4339.
//
// Solidity: function getResult(uint256 _jobId) view returns(int128[], address, uint256)
func (_OracleVerifier *OracleVerifierCallerSession) GetResult(_jobId *big.Int) ([]*big.Int, common.Address, *big.Int, error) {
	return _OracleVerifier.Contract.GetResult(&_OracleVerifier.CallOpts, _jobId)
}

// IsCompleted is a free data retrieval call binding the contract method 0x7a41984b.
//
// Solidity: function isCompleted(uint256 _jobId) view returns(bool)
func (_OracleVerifier *OracleVerifierCaller) IsCompleted(opts *bind.CallOpts, _jobId *big.Int) (bool, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "isCompleted", _jobId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsCompleted is a free data retrieval call binding the contract method 0x7a41984b.
//
// Solidity: function isCompleted(uint256 _jobId) view returns(bool)
func (_OracleVerifier *OracleVerifierSession) IsCompleted(_jobId *big.Int) (bool, error) {
	return _OracleVerifier.Contract.IsCompleted(&_OracleVerifier.CallOpts, _jobId)
}

// IsCompleted is a free data retrieval call binding the contract method 0x7a41984b.
//
// Solidity: function isCompleted(uint256 _jobId) view returns(bool)
func (_OracleVerifier *OracleVerifierCallerSession) IsCompleted(_jobId *big.Int) (bool, error) {
	return _OracleVerifier.Contract.IsCompleted(&_OracleVerifier.CallOpts, _jobId)
}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address ) view returns(bool)
func (_OracleVerifier *OracleVerifierCaller) IsOracle(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "isOracle", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address ) view returns(bool)
func (_OracleVerifier *OracleVerifierSession) IsOracle(arg0 common.Address) (bool, error) {
	return _OracleVerifier.Contract.IsOracle(&_OracleVerifier.CallOpts, arg0)
}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address ) view returns(bool)
func (_OracleVerifier *OracleVerifierCallerSession) IsOracle(arg0 common.Address) (bool, error) {
	return _OracleVerifier.Contract.IsOracle(&_OracleVerifier.CallOpts, arg0)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(address)
func (_OracleVerifier *OracleVerifierCaller) OracleQueue(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "oracleQueue")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(address)
func (_OracleVerifier *OracleVerifierSession) OracleQueue() (common.Address, error) {
	return _OracleVerifier.Contract.OracleQueue(&_OracleVerifier.CallOpts)
}

// OracleQueue is a free data retrieval call binding the contract method 0x69a4246d.
//
// Solidity: function oracleQueue() view returns(address)
func (_OracleVerifier *OracleVerifierCallerSession) OracleQueue() (common.Address, error) {
	return _OracleVerifier.Contract.OracleQueue(&_OracleVerifier.CallOpts)
}

// Results is a free data retrieval call binding the contract method 0x1b0c27da.
//
// Solidity: function results(uint256 ) view returns(address submitter, uint256 timestamp, bool saved)
func (_OracleVerifier *OracleVerifierCaller) Results(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Submitter common.Address
	Timestamp *big.Int
	Saved     bool
}, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "results", arg0)

	outstruct := new(struct {
		Submitter common.Address
		Timestamp *big.Int
		Saved     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Submitter = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Saved = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Results is a free data retrieval call binding the contract method 0x1b0c27da.
//
// Solidity: function results(uint256 ) view returns(address submitter, uint256 timestamp, bool saved)
func (_OracleVerifier *OracleVerifierSession) Results(arg0 *big.Int) (struct {
	Submitter common.Address
	Timestamp *big.Int
	Saved     bool
}, error) {
	return _OracleVerifier.Contract.Results(&_OracleVerifier.CallOpts, arg0)
}

// Results is a free data retrieval call binding the contract method 0x1b0c27da.
//
// Solidity: function results(uint256 ) view returns(address submitter, uint256 timestamp, bool saved)
func (_OracleVerifier *OracleVerifierCallerSession) Results(arg0 *big.Int) (struct {
	Submitter common.Address
	Timestamp *big.Int
	Saved     bool
}, error) {
	return _OracleVerifier.Contract.Results(&_OracleVerifier.CallOpts, arg0)
}

// UsedSeqNr is a free data retrieval call binding the contract method 0x29624283.
//
// Solidity: function usedSeqNr(uint64 ) view returns(bool)
func (_OracleVerifier *OracleVerifierCaller) UsedSeqNr(opts *bind.CallOpts, arg0 uint64) (bool, error) {
	var out []interface{}
	err := _OracleVerifier.contract.Call(opts, &out, "usedSeqNr", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedSeqNr is a free data retrieval call binding the contract method 0x29624283.
//
// Solidity: function usedSeqNr(uint64 ) view returns(bool)
func (_OracleVerifier *OracleVerifierSession) UsedSeqNr(arg0 uint64) (bool, error) {
	return _OracleVerifier.Contract.UsedSeqNr(&_OracleVerifier.CallOpts, arg0)
}

// UsedSeqNr is a free data retrieval call binding the contract method 0x29624283.
//
// Solidity: function usedSeqNr(uint64 ) view returns(bool)
func (_OracleVerifier *OracleVerifierCallerSession) UsedSeqNr(arg0 uint64) (bool, error) {
	return _OracleVerifier.Contract.UsedSeqNr(&_OracleVerifier.CallOpts, arg0)
}

// Transmit is a paid mutator transaction binding the contract method 0xf957c546.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OracleVerifier *OracleVerifierTransactor) Transmit(opts *bind.TransactOpts, configDigest [32]byte, seqNr uint64, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OracleVerifier.contract.Transact(opts, "transmit", configDigest, seqNr, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xf957c546.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OracleVerifier *OracleVerifierSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OracleVerifier.Contract.Transmit(&_OracleVerifier.TransactOpts, configDigest, seqNr, report, rs, ss, rawVs)
}

// Transmit is a paid mutator transaction binding the contract method 0xf957c546.
//
// Solidity: function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes32[] rs, bytes32[] ss, bytes32 rawVs) returns()
func (_OracleVerifier *OracleVerifierTransactorSession) Transmit(configDigest [32]byte, seqNr uint64, report []byte, rs [][32]byte, ss [][32]byte, rawVs [32]byte) (*types.Transaction, error) {
	return _OracleVerifier.Contract.Transmit(&_OracleVerifier.TransactOpts, configDigest, seqNr, report, rs, ss, rawVs)
}

// OracleVerifierJobCompletedIterator is returned from FilterJobCompleted and is used to iterate over the raw logs and unpacked data for JobCompleted events raised by the OracleVerifier contract.
type OracleVerifierJobCompletedIterator struct {
	Event *OracleVerifierJobCompleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleVerifierJobCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleVerifierJobCompleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleVerifierJobCompleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleVerifierJobCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleVerifierJobCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleVerifierJobCompleted represents a JobCompleted event raised by the OracleVerifier contract.
type OracleVerifierJobCompleted struct {
	JobId        *big.Int
	Submitter    common.Address
	VectorLength *big.Int
	Timestamp    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterJobCompleted is a free log retrieval operation binding the contract event 0xbeccb946f65bcbea37397b769e33a37dabec51c94907407fd5f2eaa81f1a0bfa.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed submitter, uint256 vectorLength, uint256 timestamp)
func (_OracleVerifier *OracleVerifierFilterer) FilterJobCompleted(opts *bind.FilterOpts, jobId []*big.Int, submitter []common.Address) (*OracleVerifierJobCompletedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var submitterRule []interface{}
	for _, submitterItem := range submitter {
		submitterRule = append(submitterRule, submitterItem)
	}

	logs, sub, err := _OracleVerifier.contract.FilterLogs(opts, "JobCompleted", jobIdRule, submitterRule)
	if err != nil {
		return nil, err
	}
	return &OracleVerifierJobCompletedIterator{contract: _OracleVerifier.contract, event: "JobCompleted", logs: logs, sub: sub}, nil
}

// WatchJobCompleted is a free log subscription operation binding the contract event 0xbeccb946f65bcbea37397b769e33a37dabec51c94907407fd5f2eaa81f1a0bfa.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed submitter, uint256 vectorLength, uint256 timestamp)
func (_OracleVerifier *OracleVerifierFilterer) WatchJobCompleted(opts *bind.WatchOpts, sink chan<- *OracleVerifierJobCompleted, jobId []*big.Int, submitter []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var submitterRule []interface{}
	for _, submitterItem := range submitter {
		submitterRule = append(submitterRule, submitterItem)
	}

	logs, sub, err := _OracleVerifier.contract.WatchLogs(opts, "JobCompleted", jobIdRule, submitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleVerifierJobCompleted)
				if err := _OracleVerifier.contract.UnpackLog(event, "JobCompleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseJobCompleted is a log parse operation binding the contract event 0xbeccb946f65bcbea37397b769e33a37dabec51c94907407fd5f2eaa81f1a0bfa.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed submitter, uint256 vectorLength, uint256 timestamp)
func (_OracleVerifier *OracleVerifierFilterer) ParseJobCompleted(log types.Log) (*OracleVerifierJobCompleted, error) {
	event := new(OracleVerifierJobCompleted)
	if err := _OracleVerifier.contract.UnpackLog(event, "JobCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
