### build binary
```
make
```

### generate test account
create 10000 test account with eth coin.
```
./txpress account --count 1000 --eth 10000
```

then generate accounts and private key in `accounts.json`, and genesis inifo in `balance.json`.


### do test
modify app.json with your case.

- **tx_count：** total transaction count to generate and send.
- **send_routine_count：** routine count to concurrent send tx.
- **speed:** send tx count per second.
- **erc20_contract:** default used erc20 token contract.
- **rpc_node:** rpc url to the chain.
- **receive_addr：** receive address for all transaction, if empty, will create a new account to receive.
- **amount：** the amount of transfer.
- **chain_id：** chain id. 
- **batch_transfer_contract：** contract used to batch send eth/token when init test accounts.

test transfer eth coin.
```
./txpress --start
```
