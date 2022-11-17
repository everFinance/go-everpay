package sdk

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	paySchema "github.com/everFinance/go-everpay/pay/schema"
	"github.com/everFinance/go-everpay/sdk/schema"
	tokSchema "github.com/everFinance/go-everpay/token/schema"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/everFinance/go-everpay/common"
)

var log = common.NewLog("sdk")

type SDK struct {
	tokens       map[string]schema.TokenInfo // tag -> TokenInfo
	feeRecipient string

	signerType string // ecc, rsa
	signer     interface{}

	AccId string
	Cli   *Client

	lastNonce    int64 // last everTx used nonce
	sendTxLocker sync.Mutex
}

func New(signer interface{}, payUrl string) (*SDK, error) {
	signerType, signerAddr, err := reflectSigner(signer)
	if err != nil {
		return nil, err
	}

	sdk := &SDK{
		signer:       signer,
		signerType:   signerType,
		AccId:        signerAddr,
		Cli:          NewClient(payUrl),
		lastNonce:    time.Now().UnixNano() / 1000000,
		sendTxLocker: sync.Mutex{},
	}
	_ = sdk.updatePayInfo()

	// sync info from everPay server every 10 mintue
	go sdk.runSyncInfo()
	return sdk, nil
}

func (s *SDK) runSyncInfo() {
	for {
		err := s.updatePayInfo()
		if err != nil {
			log.Error("can not get info from everpay", "err", err)
			time.Sleep(5 * time.Second)
			continue
		}
		time.Sleep(10 * time.Minute)
	}
}

func (s *SDK) updatePayInfo() error {
	info, err := s.Cli.GetInfo()
	if err != nil {
		return err
	}

	tokens := make(map[string]schema.TokenInfo)
	for _, t := range info.TokenList {
		tokens[t.Tag] = t
	}
	s.tokens = tokens
	s.feeRecipient = info.FeeRecipient
	return nil
}

func (s *SDK) GetTokens() map[string]schema.TokenInfo {
	return s.tokens
}

func (s *SDK) Transfer(tokenTag string, amount *big.Int, to, data string) (*paySchema.Transaction, error) {
	return s.sendTransfer(tokenTag, to, amount, data)
}

func (s *SDK) Withdraw(tokenTag string, amount *big.Int, chainType, to string) (*paySchema.Transaction, error) {
	return s.sendWithdraw(tokenTag, chainType, to, amount, "")
}

func (s *SDK) Bundle(tokenTag string, to string, amount *big.Int, bundleWithSigs paySchema.BundleWithSigs) (*paySchema.Transaction, error) {
	bundle := paySchema.BundleData{
		Bundle: bundleWithSigs,
	}
	return s.sendBundle(tokenTag, to, amount, bundle)
}

func (s *SDK) sendTransfer(tokenTag string, receiver string, amount *big.Int, data string) (*paySchema.Transaction, error) {
	tokenInfo, ok := s.tokens[tokenTag]
	if !ok {
		return nil, ErrTokenNotExist
	}
	action := tokSchema.TxActionTransfer
	fee := tokenInfo.TransferFee
	return s.sendTx(tokenInfo, action, fee, receiver, amount, data)
}

func (s *SDK) sendWithdraw(tokenTag string, targetChainType, receiver string, amount *big.Int, data string) (*paySchema.Transaction, error) {
	tokenInfo, ok := s.tokens[tokenTag]
	if !ok {
		return nil, ErrTokenNotExist
	}
	action := tokSchema.TxActionBurn
	tFee, err := s.Cli.Fee(tokenTag)
	if err != nil {
		return nil, err
	}
	fee, ok := tFee.Fee.BurnFeeMap[targetChainType]
	if !ok {
		return nil, ErrBurnFeeNotExist
	}
	if data != "" && !gjson.Valid(data) {
		return nil, ErrNotJsonData
	}

	// add targetChainType in data
	txData, err := sjson.Set(data, "targetChainType", targetChainType)
	if err != nil {
		return nil, err
	}
	return s.sendTx(tokenInfo, action, fee, receiver, amount, txData)
}

func (s *SDK) sendBundle(tokenTag string, receiver string, amount *big.Int, bundle paySchema.BundleData) (*paySchema.Transaction, error) {
	tokenInfo, ok := s.tokens[tokenTag]
	if !ok {
		return nil, ErrTokenNotExist
	}
	action := paySchema.TxActionBundle
	fee := tokenInfo.BundleFee

	data, err := json.Marshal(bundle)
	if err != nil {
		return nil, err
	}

	return s.sendTx(tokenInfo, action, fee, receiver, amount, string(data))
}

func (s *SDK) sendTx(tokenInfo schema.TokenInfo, action, fee, receiver string, amount *big.Int, data string) (*paySchema.Transaction, error) {
	s.sendTxLocker.Lock()
	defer s.sendTxLocker.Unlock()
	if amount == nil {
		amount = big.NewInt(0)
	}
	// assemble tx
	everTx := paySchema.Transaction{
		TokenSymbol:  tokenInfo.Symbol,
		Action:       strings.ToLower(action),
		From:         s.AccId,
		To:           receiver,
		Amount:       amount.String(),
		Fee:          fee,
		FeeRecipient: s.feeRecipient,
		Nonce:        fmt.Sprintf("%d", s.getNonce()),
		TokenID:      tokenInfo.ID,
		ChainType:    tokenInfo.ChainType,
		ChainID:      tokenInfo.ChainID,
		Data:         data,
		Version:      tokSchema.TxVersionV1,
		Sig:          "",
	}

	sign, err := s.Sign(everTx.String())
	if err != nil {
		log.Error("Sign failed", "error", err)
		return &everTx, err
	}
	everTx.Sig = sign

	// submit to everpay server
	if err := s.Cli.SubmitTx(everTx); err != nil {
		log.Error("submit everTx", "error", err)
		return &everTx, err
	}

	return &everTx, nil
}

// about bundleTx

// GenBundle expiration: bundle tx expiration time(s)
func GenBundle(items []paySchema.BundleItem, expiration int64) paySchema.Bundle {
	return paySchema.Bundle{
		Items:      items,
		Expiration: expiration,
		Salt:       uuid.NewString(),
		Version:    paySchema.BundleTxVersionV1,
	}
}

func (s *SDK) SignBundleData(bundleTx paySchema.Bundle) (paySchema.BundleWithSigs, error) {
	sign, err := s.Sign(bundleTx.String())
	if err != nil {
		return paySchema.BundleWithSigs{}, err
	}
	return paySchema.BundleWithSigs{
		Bundle: bundleTx,
		Sigs: map[string]string{
			s.AccId: sign,
		},
	}, nil
}

func (s *SDK) getNonce() int64 {
	for {
		newNonce := time.Now().UnixNano() / 1000000
		if newNonce > s.lastNonce {
			s.lastNonce = newNonce
			return newNonce
		}
	}
}
