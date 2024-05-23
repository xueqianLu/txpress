// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// HandOutPoolMetaData contains all meta data concerning the HandOutPool contract.
var HandOutPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"v\",\"type\":\"uint256\"}],\"name\":\"Handout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"v\",\"type\":\"uint256\"}],\"name\":\"HandoutToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"lst\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"b\",\"type\":\"bool\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// HandOutPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use HandOutPoolMetaData.ABI instead.
var HandOutPoolABI = HandOutPoolMetaData.ABI

// HandOutPool is an auto generated Go binding around an Ethereum contract.
type HandOutPool struct {
	HandOutPoolCaller     // Read-only binding to the contract
	HandOutPoolTransactor // Write-only binding to the contract
	HandOutPoolFilterer   // Log filterer for contract events
}

// HandOutPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type HandOutPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandOutPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HandOutPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandOutPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HandOutPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HandOutPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HandOutPoolSession struct {
	Contract     *HandOutPool      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HandOutPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HandOutPoolCallerSession struct {
	Contract *HandOutPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// HandOutPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HandOutPoolTransactorSession struct {
	Contract     *HandOutPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// HandOutPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type HandOutPoolRaw struct {
	Contract *HandOutPool // Generic contract binding to access the raw methods on
}

// HandOutPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HandOutPoolCallerRaw struct {
	Contract *HandOutPoolCaller // Generic read-only contract binding to access the raw methods on
}

// HandOutPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HandOutPoolTransactorRaw struct {
	Contract *HandOutPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHandOutPool creates a new instance of HandOutPool, bound to a specific deployed contract.
func NewHandOutPool(address common.Address, backend bind.ContractBackend) (*HandOutPool, error) {
	contract, err := bindHandOutPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HandOutPool{HandOutPoolCaller: HandOutPoolCaller{contract: contract}, HandOutPoolTransactor: HandOutPoolTransactor{contract: contract}, HandOutPoolFilterer: HandOutPoolFilterer{contract: contract}}, nil
}

// NewHandOutPoolCaller creates a new read-only instance of HandOutPool, bound to a specific deployed contract.
func NewHandOutPoolCaller(address common.Address, caller bind.ContractCaller) (*HandOutPoolCaller, error) {
	contract, err := bindHandOutPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HandOutPoolCaller{contract: contract}, nil
}

// NewHandOutPoolTransactor creates a new write-only instance of HandOutPool, bound to a specific deployed contract.
func NewHandOutPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*HandOutPoolTransactor, error) {
	contract, err := bindHandOutPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HandOutPoolTransactor{contract: contract}, nil
}

// NewHandOutPoolFilterer creates a new log filterer instance of HandOutPool, bound to a specific deployed contract.
func NewHandOutPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*HandOutPoolFilterer, error) {
	contract, err := bindHandOutPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HandOutPoolFilterer{contract: contract}, nil
}

// bindHandOutPool binds a generic wrapper to an already deployed contract.
func bindHandOutPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HandOutPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HandOutPool *HandOutPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HandOutPool.Contract.HandOutPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HandOutPool *HandOutPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HandOutPool.Contract.HandOutPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HandOutPool *HandOutPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HandOutPool.Contract.HandOutPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HandOutPool *HandOutPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HandOutPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HandOutPool *HandOutPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HandOutPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HandOutPool *HandOutPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HandOutPool.Contract.contract.Transact(opts, method, params...)
}

// Handout is a paid mutator transaction binding the contract method 0x676d0b79.
//
// Solidity: function Handout(address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolTransactor) Handout(opts *bind.TransactOpts, to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.contract.Transact(opts, "Handout", to, v)
}

// Handout is a paid mutator transaction binding the contract method 0x676d0b79.
//
// Solidity: function Handout(address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolSession) Handout(to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.Contract.Handout(&_HandOutPool.TransactOpts, to, v)
}

// Handout is a paid mutator transaction binding the contract method 0x676d0b79.
//
// Solidity: function Handout(address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolTransactorSession) Handout(to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.Contract.Handout(&_HandOutPool.TransactOpts, to, v)
}

// HandoutToken is a paid mutator transaction binding the contract method 0x2cc8ae57.
//
// Solidity: function HandoutToken(address token, address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolTransactor) HandoutToken(opts *bind.TransactOpts, token common.Address, to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.contract.Transact(opts, "HandoutToken", token, to, v)
}

// HandoutToken is a paid mutator transaction binding the contract method 0x2cc8ae57.
//
// Solidity: function HandoutToken(address token, address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolSession) HandoutToken(token common.Address, to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.Contract.HandoutToken(&_HandOutPool.TransactOpts, token, to, v)
}

// HandoutToken is a paid mutator transaction binding the contract method 0x2cc8ae57.
//
// Solidity: function HandoutToken(address token, address[] to, uint256 v) returns()
func (_HandOutPool *HandOutPoolTransactorSession) HandoutToken(token common.Address, to []common.Address, v *big.Int) (*types.Transaction, error) {
	return _HandOutPool.Contract.HandoutToken(&_HandOutPool.TransactOpts, token, to, v)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x84e7686b.
//
// Solidity: function setAdmin(address[] lst, bool b) returns()
func (_HandOutPool *HandOutPoolTransactor) SetAdmin(opts *bind.TransactOpts, lst []common.Address, b bool) (*types.Transaction, error) {
	return _HandOutPool.contract.Transact(opts, "setAdmin", lst, b)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x84e7686b.
//
// Solidity: function setAdmin(address[] lst, bool b) returns()
func (_HandOutPool *HandOutPoolSession) SetAdmin(lst []common.Address, b bool) (*types.Transaction, error) {
	return _HandOutPool.Contract.SetAdmin(&_HandOutPool.TransactOpts, lst, b)
}

// SetAdmin is a paid mutator transaction binding the contract method 0x84e7686b.
//
// Solidity: function setAdmin(address[] lst, bool b) returns()
func (_HandOutPool *HandOutPoolTransactorSession) SetAdmin(lst []common.Address, b bool) (*types.Transaction, error) {
	return _HandOutPool.Contract.SetAdmin(&_HandOutPool.TransactOpts, lst, b)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_HandOutPool *HandOutPoolTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HandOutPool.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_HandOutPool *HandOutPoolSession) Receive() (*types.Transaction, error) {
	return _HandOutPool.Contract.Receive(&_HandOutPool.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_HandOutPool *HandOutPoolTransactorSession) Receive() (*types.Transaction, error) {
	return _HandOutPool.Contract.Receive(&_HandOutPool.TransactOpts)
}

// HandOutPoolTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the HandOutPool contract.
type HandOutPoolTransferIterator struct {
	Event *HandOutPoolTransfer // Event containing the contract specifics and raw log

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
func (it *HandOutPoolTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HandOutPoolTransfer)
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
		it.Event = new(HandOutPoolTransfer)
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
func (it *HandOutPoolTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HandOutPoolTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HandOutPoolTransfer represents a Transfer event raised by the HandOutPool contract.
type HandOutPoolTransfer struct {
	Src   common.Address
	Dst   common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 value)
func (_HandOutPool *HandOutPoolFilterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*HandOutPoolTransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _HandOutPool.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &HandOutPoolTransferIterator{contract: _HandOutPool.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 value)
func (_HandOutPool *HandOutPoolFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *HandOutPoolTransfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _HandOutPool.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HandOutPoolTransfer)
				if err := _HandOutPool.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 value)
func (_HandOutPool *HandOutPoolFilterer) ParseTransfer(log types.Log) (*HandOutPoolTransfer, error) {
	event := new(HandOutPoolTransfer)
	if err := _HandOutPool.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
