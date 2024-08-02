package vechain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	log "github.com/sirupsen/logrus"
	"github.com/vechain/thor/api/transactions"
	"github.com/vechain/thor/tx"
	"io"
	"math/big"
	"net/http"
)

type Client struct {
	ctx     context.Context
	url     string
	httpcli *http.Client
}

func NewVechainClient(ctx context.Context, url string) *Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	hcli := new(http.Client)
	hcli.Transport = tr
	return &Client{
		ctx:     ctx,
		url:     url,
		httpcli: hcli,
	}
}

type BlockInfo struct {
	Number      int64    `json:"number"`
	Txs         []string `json:"transactions"`
	Timestamp   int64    `json:"timestamp"`
	ID          string   `json:"id"`
	Beneficiary string   `json:"beneficiary"`
}

func (c *Client) Close() {
	c.httpcli.CloseIdleConnections()
}

func (c *Client) ChainTag() (byte, error) {
	api := fmt.Sprintf("%s/blocks/0", c.url)
	res, err := c.get(api)
	if err != nil {
		return 0x00, err
	}
	var genInfo BlockInfo
	if err = json.Unmarshal(res, &genInfo); err != nil {
		log.Printf("json.Unmarshal genesis info: %v", err)
		return 0x00, err
	}
	id := common.FromHex(genInfo.ID)
	return id[31], nil
}

func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (BlockInfo, error) {
	api := fmt.Sprintf("%s/blocks/%d", c.url, number.Int64())
	res, err := c.get(api)
	if err != nil {
		return BlockInfo{}, err
	}
	var blk BlockInfo
	if err = json.Unmarshal(res, &blk); err != nil {
		log.Printf("json.Unmarshal block info: %v", err)
		return BlockInfo{}, err
	}
	return blk, nil
}

func (c *Client) FinalizedBlock(ctx context.Context) (BlockInfo, error) {
	api := fmt.Sprintf("%s/blocks/finalized", c.url)
	res, err := c.get(api)
	if err != nil {
		return BlockInfo{}, err
	}
	var blk BlockInfo
	if err = json.Unmarshal(res, &blk); err != nil {
		log.Printf("json.Unmarshal finalized block info: %v", err)
		return BlockInfo{}, err
	}
	return blk, nil
}

func (c *Client) LatestBlock(ctx context.Context) (BlockInfo, error) {
	api := fmt.Sprintf("%s/blocks/best", c.url)
	res, err := c.get(api)
	if err != nil {
		return BlockInfo{}, err
	}
	var blk BlockInfo
	if err = json.Unmarshal(res, &blk); err != nil {
		log.Printf("json.Unmarshal latest block info: %v", err)
		return BlockInfo{}, err
	}
	return blk, nil
}

func (c *Client) post(url string, obj interface{}) ([]byte, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.WithError(err).Error("post failed json.Marshal")
		return []byte{}, err
	}
	res, err := c.httpcli.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		log.WithError(err).Error("post failed http.Post")
		return nil, err
	}
	r, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("post failed io.ReadAll")
		return nil, err
	}
	defer res.Body.Close()
	return r, nil
}

func (c *Client) get(url string) ([]byte, error) {
	res, err := c.httpcli.Get(url)
	if err != nil {
		log.WithError(err).Error("get failed http.Get")
		return nil, err
	}
	r, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Error("get failed io.ReadAll")
		return nil, err

	}
	defer res.Body.Close()
	return r, nil
}

func (c *Client) SendTransaction(ctx context.Context, tx *tx.Transaction) error {
	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.WithError(err).Error("rlp.EncodeToBytes")
		return err
	}
	_, err = c.post(c.url+"/transactions", transactions.RawTx{Raw: hexutil.Encode(rlpTx)})
	if err != nil {
		log.WithError(err).Error("send transaction failed")
		return err
	}
	return nil
}

type DummyTxReceipt struct {
	Meta transactions.ReceiptMeta `json:"meta"`
}

func (c *Client) TransactionReceipt(ctx context.Context, txid string) (DummyTxReceipt, error) {
	api := fmt.Sprintf("%s/transactions/%s/receipt", c.url, txid)
	res, err := c.get(api)
	if err != nil {
		return DummyTxReceipt{}, err
	}
	var receipt DummyTxReceipt
	if err = json.Unmarshal(res, &receipt); err != nil {
		log.Printf("json.Unmarshal receipt: %v", err)
		return DummyTxReceipt{}, err
	}
	return receipt, nil
}
