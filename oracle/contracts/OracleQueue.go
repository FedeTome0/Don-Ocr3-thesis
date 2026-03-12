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

// OracleQueueMetaData contains all meta data concerning the OracleQueue contract.
var OracleQueueMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_queryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_oracleReward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"}],\"name\":\"LogNewCustomerRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"jobId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"}],\"name\":\"LogNewJobForOracles\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"approveJob\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"customerQueue\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"payment\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isProcessed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"modelCreator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleJobCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracleQueue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"originalRequestId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"ipfsCid\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleVerifierAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"queryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_ipfsCid\",\"type\":\"string\"}],\"name\":\"requestAttribution\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"rewardOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_verifier\",\"type\":\"address\"}],\"name\":\"setVerifierAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OracleQueueABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleQueueMetaData.ABI instead.
var OracleQueueABI = OracleQueueMetaData.ABI

// OracleQueue is an auto generated Go binding around an Ethereum contract.
type OracleQueue struct {
	OracleQueueCaller     // Read-only binding to the contract
	OracleQueueTransactor // Write-only binding to the contract
	OracleQueueFilterer   // Log filterer for contract events
}

// OracleQueueCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleQueueCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleQueueTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleQueueTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleQueueFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleQueueFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleQueueSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleQueueSession struct {
	Contract     *OracleQueue      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleQueueCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleQueueCallerSession struct {
	Contract *OracleQueueCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// OracleQueueTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleQueueTransactorSession struct {
	Contract     *OracleQueueTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// OracleQueueRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleQueueRaw struct {
	Contract *OracleQueue // Generic contract binding to access the raw methods on
}

// OracleQueueCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleQueueCallerRaw struct {
	Contract *OracleQueueCaller // Generic read-only contract binding to access the raw methods on
}

// OracleQueueTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleQueueTransactorRaw struct {
	Contract *OracleQueueTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleQueue creates a new instance of OracleQueue, bound to a specific deployed contract.
func NewOracleQueue(address common.Address, backend bind.ContractBackend) (*OracleQueue, error) {
	contract, err := bindOracleQueue(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleQueue{OracleQueueCaller: OracleQueueCaller{contract: contract}, OracleQueueTransactor: OracleQueueTransactor{contract: contract}, OracleQueueFilterer: OracleQueueFilterer{contract: contract}}, nil
}

// NewOracleQueueCaller creates a new read-only instance of OracleQueue, bound to a specific deployed contract.
func NewOracleQueueCaller(address common.Address, caller bind.ContractCaller) (*OracleQueueCaller, error) {
	contract, err := bindOracleQueue(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleQueueCaller{contract: contract}, nil
}

// NewOracleQueueTransactor creates a new write-only instance of OracleQueue, bound to a specific deployed contract.
func NewOracleQueueTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleQueueTransactor, error) {
	contract, err := bindOracleQueue(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleQueueTransactor{contract: contract}, nil
}

// NewOracleQueueFilterer creates a new log filterer instance of OracleQueue, bound to a specific deployed contract.
func NewOracleQueueFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleQueueFilterer, error) {
	contract, err := bindOracleQueue(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleQueueFilterer{contract: contract}, nil
}

// bindOracleQueue binds a generic wrapper to an already deployed contract.
func bindOracleQueue(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleQueueMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleQueue *OracleQueueRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleQueue.Contract.OracleQueueCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleQueue *OracleQueueRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleQueue.Contract.OracleQueueTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleQueue *OracleQueueRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleQueue.Contract.OracleQueueTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleQueue *OracleQueueCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleQueue.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleQueue *OracleQueueTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleQueue.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleQueue *OracleQueueTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleQueue.Contract.contract.Transact(opts, method, params...)
}

// CustomerQueue is a free data retrieval call binding the contract method 0xc8c6bec6.
//
// Solidity: function customerQueue(uint256 ) view returns(string ipfsCid, address requester, uint256 payment, bool isProcessed)
func (_OracleQueue *OracleQueueCaller) CustomerQueue(opts *bind.CallOpts, arg0 *big.Int) (struct {
	IpfsCid     string
	Requester   common.Address
	Payment     *big.Int
	IsProcessed bool
}, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "customerQueue", arg0)

	outstruct := new(struct {
		IpfsCid     string
		Requester   common.Address
		Payment     *big.Int
		IsProcessed bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IpfsCid = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Requester = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Payment = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.IsProcessed = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// CustomerQueue is a free data retrieval call binding the contract method 0xc8c6bec6.
//
// Solidity: function customerQueue(uint256 ) view returns(string ipfsCid, address requester, uint256 payment, bool isProcessed)
func (_OracleQueue *OracleQueueSession) CustomerQueue(arg0 *big.Int) (struct {
	IpfsCid     string
	Requester   common.Address
	Payment     *big.Int
	IsProcessed bool
}, error) {
	return _OracleQueue.Contract.CustomerQueue(&_OracleQueue.CallOpts, arg0)
}

// CustomerQueue is a free data retrieval call binding the contract method 0xc8c6bec6.
//
// Solidity: function customerQueue(uint256 ) view returns(string ipfsCid, address requester, uint256 payment, bool isProcessed)
func (_OracleQueue *OracleQueueCallerSession) CustomerQueue(arg0 *big.Int) (struct {
	IpfsCid     string
	Requester   common.Address
	Payment     *big.Int
	IsProcessed bool
}, error) {
	return _OracleQueue.Contract.CustomerQueue(&_OracleQueue.CallOpts, arg0)
}

// ModelCreator is a free data retrieval call binding the contract method 0xd24f42be.
//
// Solidity: function modelCreator() view returns(address)
func (_OracleQueue *OracleQueueCaller) ModelCreator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "modelCreator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ModelCreator is a free data retrieval call binding the contract method 0xd24f42be.
//
// Solidity: function modelCreator() view returns(address)
func (_OracleQueue *OracleQueueSession) ModelCreator() (common.Address, error) {
	return _OracleQueue.Contract.ModelCreator(&_OracleQueue.CallOpts)
}

// ModelCreator is a free data retrieval call binding the contract method 0xd24f42be.
//
// Solidity: function modelCreator() view returns(address)
func (_OracleQueue *OracleQueueCallerSession) ModelCreator() (common.Address, error) {
	return _OracleQueue.Contract.ModelCreator(&_OracleQueue.CallOpts)
}

// OracleJobCounter is a free data retrieval call binding the contract method 0xe704095e.
//
// Solidity: function oracleJobCounter() view returns(uint256)
func (_OracleQueue *OracleQueueCaller) OracleJobCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "oracleJobCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OracleJobCounter is a free data retrieval call binding the contract method 0xe704095e.
//
// Solidity: function oracleJobCounter() view returns(uint256)
func (_OracleQueue *OracleQueueSession) OracleJobCounter() (*big.Int, error) {
	return _OracleQueue.Contract.OracleJobCounter(&_OracleQueue.CallOpts)
}

// OracleJobCounter is a free data retrieval call binding the contract method 0xe704095e.
//
// Solidity: function oracleJobCounter() view returns(uint256)
func (_OracleQueue *OracleQueueCallerSession) OracleJobCounter() (*big.Int, error) {
	return _OracleQueue.Contract.OracleJobCounter(&_OracleQueue.CallOpts)
}

// OracleQueue is a free data retrieval call binding the contract method 0xec62b86a.
//
// Solidity: function oracleQueue(uint256 ) view returns(uint256 originalRequestId, string ipfsCid)
func (_OracleQueue *OracleQueueCaller) OracleQueue(opts *bind.CallOpts, arg0 *big.Int) (struct {
	OriginalRequestId *big.Int
	IpfsCid           string
}, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "oracleQueue", arg0)

	outstruct := new(struct {
		OriginalRequestId *big.Int
		IpfsCid           string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OriginalRequestId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IpfsCid = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// OracleQueue is a free data retrieval call binding the contract method 0xec62b86a.
//
// Solidity: function oracleQueue(uint256 ) view returns(uint256 originalRequestId, string ipfsCid)
func (_OracleQueue *OracleQueueSession) OracleQueue(arg0 *big.Int) (struct {
	OriginalRequestId *big.Int
	IpfsCid           string
}, error) {
	return _OracleQueue.Contract.OracleQueue(&_OracleQueue.CallOpts, arg0)
}

// OracleQueue is a free data retrieval call binding the contract method 0xec62b86a.
//
// Solidity: function oracleQueue(uint256 ) view returns(uint256 originalRequestId, string ipfsCid)
func (_OracleQueue *OracleQueueCallerSession) OracleQueue(arg0 *big.Int) (struct {
	OriginalRequestId *big.Int
	IpfsCid           string
}, error) {
	return _OracleQueue.Contract.OracleQueue(&_OracleQueue.CallOpts, arg0)
}

// OracleReward is a free data retrieval call binding the contract method 0x21873631.
//
// Solidity: function oracleReward() view returns(uint256)
func (_OracleQueue *OracleQueueCaller) OracleReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "oracleReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OracleReward is a free data retrieval call binding the contract method 0x21873631.
//
// Solidity: function oracleReward() view returns(uint256)
func (_OracleQueue *OracleQueueSession) OracleReward() (*big.Int, error) {
	return _OracleQueue.Contract.OracleReward(&_OracleQueue.CallOpts)
}

// OracleReward is a free data retrieval call binding the contract method 0x21873631.
//
// Solidity: function oracleReward() view returns(uint256)
func (_OracleQueue *OracleQueueCallerSession) OracleReward() (*big.Int, error) {
	return _OracleQueue.Contract.OracleReward(&_OracleQueue.CallOpts)
}

// OracleVerifierAddress is a free data retrieval call binding the contract method 0x42e7c0b0.
//
// Solidity: function oracleVerifierAddress() view returns(address)
func (_OracleQueue *OracleQueueCaller) OracleVerifierAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "oracleVerifierAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleVerifierAddress is a free data retrieval call binding the contract method 0x42e7c0b0.
//
// Solidity: function oracleVerifierAddress() view returns(address)
func (_OracleQueue *OracleQueueSession) OracleVerifierAddress() (common.Address, error) {
	return _OracleQueue.Contract.OracleVerifierAddress(&_OracleQueue.CallOpts)
}

// OracleVerifierAddress is a free data retrieval call binding the contract method 0x42e7c0b0.
//
// Solidity: function oracleVerifierAddress() view returns(address)
func (_OracleQueue *OracleQueueCallerSession) OracleVerifierAddress() (common.Address, error) {
	return _OracleQueue.Contract.OracleVerifierAddress(&_OracleQueue.CallOpts)
}

// QueryFee is a free data retrieval call binding the contract method 0xfdd26881.
//
// Solidity: function queryFee() view returns(uint256)
func (_OracleQueue *OracleQueueCaller) QueryFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "queryFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QueryFee is a free data retrieval call binding the contract method 0xfdd26881.
//
// Solidity: function queryFee() view returns(uint256)
func (_OracleQueue *OracleQueueSession) QueryFee() (*big.Int, error) {
	return _OracleQueue.Contract.QueryFee(&_OracleQueue.CallOpts)
}

// QueryFee is a free data retrieval call binding the contract method 0xfdd26881.
//
// Solidity: function queryFee() view returns(uint256)
func (_OracleQueue *OracleQueueCallerSession) QueryFee() (*big.Int, error) {
	return _OracleQueue.Contract.QueryFee(&_OracleQueue.CallOpts)
}

// RequestCounter is a free data retrieval call binding the contract method 0x973a814e.
//
// Solidity: function requestCounter() view returns(uint256)
func (_OracleQueue *OracleQueueCaller) RequestCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OracleQueue.contract.Call(opts, &out, "requestCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestCounter is a free data retrieval call binding the contract method 0x973a814e.
//
// Solidity: function requestCounter() view returns(uint256)
func (_OracleQueue *OracleQueueSession) RequestCounter() (*big.Int, error) {
	return _OracleQueue.Contract.RequestCounter(&_OracleQueue.CallOpts)
}

// RequestCounter is a free data retrieval call binding the contract method 0x973a814e.
//
// Solidity: function requestCounter() view returns(uint256)
func (_OracleQueue *OracleQueueCallerSession) RequestCounter() (*big.Int, error) {
	return _OracleQueue.Contract.RequestCounter(&_OracleQueue.CallOpts)
}

// ApproveJob is a paid mutator transaction binding the contract method 0x4bd23b9e.
//
// Solidity: function approveJob(uint256 _requestId) returns()
func (_OracleQueue *OracleQueueTransactor) ApproveJob(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _OracleQueue.contract.Transact(opts, "approveJob", _requestId)
}

// ApproveJob is a paid mutator transaction binding the contract method 0x4bd23b9e.
//
// Solidity: function approveJob(uint256 _requestId) returns()
func (_OracleQueue *OracleQueueSession) ApproveJob(_requestId *big.Int) (*types.Transaction, error) {
	return _OracleQueue.Contract.ApproveJob(&_OracleQueue.TransactOpts, _requestId)
}

// ApproveJob is a paid mutator transaction binding the contract method 0x4bd23b9e.
//
// Solidity: function approveJob(uint256 _requestId) returns()
func (_OracleQueue *OracleQueueTransactorSession) ApproveJob(_requestId *big.Int) (*types.Transaction, error) {
	return _OracleQueue.Contract.ApproveJob(&_OracleQueue.TransactOpts, _requestId)
}

// RequestAttribution is a paid mutator transaction binding the contract method 0x853edc03.
//
// Solidity: function requestAttribution(string _ipfsCid) payable returns()
func (_OracleQueue *OracleQueueTransactor) RequestAttribution(opts *bind.TransactOpts, _ipfsCid string) (*types.Transaction, error) {
	return _OracleQueue.contract.Transact(opts, "requestAttribution", _ipfsCid)
}

// RequestAttribution is a paid mutator transaction binding the contract method 0x853edc03.
//
// Solidity: function requestAttribution(string _ipfsCid) payable returns()
func (_OracleQueue *OracleQueueSession) RequestAttribution(_ipfsCid string) (*types.Transaction, error) {
	return _OracleQueue.Contract.RequestAttribution(&_OracleQueue.TransactOpts, _ipfsCid)
}

// RequestAttribution is a paid mutator transaction binding the contract method 0x853edc03.
//
// Solidity: function requestAttribution(string _ipfsCid) payable returns()
func (_OracleQueue *OracleQueueTransactorSession) RequestAttribution(_ipfsCid string) (*types.Transaction, error) {
	return _OracleQueue.Contract.RequestAttribution(&_OracleQueue.TransactOpts, _ipfsCid)
}

// RewardOracle is a paid mutator transaction binding the contract method 0x1fafb331.
//
// Solidity: function rewardOracle(address _oracle) returns()
func (_OracleQueue *OracleQueueTransactor) RewardOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _OracleQueue.contract.Transact(opts, "rewardOracle", _oracle)
}

// RewardOracle is a paid mutator transaction binding the contract method 0x1fafb331.
//
// Solidity: function rewardOracle(address _oracle) returns()
func (_OracleQueue *OracleQueueSession) RewardOracle(_oracle common.Address) (*types.Transaction, error) {
	return _OracleQueue.Contract.RewardOracle(&_OracleQueue.TransactOpts, _oracle)
}

// RewardOracle is a paid mutator transaction binding the contract method 0x1fafb331.
//
// Solidity: function rewardOracle(address _oracle) returns()
func (_OracleQueue *OracleQueueTransactorSession) RewardOracle(_oracle common.Address) (*types.Transaction, error) {
	return _OracleQueue.Contract.RewardOracle(&_OracleQueue.TransactOpts, _oracle)
}

// SetVerifierAddress is a paid mutator transaction binding the contract method 0x17e95526.
//
// Solidity: function setVerifierAddress(address _verifier) returns()
func (_OracleQueue *OracleQueueTransactor) SetVerifierAddress(opts *bind.TransactOpts, _verifier common.Address) (*types.Transaction, error) {
	return _OracleQueue.contract.Transact(opts, "setVerifierAddress", _verifier)
}

// SetVerifierAddress is a paid mutator transaction binding the contract method 0x17e95526.
//
// Solidity: function setVerifierAddress(address _verifier) returns()
func (_OracleQueue *OracleQueueSession) SetVerifierAddress(_verifier common.Address) (*types.Transaction, error) {
	return _OracleQueue.Contract.SetVerifierAddress(&_OracleQueue.TransactOpts, _verifier)
}

// SetVerifierAddress is a paid mutator transaction binding the contract method 0x17e95526.
//
// Solidity: function setVerifierAddress(address _verifier) returns()
func (_OracleQueue *OracleQueueTransactorSession) SetVerifierAddress(_verifier common.Address) (*types.Transaction, error) {
	return _OracleQueue.Contract.SetVerifierAddress(&_OracleQueue.TransactOpts, _verifier)
}

// OracleQueueLogNewCustomerRequestIterator is returned from FilterLogNewCustomerRequest and is used to iterate over the raw logs and unpacked data for LogNewCustomerRequest events raised by the OracleQueue contract.
type OracleQueueLogNewCustomerRequestIterator struct {
	Event *OracleQueueLogNewCustomerRequest // Event containing the contract specifics and raw log

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
func (it *OracleQueueLogNewCustomerRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleQueueLogNewCustomerRequest)
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
		it.Event = new(OracleQueueLogNewCustomerRequest)
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
func (it *OracleQueueLogNewCustomerRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleQueueLogNewCustomerRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleQueueLogNewCustomerRequest represents a LogNewCustomerRequest event raised by the OracleQueue contract.
type OracleQueueLogNewCustomerRequest struct {
	RequestId *big.Int
	IpfsCid   string
	Requester common.Address
	Payment   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogNewCustomerRequest is a free log retrieval operation binding the contract event 0x4498100bf5e1c1ca57375bd8b4f12c7cf4b4459bbe7c9ae6b6d55c44987bf001.
//
// Solidity: event LogNewCustomerRequest(uint256 indexed requestId, string ipfsCid, address requester, uint256 payment)
func (_OracleQueue *OracleQueueFilterer) FilterLogNewCustomerRequest(opts *bind.FilterOpts, requestId []*big.Int) (*OracleQueueLogNewCustomerRequestIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OracleQueue.contract.FilterLogs(opts, "LogNewCustomerRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleQueueLogNewCustomerRequestIterator{contract: _OracleQueue.contract, event: "LogNewCustomerRequest", logs: logs, sub: sub}, nil
}

// WatchLogNewCustomerRequest is a free log subscription operation binding the contract event 0x4498100bf5e1c1ca57375bd8b4f12c7cf4b4459bbe7c9ae6b6d55c44987bf001.
//
// Solidity: event LogNewCustomerRequest(uint256 indexed requestId, string ipfsCid, address requester, uint256 payment)
func (_OracleQueue *OracleQueueFilterer) WatchLogNewCustomerRequest(opts *bind.WatchOpts, sink chan<- *OracleQueueLogNewCustomerRequest, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _OracleQueue.contract.WatchLogs(opts, "LogNewCustomerRequest", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleQueueLogNewCustomerRequest)
				if err := _OracleQueue.contract.UnpackLog(event, "LogNewCustomerRequest", log); err != nil {
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

// ParseLogNewCustomerRequest is a log parse operation binding the contract event 0x4498100bf5e1c1ca57375bd8b4f12c7cf4b4459bbe7c9ae6b6d55c44987bf001.
//
// Solidity: event LogNewCustomerRequest(uint256 indexed requestId, string ipfsCid, address requester, uint256 payment)
func (_OracleQueue *OracleQueueFilterer) ParseLogNewCustomerRequest(log types.Log) (*OracleQueueLogNewCustomerRequest, error) {
	event := new(OracleQueueLogNewCustomerRequest)
	if err := _OracleQueue.contract.UnpackLog(event, "LogNewCustomerRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleQueueLogNewJobForOraclesIterator is returned from FilterLogNewJobForOracles and is used to iterate over the raw logs and unpacked data for LogNewJobForOracles events raised by the OracleQueue contract.
type OracleQueueLogNewJobForOraclesIterator struct {
	Event *OracleQueueLogNewJobForOracles // Event containing the contract specifics and raw log

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
func (it *OracleQueueLogNewJobForOraclesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleQueueLogNewJobForOracles)
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
		it.Event = new(OracleQueueLogNewJobForOracles)
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
func (it *OracleQueueLogNewJobForOraclesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleQueueLogNewJobForOraclesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleQueueLogNewJobForOracles represents a LogNewJobForOracles event raised by the OracleQueue contract.
type OracleQueueLogNewJobForOracles struct {
	JobId   *big.Int
	IpfsCid string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLogNewJobForOracles is a free log retrieval operation binding the contract event 0x5af5066ee921684d5e820a56d6e5abdfb117e07345b9e29718782ba6251f68a5.
//
// Solidity: event LogNewJobForOracles(uint256 indexed jobId, string ipfsCid)
func (_OracleQueue *OracleQueueFilterer) FilterLogNewJobForOracles(opts *bind.FilterOpts, jobId []*big.Int) (*OracleQueueLogNewJobForOraclesIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _OracleQueue.contract.FilterLogs(opts, "LogNewJobForOracles", jobIdRule)
	if err != nil {
		return nil, err
	}
	return &OracleQueueLogNewJobForOraclesIterator{contract: _OracleQueue.contract, event: "LogNewJobForOracles", logs: logs, sub: sub}, nil
}

// WatchLogNewJobForOracles is a free log subscription operation binding the contract event 0x5af5066ee921684d5e820a56d6e5abdfb117e07345b9e29718782ba6251f68a5.
//
// Solidity: event LogNewJobForOracles(uint256 indexed jobId, string ipfsCid)
func (_OracleQueue *OracleQueueFilterer) WatchLogNewJobForOracles(opts *bind.WatchOpts, sink chan<- *OracleQueueLogNewJobForOracles, jobId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _OracleQueue.contract.WatchLogs(opts, "LogNewJobForOracles", jobIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleQueueLogNewJobForOracles)
				if err := _OracleQueue.contract.UnpackLog(event, "LogNewJobForOracles", log); err != nil {
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

// ParseLogNewJobForOracles is a log parse operation binding the contract event 0x5af5066ee921684d5e820a56d6e5abdfb117e07345b9e29718782ba6251f68a5.
//
// Solidity: event LogNewJobForOracles(uint256 indexed jobId, string ipfsCid)
func (_OracleQueue *OracleQueueFilterer) ParseLogNewJobForOracles(log types.Log) (*OracleQueueLogNewJobForOracles, error) {
	event := new(OracleQueueLogNewJobForOracles)
	if err := _OracleQueue.contract.UnpackLog(event, "LogNewJobForOracles", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
