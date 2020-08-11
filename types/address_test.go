package types_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/pocblockchain/pocc/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"gopkg.in/yaml.v2"
)

var invalidStrs = []string{
	"hello, world!",
	"0xAA",
	"AAA",
	types.Bech32PrefixAccAddr + "AB0C",
	types.Bech32PrefixAccPub + "1234",
	types.Bech32PrefixValAddr + "5678",
	types.Bech32PrefixValPub + "BBAB",
	types.Bech32PrefixConsAddr + "FF04",
	types.Bech32PrefixConsPub + "6789",
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

func TestEmptyAddresses(t *testing.T) {
	require.Equal(t, (types.AccAddress{}).String(), "")
	require.Equal(t, (types.ValAddress{}).String(), "")
	require.Equal(t, (types.ConsAddress{}).String(), "")

	accAddr, err := types.AccAddressFromBech32("")
	require.True(t, accAddr.Empty())
	require.Nil(t, err)

	valAddr, err := types.ValAddressFromBech32("")
	require.True(t, valAddr.Empty())
	require.Nil(t, err)

	consAddr, err := types.ConsAddressFromBech32("")
	require.True(t, consAddr.Empty())
	require.Nil(t, err)
}

func TestRandBech32PubkeyConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		rand.Read(pub[:])

		mustBech32AccPub := types.MustBech32ifyAccPub(pub)
		bech32AccPub, err := types.Bech32ifyAccPub(pub)
		require.Nil(t, err)
		require.Equal(t, bech32AccPub, mustBech32AccPub)

		mustBech32ValPub := types.MustBech32ifyValPub(pub)
		bech32ValPub, err := types.Bech32ifyValPub(pub)
		require.Nil(t, err)
		require.Equal(t, bech32ValPub, mustBech32ValPub)

		mustBech32ConsPub := types.MustBech32ifyConsPub(pub)
		bech32ConsPub, err := types.Bech32ifyConsPub(pub)
		require.Nil(t, err)
		require.Equal(t, bech32ConsPub, mustBech32ConsPub)

		mustAccPub := types.MustGetAccPubKeyBech32(bech32AccPub)
		accPub, err := types.GetAccPubKeyBech32(bech32AccPub)
		require.Nil(t, err)
		require.Equal(t, accPub, mustAccPub)

		mustValPub := types.MustGetValPubKeyBech32(bech32ValPub)
		valPub, err := types.GetValPubKeyBech32(bech32ValPub)
		require.Nil(t, err)
		require.Equal(t, valPub, mustValPub)

		mustConsPub := types.MustGetConsPubKeyBech32(bech32ConsPub)
		consPub, err := types.GetConsPubKeyBech32(bech32ConsPub)
		require.Nil(t, err)
		require.Equal(t, consPub, mustConsPub)

		require.Equal(t, valPub, accPub)
		require.Equal(t, valPub, consPub)
	}
}

func TestYAMLMarshalers(t *testing.T) {
	addr := secp256k1.GenPrivKey().PubKey().Address()

	acc := types.AccAddress(addr)
	val := types.ValAddress(addr)
	cons := types.ConsAddress(addr)

	got, _ := yaml.Marshal(&acc)
	require.Equal(t, acc.String()+"\n", string(got))

	got, _ = yaml.Marshal(&val)
	require.Equal(t, val.String()+"\n", string(got))

	got, _ = yaml.Marshal(&cons)
	require.Equal(t, cons.String()+"\n", string(got))
}

func TestRandBech32AccAddrConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		rand.Read(pub[:])

		acc := types.AccAddress(pub.Address())
		res := types.AccAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.AccAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AccAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AccAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AccAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.AccAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestValAddr(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.ValAddress(pub.Address())
		res := types.ValAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.ValAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.ValAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.ValAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.ValAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.ValAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestConsAddress(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 20; i++ {
		rand.Read(pub[:])

		acc := types.ConsAddress(pub.Address())
		res := types.ConsAddress{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err := types.ConsAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.ConsAddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.ConsAddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.ConsAddressFromBech32(str)
		require.NotNil(t, err)

		err = (*types.ConsAddress)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

var testData = []struct {
	hexAddr    string
	bech32Addr string
}{
	{"B5AD24DD9E5D60E1F0734AF2D819FF9A198A2A38", "poc1kkkjfhv7t4swrurnftedsx0lngvc523cfa00s3"},
	{"549C15831315AD56F89C0EDDF9D852B6CB7605E3", "poc12jwptqcnzkk4d7yupmwlnkzjkm9hvp0rr0chxr"},
	{"7FDE0B354783BC8607222CD994C805E303A80BCD", "poc10l0qkd28sw7gvpez9nvefjq9uvp6sz7d75mycf"},
	{"C48C652CD2A293D7F57BB36FD3B30AD87C857923", "poc1cjxx2txj52fa0atmkdha8vc2mp7g27frfcxh7h"},
	{"FEBF0CA4CB4897C9A27A54275E612FEF275752AE", "poc1l6lsefxtfztungn62sn4ucf0aun4w54wx098yt"},
	{"3BDA7843C6CE02FB1B274DF18F58E04354750EB8", "poc180d8ss7xecp0kxe8fhcc7k8qgd282r4ckmwfem"},
	{"CCE353A7008DD9E838691E5921D935848A0410F8", "poc1en348fcq3hv7swrfrevjrkf4sj9qgy8clylht8"},
	{"4DBC8579C8A7453E7547A496AB07FB48F435B1F0", "poc1fk7g27wg5aznua285jt2kplmfr6rtv0sxn42gh"},
	{"98488C3C1BDF59D448A52EC97A9410F164834CF3", "poc1npygc0qmmavagj999myh49qs79jgxn8n26fxeg"},
	{"2D44FDEFC054FD718550B60210C07982012C8D11", "poc194z0mm7q2n7hrp2skcpppsresgqjerg3px5ynk"},
	{"935FCC04364A3C68CFD4E014CC21BBC1A29845EB", "poc1jd0ucppkfg7x3n75uq2vcgdmcx3fs30tan8jtw"},
	{"03A62A1613A098DF0C607CEBFEEC7475118E02E8", "poc1qwnz59sn5zvd7rrq0n4lamr5w5gcuqhgprvekr"},
	{"F81E3C014E639ACA87BB7A7F8A38C24F355DB3DB", "poc1lq0rcq2wvwdv4pam0flc5wxzfu64mv7m68p5wu"},
}

func TestAccAddressToBech32(t *testing.T) {
	for _, d := range testData {
		addr1, err := types.AccAddressFromHex(d.hexAddr)
		require.Nil(t, err)
		addr2, err := types.AccAddressFromBech32(d.bech32Addr)
		require.Nil(t, err)
		require.Equal(t, addr1, addr2)
	}
}

func TestAddressesEquals1(t *testing.T) {
	for i := 0; i < 100; i++ {
		privKey := ed25519.GenPrivKey()
		pubKey := privKey.PubKey()
		accAddr := types.AccAddress(pubKey.Address())
		valAddr := types.ValAddress(pubKey.Address())
		consAddr := types.ConsAddress(pubKey.Address())
		require.True(t, accAddr.Equals(valAddr))
		require.True(t, accAddr.Equals(consAddr))
		require.True(t, valAddr.Equals(consAddr))
		msg := crypto.CRandBytes(128)
		signature, _ := privKey.Sign(msg)
		require.True(t, pubKey.VerifyBytes(msg, signature))
		t.Logf("accAddr Byte:%v, accAddr:%v, valAddr:%v, consAddr:%v", accAddr, accAddr.String(), valAddr.String(), consAddr.String())
		//	conspubkeystr := types.MustBech32ifyConsPub(pubKey)
		//	t.Logf("conspubkey %v", conspubkeystr)
	}
}
func testMarshalaccAddr(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

func TestAccAddressMarshal(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	for i := 0; i < 20; i++ {
		rand.Read(pub[:])
		accAddr := types.AccAddress(pub.Address())
		res := types.AccAddress{}
		testMarshalaccAddr(t, &accAddr, &res, accAddr.MarshalJSON, (&res).UnmarshalJSON)
		testMarshalaccAddr(t, &accAddr, &res, accAddr.Marshal, (&res).Unmarshal)
		str := accAddr.String()
		res, err := types.AccAddressFromBech32(str)
		require.Nil(t, err)
		require.Equal(t, accAddr, res)
		hexStr := hex.EncodeToString(accAddr)
		res, err = types.AccAddressFromHex(hexStr)
		require.Nil(t, err)
		require.Equal(t, accAddr, res)
	}
}

func TestAddressesEquals2(t *testing.T) {
	hexAddrStr := "B5AD24DD9E5D60E1F0734AF2D819FF9A198A2A38"
	accAddr, err := types.AccAddressFromHex(hexAddrStr)
	require.Nil(t, err)
	ValAddr, err := types.ValAddressFromHex(hexAddrStr)
	require.Nil(t, err)
	require.True(t, accAddr.Equals(ValAddr))
	ConAddr, err := types.ConsAddressFromHex(hexAddrStr)
	require.Nil(t, err)
	require.True(t, accAddr.Equals(ConAddr))
}

func TestAccAddressToValiadatorAddress(t *testing.T) {
	for i := 0; i < 10; i++ {
		pubKey := secp256k1.GenPrivKey().PubKey()
		accAddr := types.AccAddressFromPubKey(pubKey)
		valOperAddrStr := types.ValAddress(pubKey.Address().Bytes()).String()
		valOperAddr, err := types.ValAddressFromBech32(valOperAddrStr)
		require.Nil(t, err)
		require.Equal(t, accAddr, types.AccAddress(valOperAddr))
		valConsAddrStr := types.ConsAddress(pubKey.Address().Bytes()).String()
		valConsAddr, err := types.ConsAddressFromBech32(valConsAddrStr)
		require.Nil(t, err)
		require.Equal(t, accAddr, types.AccAddress(valConsAddr))
	}
}

//user defined prefix
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestConfiguredPrefix(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	for length := 1; length < 10; length++ {
		for times := 1; times < 20; times++ {
			rand.Read(pub[:])
			// Test if randomly generated prefix of a given length works
			prefix := RandString(length)

			// Assuming that GetConfig is not sealed.
			config := types.GetConfig()
			config.SetBech32PrefixForAccount(
				prefix+types.PrefixAccount,
				prefix+types.PrefixPublic)

			acc := types.AccAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				acc.String(),
				prefix+types.PrefixAccount), acc.String())

			bech32Pub := types.MustBech32ifyAccPub(pub)
			require.True(t, strings.HasPrefix(
				bech32Pub,
				prefix+types.PrefixPublic))

			config.SetBech32PrefixForValidator(
				prefix+types.PrefixValidator+types.PrefixAddress,
				prefix+types.PrefixValidator+types.PrefixPublic)

			val := types.ValAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				val.String(),
				prefix+types.PrefixValidator+types.PrefixAddress))

			bech32ValPub := types.MustBech32ifyValPub(pub)
			require.True(t, strings.HasPrefix(
				bech32ValPub,
				prefix+types.PrefixValidator+types.PrefixPublic))

			config.SetBech32PrefixForConsensusNode(
				prefix+types.PrefixConsensus+types.PrefixAddress,
				prefix+types.PrefixConsensus+types.PrefixPublic)

			cons := types.ConsAddress(pub.Address())
			require.True(t, strings.HasPrefix(
				cons.String(),
				prefix+types.PrefixConsensus+types.PrefixAddress))

			bech32ConsPub := types.MustBech32ifyConsPub(pub)
			require.True(t, strings.HasPrefix(
				bech32ConsPub,
				prefix+types.PrefixConsensus+types.PrefixPublic))
		}

	}
}

func TestAddressInterface(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])

	addrs := []types.Address{
		types.ConsAddress(pub.Address()),
		types.ValAddress(pub.Address()),
		types.AccAddress(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case types.AccAddress:
			_, err := types.AccAddressFromBech32(addr.String())
			require.Nil(t, err)
		case types.ValAddress:
			_, err := types.ValAddressFromBech32(addr.String())
			require.Nil(t, err)
		case types.ConsAddress:
			_, err := types.ConsAddressFromBech32(addr.String())
			require.Nil(t, err)
		default:
			t.Fail()
		}
	}

}

func TestCustomAddressVerifier(t *testing.T) {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	accBech := types.AccAddress(addr).String()
	valBech := types.ValAddress(addr).String()
	consBech := types.ConsAddress(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := types.VerifyAddressFormat(addr)
	require.NotNil(t, err)
	_, err = types.AccAddressFromBech32(accBech)
	require.NotNil(t, err)
	_, err = types.ValAddressFromBech32(valBech)
	require.NotNil(t, err)
	_, err = types.ConsAddressFromBech32(consBech)
	require.NotNil(t, err)

	// Set a custom address verifier that accepts 10 or 20 byte addresses
	types.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == 10 || n == types.AddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic accepts this 10 byte address
	err = types.VerifyAddressFormat(addr)
	require.Nil(t, err)
	_, err = types.AccAddressFromBech32(accBech)
	require.Nil(t, err)
	_, err = types.ValAddressFromBech32(valBech)
	require.Nil(t, err)
	_, err = types.ConsAddressFromBech32(consBech)
	require.Nil(t, err)
}
