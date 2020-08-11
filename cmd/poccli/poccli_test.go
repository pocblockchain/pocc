package main

import (
	"github.com/pocblockchain/pocc/client/lcd"
	"github.com/pocblockchain/pocc/pocapp"
	"github.com/pocblockchain/pocc/x/auth/client/rest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRestEncode(t *testing.T) {
	cdc := pocapp.MakeCodec()
	rs := lcd.NewRestServer(cdc)
	registerRoutes(rs)
	rs.RegisterSwaggerUI()

	in := []byte("{\"tx\":{\"type\":\"poc/StdTx\",\"value\":{\"msg\":[{\"type\":\"poc/MsgSend\",\"value\":{\"from_address\":\"poc1uhddf9ncrrxs6qfl64cqz8kj4ez8rdhgpn5a8j\",\"to_address\":\"poc1mq3cnrrm6cluj06u05vjxgy2snwy8jeajncx5a\",\"amount\":[{\"denom\":\"poc\",\"amount\":\"1000000000000000000\"}]}}],\"fee\":{\"amount\":[{\"denom\":\"poc\",\"amount\":\"200000000000000000000\"}],\"gas\":\"200000\"},\"signatures\":[{\"pub_key\":{\"type\":\"tendermint/PubKeySecp256k1\",\"value\":\"AusYTPffSVKqhtaPoMoa0AeGZpfcWFVP88rvY074JNQI\"},\"signature\":\"Sy0MGJZTOz4onj6ia/aEtw920dlu9NfB0JaM4f0fzHhKsufm7ir1i9cJLRES4Zum7XB8aSe8qWBiHTJ5Gb/OOw==\"}],\"memo\":\"\"}},\"mode\":\"block\"}")
	t.Logf("in:%s\n", in)
	t.Logf("bz:%v\n", in)
	req := rest.BroadcastReq{}

	err := cdc.UnmarshalJSON(in, &req)
	require.Nil(t, err)

	t.Logf("req:%v", req)
}

/*
   TestRestEncode: poccli_test.go:18: in:{"tx":                {"type":"poc/StdTx","value":{"msg":[{"type":"poc/MsgSend","value":{"from_address":"poc1uhddf9ncrrxs6qfl64cqz8kj4ez8rdhgpn5a8j","to_address":"poc1mq3cnrrm6cluj06u05vjxgy2snwy8jeajncx5a","amount":[{"denom":"poc","amount":"1000000000000000000"}]}}],"fee":{"amount":[{"denom":"poc","amount":"200000000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AusYTPffSVKqhtaPoMoa0AeGZpfcWFVP88rvY074JNQI"},"signature":"Sy0MGJZTOz4onj6ia/aEtw920dlu9NfB0JaM4f0fzHhKsufm7ir1i9cJLRES4Zum7XB8aSe8qWBiHTJ5Gb/OOw=="}],"memo":""}},"mode":"block"}
   TestRestEncode: poccli_test.go:19: bz:[123 34 116 120 34 58 123 34 116 121 112 101 34 58 34 112 111 99 47 83 116 100 84 120 34 44 34 118 97 108 117 101 34 58 123 34 109 115 103 34 58 91 123 34 116 121 112 101 34 58 34 112 111 99 47 77 115 103 83 101 110 100 34 44 34 118 97 108 117 101 34 58 123 34 102 114 111 109 95 97 100 100 114 101 115 115 34 58 34 112 111 99 49 117 104 100 100 102 57 110 99 114 114 120 115 54 113 102 108 54 52 99 113 122 56 107 106 52 101 122 56 114 100 104 103 112 110 53 97 56 106 34 44 34 116 111 95 97 100 100 114 101 115 115 34 58 34 112 111 99 49 109 113 51 99 110 114 114 109 54 99 108 117 106 48 54 117 48 53 118 106 120 103 121 50 115 110 119 121 56 106 101 97 106 110 99 120 53 97 34 44 34 97 109 111 117 110 116 34 58 91 123 34 100 101 110 111 109 34 58 34 112 111 99 34 44 34 97 109 111 117 110 116 34 58 34 49 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 34 125 93 125 125 93 44 34 102 101 101 34 58 123 34 97 109 111 117 110 116 34 58 91 123 34 100 101 110 111 109 34 58 34 112 111 99 34 44 34 97 109 111 117 110 116 34 58 34 50 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 34 125 93 44 34 103 97 115 34 58 34 50 48 48 48 48 48 34 125 44 34 115 105 103 110 97 116 117 114 101 115 34 58 91 123 34 112 117 98 95 107 101 121 34 58 123 34 116 121 112 101 34 58 34 116 101 110 100 101 114 109 105 110 116 47 80 117 98 75 101 121 83 101 99 112 50 53 54 107 49 34 44 34 118 97 108 117 101 34 58 34 65 117 115 89 84 80 102 102 83 86 75 113 104 116 97 80 111 77 111 97 48 65 101 71 90 112 102 99 87 70 86 80 56 56 114 118 89 48 55 52 74 78 81 73 34 125 44 34 115 105 103 110 97 116 117 114 101 34 58 34 83 121 48 77 71 74 90 84 79 122 52 111 110 106 54 105 97 47 97 69 116 119 57 50 48 100 108 117 57 78 102 66 48 74 97 77 52 102 48 102 122 72 104 75 115 117 102 109 55 105 114 49 105 57 99 74 76 82 69 83 52 90 117 109 55 88 66 56 97 83 101 56 113 87 66 105 72 84 74 53 71 98 47 79 79 119 61 61 34 125 93 44 34 109 101 109 111 34 58 34 34 125 125 44 34 109 111 100 101 34 58 34 98 108 111 99 107 34 125]
x                                                                                  e   "  :  "  p   o   c /   S  t   d  T   x   "  ,  "  v

*/
