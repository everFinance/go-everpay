package sdk

import (
	"github.com/everFinance/goether"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestSDK_Transfer(t *testing.T) {
	// signer, err :=goether.NewSigner("ad1dcf8f1c449e7af21a7b8341eba5f053055819fff9948f1251ea94a0184cae")
	// assert.NoError(t, err)
	// testSDK , err := New(signer , "https://api-dev.everpay.io")
	// assert.NoError(t, err)
	// to := "0xa2026731B31E4DFBa78314bDBfBFDC8cF5F761F8"
	// amount := big.NewInt(2000000)
	// result, err := testSDK.Transfer("usdt", amount, to,`{"msg": "hello"}`)
	// assert.NoError(t, err)
	// t.Log(result.HexHash())
}

func TestClient_SubmitBundleTx(t *testing.T) {
	// addr01 := "0x3D7e9DFbc58952FdACEe2a5C69367C8478474D82"
	// priv01 := "ad1dcf8f1c449e7af21a7b8341eba5f053055819fff9948f1251ea94a0184cae"
	// priv02 := "338f76e7463ed64f98e883aa0f522c92cc5881cbce113894559d703d515a55e1"
	// addr02 := "0xf392A4e8DDbfBD7782407561B8Beab911c36d59A"
	//
	// addr03 := "cSYOy8-p1QFenktkDBFyRM3cwZSTrQ_J4EsELLho_UE"
	//
	// items := []paySchema.BundleItem{
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr02,
	// 		Amount:  "99999",
	// 	},
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr03,
	// 		Amount:  "888888",
	// 	},
	// 	{
	// 		Tag:     "ethereum-usdt-0xd85476c906b5301e8e9eb58d174a6f96b9dfc5ee",
	// 		ChainID: "42",
	// 		From:    addr01,
	// 		To:      addr03,
	// 		Amount:  "12345",
	// 	},
	// 	{
	// 		Tag:     "ethereum-eth-0x0000000000000000000000000000000000000000",
	// 		ChainID: "42",
	// 		From:    addr02,
	// 		To:      addr03,
	// 		Amount:  "6666",
	// 	},
	// }
	//
	// txNonce := time.Now().UnixNano() / 1e6
	// expiration :=  txNonce/1000 + 1000
	// bundle := GenBundle(items, expiration)
	//
	// signer01 , _ := goether.NewSigner(priv01)
	// signer02 , _ := goether.NewSigner(priv02)
	// sdk01 ,err := New(signer01,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	// sdk02 ,err:= New(signer02,"https://api-dev.everpay.io")
	// assert.NoError(t, err)
	//
	// bundleData01, err := sdk01.SignBundleData(bundle)
	// assert.NoError(t, err)
	// bundleData02, err := sdk02.SignBundleData(bundle)
	// assert.NoError(t, err)
	//
	// bundleSigs := paySchema.BundleWithSigs{
	// 	Bundle: bundle,
	// 	Sigs: map[string]string{
	// 		sdk01.AccId: bundleData01.Sigs[sdk01.AccId],
	// 		sdk02.AccId: bundleData02.Sigs[sdk02.AccId],
	// 	},
	// }
	//
	// res, err := sdk01.Bundle("ETH",addr01,nil,bundleSigs)
	// assert.NoError(t, err)
	// t.Log(res.HexHash())
}

func TestBurnTx(t *testing.T) {
	payUrl := "https://api-dev.everpay.io"
	signer, err := goether.NewSigner("338f76e7463ed64f98e883aa0f522c92cc5881cbce113894559d703d515a55e1")
	if err != nil {
		panic(err)
	}
	t.Log(signer.Address)
	sdk, err := New(signer, payUrl)
	assert.NoError(t, err)

	tokenTag := "bsc-ddd-0xb5eadfdbdb40257d1d24a1432faa2503a867c270"
	targetChain := "moon"
	to := "0x63d5bdaba94a5fccb8980bdc738e017df5f9a1fd"
	amount := big.NewInt(6000000000000000000)
	res, err := sdk.Withdraw(tokenTag, amount, targetChain, to)
	assert.NoError(t, err)
	t.Log(res.HexHash())
}
