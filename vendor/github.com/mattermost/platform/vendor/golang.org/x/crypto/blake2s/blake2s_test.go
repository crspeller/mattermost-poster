// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blake2s

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestHashes(t *testing.T) {
	defer func(sse2, ssse3, sse4 bool) {
		useSSE2, useSSSE3, useSSE4 = sse2, ssse3, sse4
	}(useSSE2, useSSSE3, useSSE4)

	if useSSE4 {
		t.Log("SSE4 version")
		testHashes(t)
		testHashes128(t)
		useSSE4 = false
	}
	if useSSSE3 {
		t.Log("SSSE3 version")
		testHashes(t)
		testHashes128(t)
		useSSSE3 = false
	}
	if useSSE2 {
		t.Log("SSE2 version")
		testHashes(t)
		testHashes128(t)
		useSSE2 = false
	}
	if useGeneric {
		t.Log("generic version")
		testHashes(t)
		testHashes128(t)
	}
}

func testHashes(t *testing.T) {
	key, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")

	input := make([]byte, 255)
	for i := range input {
		input[i] = byte(i)
	}

	for i, expectedHex := range hashes {
		h, err := New256(key)
		if err != nil {
			t.Fatalf("#%d: error from New256: %v", i, err)
		}

		h.Write(input[:i])
		sum := h.Sum(nil)

		if gotHex := fmt.Sprintf("%x", sum); gotHex != expectedHex {
			t.Fatalf("#%d (single write): got %s, wanted %s", i, gotHex, expectedHex)
		}

		h.Reset()
		for j := 0; j < i; j++ {
			h.Write(input[j : j+1])
		}

		sum = h.Sum(sum[:0])
		if gotHex := fmt.Sprintf("%x", sum); gotHex != expectedHex {
			t.Fatalf("#%d (byte-by-byte): got %s, wanted %s", i, gotHex, expectedHex)
		}
	}
}

func testHashes128(t *testing.T) {
	key, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")

	input := make([]byte, 255)
	for i := range input {
		input[i] = byte(i)
	}

	for i, expectedHex := range hashes128 {
		h, err := New128(key)
		if err != nil {
			t.Fatalf("#%d: error from New128: %v", i, err)
		}

		h.Write(input[:i])
		sum := h.Sum(nil)

		if gotHex := fmt.Sprintf("%x", sum); gotHex != expectedHex {
			t.Fatalf("#%d (single write): got %s, wanted %s", i, gotHex, expectedHex)
		}

		h.Reset()
		for j := 0; j < i; j++ {
			h.Write(input[j : j+1])
		}

		sum = h.Sum(sum[:0])
		if gotHex := fmt.Sprintf("%x", sum); gotHex != expectedHex {
			t.Fatalf("#%d (byte-by-byte): got %s, wanted %s", i, gotHex, expectedHex)
		}
	}
}

// Benchmarks

func benchmarkSum(b *testing.B, size int) {
	data := make([]byte, size)
	b.SetBytes(int64(size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum256(data)
	}
}

func benchmarkWrite(b *testing.B, size int) {
	data := make([]byte, size)
	h, _ := New256(nil)
	b.SetBytes(int64(size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(data)
	}
}

func BenchmarkWrite64(b *testing.B) { benchmarkWrite(b, 64) }
func BenchmarkWrite1K(b *testing.B) { benchmarkWrite(b, 1024) }

func BenchmarkSum64(b *testing.B) { benchmarkSum(b, 64) }
func BenchmarkSum1K(b *testing.B) { benchmarkSum(b, 1024) }

// hashes is taken from https://blake2.net/blake2s-test.txt
var hashes = []string{
	"48a8997da407876b3d79c0d92325ad3b89cbb754d86ab71aee047ad345fd2c49",
	"40d15fee7c328830166ac3f918650f807e7e01e177258cdc0a39b11f598066f1",
	"6bb71300644cd3991b26ccd4d274acd1adeab8b1d7914546c1198bbe9fc9d803",
	"1d220dbe2ee134661fdf6d9e74b41704710556f2f6e5a091b227697445dbea6b",
	"f6c3fbadb4cc687a0064a5be6e791bec63b868ad62fba61b3757ef9ca52e05b2",
	"49c1f21188dfd769aea0e911dd6b41f14dab109d2b85977aa3088b5c707e8598",
	"fdd8993dcd43f696d44f3cea0ff35345234ec8ee083eb3cada017c7f78c17143",
	"e6c8125637438d0905b749f46560ac89fd471cf8692e28fab982f73f019b83a9",
	"19fc8ca6979d60e6edd3b4541e2f967ced740df6ec1eaebbfe813832e96b2974",
	"a6ad777ce881b52bb5a4421ab6cdd2dfba13e963652d4d6d122aee46548c14a7",
	"f5c4b2ba1a00781b13aba0425242c69cb1552f3f71a9a3bb22b4a6b4277b46dd",
	"e33c4c9bd0cc7e45c80e65c77fa5997fec7002738541509e68a9423891e822a3",
	"fba16169b2c3ee105be6e1e650e5cbf40746b6753d036ab55179014ad7ef6651",
	"f5c4bec6d62fc608bf41cc115f16d61c7efd3ff6c65692bbe0afffb1fede7475",
	"a4862e76db847f05ba17ede5da4e7f91b5925cf1ad4ba12732c3995742a5cd6e",
	"65f4b860cd15b38ef814a1a804314a55be953caa65fd758ad989ff34a41c1eea",
	"19ba234f0a4f38637d1839f9d9f76ad91c8522307143c97d5f93f69274cec9a7",
	"1a67186ca4a5cb8e65fca0e2ecbc5ddc14ae381bb8bffeb9e0a103449e3ef03c",
	"afbea317b5a2e89c0bd90ccf5d7fd0ed57fe585e4be3271b0a6bf0f5786b0f26",
	"f1b01558ce541262f5ec34299d6fb4090009e3434be2f49105cf46af4d2d4124",
	"13a0a0c86335635eaa74ca2d5d488c797bbb4f47dc07105015ed6a1f3309efce",
	"1580afeebebb346f94d59fe62da0b79237ead7b1491f5667a90e45edf6ca8b03",
	"20be1a875b38c573dd7faaa0de489d655c11efb6a552698e07a2d331b5f655c3",
	"be1fe3c4c04018c54c4a0f6b9a2ed3c53abe3a9f76b4d26de56fc9ae95059a99",
	"e3e3ace537eb3edd8463d9ad3582e13cf86533ffde43d668dd2e93bbdbd7195a",
	"110c50c0bf2c6e7aeb7e435d92d132ab6655168e78a2decdec3330777684d9c1",
	"e9ba8f505c9c80c08666a701f3367e6cc665f34b22e73c3c0417eb1c2206082f",
	"26cd66fca02379c76df12317052bcafd6cd8c3a7b890d805f36c49989782433a",
	"213f3596d6e3a5d0e9932cd2159146015e2abc949f4729ee2632fe1edb78d337",
	"1015d70108e03be1c702fe97253607d14aee591f2413ea6787427b6459ff219a",
	"3ca989de10cfe609909472c8d35610805b2f977734cf652cc64b3bfc882d5d89",
	"b6156f72d380ee9ea6acd190464f2307a5c179ef01fd71f99f2d0f7a57360aea",
	"c03bc642b20959cbe133a0303e0c1abff3e31ec8e1a328ec8565c36decff5265",
	"2c3e08176f760c6264c3a2cd66fec6c3d78de43fc192457b2a4a660a1e0eb22b",
	"f738c02f3c1b190c512b1a32deabf353728e0e9ab034490e3c3409946a97aeec",
	"8b1880df301cc963418811088964839287ff7fe31c49ea6ebd9e48bdeee497c5",
	"1e75cb21c60989020375f1a7a242839f0b0b68973a4c2a05cf7555ed5aaec4c1",
	"62bf8a9c32a5bccf290b6c474d75b2a2a4093f1a9e27139433a8f2b3bce7b8d7",
	"166c8350d3173b5e702b783dfd33c66ee0432742e9b92b997fd23c60dc6756ca",
	"044a14d822a90cacf2f5a101428adc8f4109386ccb158bf905c8618b8ee24ec3",
	"387d397ea43a994be84d2d544afbe481a2000f55252696bba2c50c8ebd101347",
	"56f8ccf1f86409b46ce36166ae9165138441577589db08cbc5f66ca29743b9fd",
	"9706c092b04d91f53dff91fa37b7493d28b576b5d710469df79401662236fc03",
	"877968686c068ce2f7e2adcff68bf8748edf3cf862cfb4d3947a3106958054e3",
	"8817e5719879acf7024787eccdb271035566cfa333e049407c0178ccc57a5b9f",
	"8938249e4b50cadaccdf5b18621326cbb15253e33a20f5636e995d72478de472",
	"f164abba4963a44d107257e3232d90aca5e66a1408248c51741e991db5227756",
	"d05563e2b1cba0c4a2a1e8bde3a1a0d9f5b40c85a070d6f5fb21066ead5d0601",
	"03fbb16384f0a3866f4c3117877666efbf124597564b293d4aab0d269fabddfa",
	"5fa8486ac0e52964d1881bbe338eb54be2f719549224892057b4da04ba8b3475",
	"cdfabcee46911111236a31708b2539d71fc211d9b09c0d8530a11e1dbf6eed01",
	"4f82de03b9504793b82a07a0bdcdff314d759e7b62d26b784946b0d36f916f52",
	"259ec7f173bcc76a0994c967b4f5f024c56057fb79c965c4fae41875f06a0e4c",
	"193cc8e7c3e08bb30f5437aa27ade1f142369b246a675b2383e6da9b49a9809e",
	"5c10896f0e2856b2a2eee0fe4a2c1633565d18f0e93e1fab26c373e8f829654d",
	"f16012d93f28851a1eb989f5d0b43f3f39ca73c9a62d5181bff237536bd348c3",
	"2966b3cfae1e44ea996dc5d686cf25fa053fb6f67201b9e46eade85d0ad6b806",
	"ddb8782485e900bc60bcf4c33a6fd585680cc683d516efa03eb9985fad8715fb",
	"4c4d6e71aea05786413148fc7a786b0ecaf582cff1209f5a809fba8504ce662c",
	"fb4c5e86d7b2229b99b8ba6d94c247ef964aa3a2bae8edc77569f28dbbff2d4e",
	"e94f526de9019633ecd54ac6120f23958d7718f1e7717bf329211a4faeed4e6d",
	"cbd6660a10db3f23f7a03d4b9d4044c7932b2801ac89d60bc9eb92d65a46c2a0",
	"8818bbd3db4dc123b25cbba5f54c2bc4b3fcf9bf7d7a7709f4ae588b267c4ece",
	"c65382513f07460da39833cb666c5ed82e61b9e998f4b0c4287cee56c3cc9bcd",
	"8975b0577fd35566d750b362b0897a26c399136df07bababbde6203ff2954ed4",
	"21fe0ceb0052be7fb0f004187cacd7de67fa6eb0938d927677f2398c132317a8",
	"2ef73f3c26f12d93889f3c78b6a66c1d52b649dc9e856e2c172ea7c58ac2b5e3",
	"388a3cd56d73867abb5f8401492b6e2681eb69851e767fd84210a56076fb3dd3",
	"af533e022fc9439e4e3cb838ecd18692232adf6fe9839526d3c3dd1b71910b1a",
	"751c09d41a9343882a81cd13ee40818d12eb44c6c7f40df16e4aea8fab91972a",
	"5b73ddb68d9d2b0aa265a07988d6b88ae9aac582af83032f8a9b21a2e1b7bf18",
	"3da29126c7c5d7f43e64242a79feaa4ef3459cdeccc898ed59a97f6ec93b9dab",
	"566dc920293da5cb4fe0aa8abda8bbf56f552313bff19046641e3615c1e3ed3f",
	"4115bea02f73f97f629e5c5590720c01e7e449ae2a6697d4d2783321303692f9",
	"4ce08f4762468a7670012164878d68340c52a35e66c1884d5c864889abc96677",
	"81ea0b7804124e0c22ea5fc71104a2afcb52a1fa816f3ecb7dcb5d9dea1786d0",
	"fe362733b05f6bedaf9379d7f7936ede209b1f8323c3922549d9e73681b5db7b",
	"eff37d30dfd20359be4e73fdf40d27734b3df90a97a55ed745297294ca85d09f",
	"172ffc67153d12e0ca76a8b6cd5d4731885b39ce0cac93a8972a18006c8b8baf",
	"c47957f1cc88e83ef9445839709a480a036bed5f88ac0fcc8e1e703ffaac132c",
	"30f3548370cfdceda5c37b569b6175e799eef1a62aaa943245ae7669c227a7b5",
	"c95dcb3cf1f27d0eef2f25d2413870904a877c4a56c2de1e83e2bc2ae2e46821",
	"d5d0b5d705434cd46b185749f66bfb5836dcdf6ee549a2b7a4aee7f58007caaf",
	"bbc124a712f15d07c300e05b668389a439c91777f721f8320c1c9078066d2c7e",
	"a451b48c35a6c7854cfaae60262e76990816382ac0667e5a5c9e1b46c4342ddf",
	"b0d150fb55e778d01147f0b5d89d99ecb20ff07e5e6760d6b645eb5b654c622b",
	"34f737c0ab219951eee89a9f8dac299c9d4c38f33fa494c5c6eefc92b6db08bc",
	"1a62cc3a00800dcbd99891080c1e098458193a8cc9f970ea99fbeff00318c289",
	"cfce55ebafc840d7ae48281c7fd57ec8b482d4b704437495495ac414cf4a374b",
	"6746facf71146d999dabd05d093ae586648d1ee28e72617b99d0f0086e1e45bf",
	"571ced283b3f23b4e750bf12a2caf1781847bd890e43603cdc5976102b7bb11b",
	"cfcb765b048e35022c5d089d26e85a36b005a2b80493d03a144e09f409b6afd1",
	"4050c7a27705bb27f42089b299f3cbe5054ead68727e8ef9318ce6f25cd6f31d",
	"184070bd5d265fbdc142cd1c5cd0d7e414e70369a266d627c8fba84fa5e84c34",
	"9edda9a4443902a9588c0d0ccc62b930218479a6841e6fe7d43003f04b1fd643",
	"e412feef7908324a6da1841629f35d3d358642019310ec57c614836b63d30763",
	"1a2b8edff3f9acc1554fcbae3cf1d6298c6462e22e5eb0259684f835012bd13f",
	"288c4ad9b9409762ea07c24a41f04f69a7d74bee2d95435374bde946d7241c7b",
	"805691bb286748cfb591d3aebe7e6f4e4dc6e2808c65143cc004e4eb6fd09d43",
	"d4ac8d3a0afc6cfa7b460ae3001baeb36dadb37da07d2e8ac91822df348aed3d",
	"c376617014d20158bced3d3ba552b6eccf84e62aa3eb650e90029c84d13eea69",
	"c41f09f43cecae7293d6007ca0a357087d5ae59be500c1cd5b289ee810c7b082",
	"03d1ced1fba5c39155c44b7765cb760c78708dcfc80b0bd8ade3a56da8830b29",
	"09bde6f152218dc92c41d7f45387e63e5869d807ec70b821405dbd884b7fcf4b",
	"71c9036e18179b90b37d39e9f05eb89cc5fc341fd7c477d0d7493285faca08a4",
	"5916833ebb05cd919ca7fe83b692d3205bef72392b2cf6bb0a6d43f994f95f11",
	"f63aab3ec641b3b024964c2b437c04f6043c4c7e0279239995401958f86bbe54",
	"f172b180bfb09740493120b6326cbdc561e477def9bbcfd28cc8c1c5e3379a31",
	"cb9b89cc18381dd9141ade588654d4e6a231d5bf49d4d59ac27d869cbe100cf3",
	"7bd8815046fdd810a923e1984aaebdcdf84d87c8992d68b5eeb460f93eb3c8d7",
	"607be66862fd08ee5b19facac09dfdbcd40c312101d66e6ebd2b841f1b9a9325",
	"9fe03bbe69ab1834f5219b0da88a08b30a66c5913f0151963c360560db0387b3",
	"90a83585717b75f0e9b725e055eeeeb9e7a028ea7e6cbc07b20917ec0363e38c",
	"336ea0530f4a7469126e0218587ebbde3358a0b31c29d200f7dc7eb15c6aadd8",
	"a79e76dc0abca4396f0747cd7b748df913007626b1d659da0c1f78b9303d01a3",
	"44e78a773756e0951519504d7038d28d0213a37e0ce375371757bc996311e3b8",
	"77ac012a3f754dcfeab5eb996be9cd2d1f96111b6e49f3994df181f28569d825",
	"ce5a10db6fccdaf140aaa4ded6250a9c06e9222bc9f9f3658a4aff935f2b9f3a",
	"ecc203a7fe2be4abd55bb53e6e673572e0078da8cd375ef430cc97f9f80083af",
	"14a5186de9d7a18b0412b8563e51cc5433840b4a129a8ff963b33a3c4afe8ebb",
	"13f8ef95cb86e6a638931c8e107673eb76ba10d7c2cd70b9d9920bbeed929409",
	"0b338f4ee12f2dfcb78713377941e0b0632152581d1332516e4a2cab1942cca4",
	"eaab0ec37b3b8ab796e9f57238de14a264a076f3887d86e29bb5906db5a00e02",
	"23cb68b8c0e6dc26dc27766ddc0a13a99438fd55617aa4095d8f969720c872df",
	"091d8ee30d6f2968d46b687dd65292665742de0bb83dcc0004c72ce10007a549",
	"7f507abc6d19ba00c065a876ec5657868882d18a221bc46c7a6912541f5bc7ba",
	"a0607c24e14e8c223db0d70b4d30ee88014d603f437e9e02aa7dafa3cdfbad94",
	"ddbfea75cc467882eb3483ce5e2e756a4f4701b76b445519e89f22d60fa86e06",
	"0c311f38c35a4fb90d651c289d486856cd1413df9b0677f53ece2cd9e477c60a",
	"46a73a8dd3e70f59d3942c01df599def783c9da82fd83222cd662b53dce7dbdf",
	"ad038ff9b14de84a801e4e621ce5df029dd93520d0c2fa38bff176a8b1d1698c",
	"ab70c5dfbd1ea817fed0cd067293abf319e5d7901c2141d5d99b23f03a38e748",
	"1fffda67932b73c8ecaf009a3491a026953babfe1f663b0697c3c4ae8b2e7dcb",
	"b0d2cc19472dd57f2b17efc03c8d58c2283dbb19da572f7755855aa9794317a0",
	"a0d19a6ee33979c325510e276622df41f71583d07501b87071129a0ad94732a5",
	"724642a7032d1062b89e52bea34b75df7d8fe772d9fe3c93ddf3c4545ab5a99b",
	"ade5eaa7e61f672d587ea03dae7d7b55229c01d06bc0a5701436cbd18366a626",
	"013b31ebd228fcdda51fabb03bb02d60ac20ca215aafa83bdd855e3755a35f0b",
	"332ed40bb10dde3c954a75d7b8999d4b26a1c063c1dc6e32c1d91bab7bbb7d16",
	"c7a197b3a05b566bcc9facd20e441d6f6c2860ac9651cd51d6b9d2cdeeea0390",
	"bd9cf64ea8953c037108e6f654914f3958b68e29c16700dc184d94a21708ff60",
	"8835b0ac021151df716474ce27ce4d3c15f0b2dab48003cf3f3efd0945106b9a",
	"3bfefa3301aa55c080190cffda8eae51d9af488b4c1f24c3d9a75242fd8ea01d",
	"08284d14993cd47d53ebaecf0df0478cc182c89c00e1859c84851686ddf2c1b7",
	"1ed7ef9f04c2ac8db6a864db131087f27065098e69c3fe78718d9b947f4a39d0",
	"c161f2dcd57e9c1439b31a9dd43d8f3d7dd8f0eb7cfac6fb25a0f28e306f0661",
	"c01969ad34c52caf3dc4d80d19735c29731ac6e7a92085ab9250c48dea48a3fc",
	"1720b3655619d2a52b3521ae0e49e345cb3389ebd6208acaf9f13fdacca8be49",
	"756288361c83e24c617cf95c905b22d017cdc86f0bf1d658f4756c7379873b7f",
	"e7d0eda3452693b752abcda1b55e276f82698f5f1605403eff830bea0071a394",
	"2c82ecaa6b84803e044af63118afe544687cb6e6c7df49ed762dfd7c8693a1bc",
	"6136cbf4b441056fa1e2722498125d6ded45e17b52143959c7f4d4e395218ac2",
	"721d3245aafef27f6a624f47954b6c255079526ffa25e9ff77e5dcff473b1597",
	"9dd2fbd8cef16c353c0ac21191d509eb28dd9e3e0d8cea5d26ca839393851c3a",
	"b2394ceacdebf21bf9df2ced98e58f1c3a4bbbff660dd900f62202d6785cc46e",
	"57089f222749ad7871765f062b114f43ba20ec56422a8b1e3f87192c0ea718c6",
	"e49a9459961cd33cdf4aae1b1078a5dea7c040e0fea340c93a724872fc4af806",
	"ede67f720effd2ca9c88994152d0201dee6b0a2d2c077aca6dae29f73f8b6309",
	"e0f434bf22e3088039c21f719ffc67f0f2cb5e98a7a0194c76e96bf4e8e17e61",
	"277c04e2853484a4eba910ad336d01b477b67cc200c59f3c8d77eef8494f29cd",
	"156d5747d0c99c7f27097d7b7e002b2e185cb72d8dd7eb424a0321528161219f",
	"20ddd1ed9b1ca803946d64a83ae4659da67fba7a1a3eddb1e103c0f5e03e3a2c",
	"f0af604d3dabbf9a0f2a7d3dda6bd38bba72c6d09be494fcef713ff10189b6e6",
	"9802bb87def4cc10c4a5fd49aa58dfe2f3fddb46b4708814ead81d23ba95139b",
	"4f8ce1e51d2fe7f24043a904d898ebfc91975418753413aa099b795ecb35cedb",
	"bddc6514d7ee6ace0a4ac1d0e068112288cbcf560454642705630177cba608bd",
	"d635994f6291517b0281ffdd496afa862712e5b3c4e52e4cd5fdae8c0e72fb08",
	"878d9ca600cf87e769cc305c1b35255186615a73a0da613b5f1c98dbf81283ea",
	"a64ebe5dc185de9fdde7607b6998702eb23456184957307d2fa72e87a47702d6",
	"ce50eab7b5eb52bdc9ad8e5a480ab780ca9320e44360b1fe37e03f2f7ad7de01",
	"eeddb7c0db6e30abe66d79e327511e61fcebbc29f159b40a86b046ecf0513823",
	"787fc93440c1ec96b5ad01c16cf77916a1405f9426356ec921d8dff3ea63b7e0",
	"7f0d5eab47eefda696c0bf0fbf86ab216fce461e9303aba6ac374120e890e8df",
	"b68004b42f14ad029f4c2e03b1d5eb76d57160e26476d21131bef20ada7d27f4",
	"b0c4eb18ae250b51a41382ead92d0dc7455f9379fc9884428e4770608db0faec",
	"f92b7a870c059f4d46464c824ec96355140bdce681322cc3a992ff103e3fea52",
	"5364312614813398cc525d4c4e146edeb371265fba19133a2c3d2159298a1742",
	"f6620e68d37fb2af5000fc28e23b832297ecd8bce99e8be4d04e85309e3d3374",
	"5316a27969d7fe04ff27b283961bffc3bf5dfb32fb6a89d101c6c3b1937c2871",
	"81d1664fdf3cb33c24eebac0bd64244b77c4abea90bbe8b5ee0b2aafcf2d6a53",
	"345782f295b0880352e924a0467b5fbc3e8f3bfbc3c7e48b67091fb5e80a9442",
	"794111ea6cd65e311f74ee41d476cb632ce1e4b051dc1d9e9d061a19e1d0bb49",
	"2a85daf6138816b99bf8d08ba2114b7ab07975a78420c1a3b06a777c22dd8bcb",
	"89b0d5f289ec16401a069a960d0b093e625da3cf41ee29b59b930c5820145455",
	"d0fdcb543943fc27d20864f52181471b942cc77ca675bcb30df31d358ef7b1eb",
	"b17ea8d77063c709d4dc6b879413c343e3790e9e62ca85b7900b086f6b75c672",
	"e71a3e2c274db842d92114f217e2c0eac8b45093fdfd9df4ca7162394862d501",
	"c0476759ab7aa333234f6b44f5fd858390ec23694c622cb986e769c78edd733e",
	"9ab8eabb1416434d85391341d56993c55458167d4418b19a0f2ad8b79a83a75b",
	"7992d0bbb15e23826f443e00505d68d3ed7372995a5c3e498654102fbcd0964e",
	"c021b30085151435df33b007ccecc69df1269f39ba25092bed59d932ac0fdc28",
	"91a25ec0ec0d9a567f89c4bfe1a65a0e432d07064b4190e27dfb81901fd3139b",
	"5950d39a23e1545f301270aa1a12f2e6c453776e4d6355de425cc153f9818867",
	"d79f14720c610af179a3765d4b7c0968f977962dbf655b521272b6f1e194488e",
	"e9531bfc8b02995aeaa75ba27031fadbcbf4a0dab8961d9296cd7e84d25d6006",
	"34e9c26a01d7f16181b454a9d1623c233cb99d31c694656e9413aca3e918692f",
	"d9d7422f437bd439ddd4d883dae2a08350173414be78155133fff1964c3d7972",
	"4aee0c7aaf075414ff1793ead7eaca601775c615dbd60b640b0a9f0ce505d435",
	"6bfdd15459c83b99f096bfb49ee87b063d69c1974c6928acfcfb4099f8c4ef67",
	"9fd1c408fd75c336193a2a14d94f6af5adf050b80387b4b010fb29f4cc72707c",
	"13c88480a5d00d6c8c7ad2110d76a82d9b70f4fa6696d4e5dd42a066dcaf9920",
	"820e725ee25fe8fd3a8d5abe4c46c3ba889de6fa9191aa22ba67d5705421542b",
	"32d93a0eb02f42fbbcaf2bad0085b282e46046a4df7ad10657c9d6476375b93e",
	"adc5187905b1669cd8ec9c721e1953786b9d89a9bae30780f1e1eab24a00523c",
	"e90756ff7f9ad810b239a10ced2cf9b2284354c1f8c7e0accc2461dc796d6e89",
	"1251f76e56978481875359801db589a0b22f86d8d634dc04506f322ed78f17e8",
	"3afa899fd980e73ecb7f4d8b8f291dc9af796bc65d27f974c6f193c9191a09fd",
	"aa305be26e5deddc3c1010cbc213f95f051c785c5b431e6a7cd048f161787528",
	"8ea1884ff32e9d10f039b407d0d44e7e670abd884aeee0fb757ae94eaa97373d",
	"d482b2155d4dec6b4736a1f1617b53aaa37310277d3fef0c37ad41768fc235b4",
	"4d413971387e7a8898a8dc2a27500778539ea214a2dfe9b3d7e8ebdce5cf3db3",
	"696e5d46e6c57e8796e4735d08916e0b7929b3cf298c296d22e9d3019653371c",
	"1f5647c1d3b088228885865c8940908bf40d1a8272821973b160008e7a3ce2eb",
	"b6e76c330f021a5bda65875010b0edf09126c0f510ea849048192003aef4c61c",
	"3cd952a0beada41abb424ce47f94b42be64e1ffb0fd0782276807946d0d0bc55",
	"98d92677439b41b7bb513312afb92bcc8ee968b2e3b238cecb9b0f34c9bb63d0",
	"ecbca2cf08ae57d517ad16158a32bfa7dc0382eaeda128e91886734c24a0b29d",
	"942cc7c0b52e2b16a4b89fa4fc7e0bf609e29a08c1a8543452b77c7bfd11bb28",
	"8a065d8b61a0dffb170d5627735a76b0e9506037808cba16c345007c9f79cf8f",
	"1b9fa19714659c78ff413871849215361029ac802b1cbcd54e408bd87287f81f",
	"8dab071bcd6c7292a9ef727b4ae0d86713301da8618d9a48adce55f303a869a1",
	"8253e3e7c7b684b9cb2beb014ce330ff3d99d17abbdbabe4f4d674ded53ffc6b",
	"f195f321e9e3d6bd7d074504dd2ab0e6241f92e784b1aa271ff648b1cab6d7f6",
	"27e4cc72090f241266476a7c09495f2db153d5bcbd761903ef79275ec56b2ed8",
	"899c2405788e25b99a1846355e646d77cf400083415f7dc5afe69d6e17c00023",
	"a59b78c4905744076bfee894de707d4f120b5c6893ea0400297d0bb834727632",
	"59dc78b105649707a2bb4419c48f005400d3973de3736610230435b10424b24f",
	"c0149d1d7e7a6353a6d906efe728f2f329fe14a4149a3ea77609bc42b975ddfa",
	"a32f241474a6c16932e9243be0cf09bcdc7e0ca0e7a6a1b9b1a0f01e41502377",
	"b239b2e4f81841361c1339f68e2c359f929af9ad9f34e01aab4631ad6d5500b0",
	"85fb419c7002a3e0b4b6ea093b4c1ac6936645b65dac5ac15a8528b7b94c1754",
	"9619720625f190b93a3fad186ab314189633c0d3a01e6f9bc8c4a8f82f383dbf",
	"7d620d90fe69fa469a6538388970a1aa09bb48a2d59b347b97e8ce71f48c7f46",
	"294383568596fb37c75bbacd979c5ff6f20a556bf8879cc72924855df9b8240e",
	"16b18ab314359c2b833c1c6986d48c55a9fc97cde9a3c1f10a3177140f73f738",
	"8cbbdd14bc33f04cf45813e4a153a273d36adad5ce71f499eeb87fb8ac63b729",
	"69c9a498db174ecaefcc5a3ac9fdedf0f813a5bec727f1e775babdec7718816e",
	"b462c3be40448f1d4f80626254e535b08bc9cdcff599a768578d4b2881a8e3f0",
	"553e9d9c5f360ac0b74a7d44e5a391dad4ced03e0c24183b7e8ecabdf1715a64",
	"7a7c55a56fa9ae51e655e01975d8a6ff4ae9e4b486fcbe4eac044588f245ebea",
	"2afdf3c82abc4867f5de111286c2b3be7d6e48657ba923cfbf101a6dfcf9db9a",
	"41037d2edcdce0c49b7fb4a6aa0999ca66976c7483afe631d4eda283144f6dfc",
	"c4466f8497ca2eeb4583a0b08e9d9ac74395709fda109d24f2e4462196779c5d",
	"75f609338aa67d969a2ae2a2362b2da9d77c695dfd1df7224a6901db932c3364",
	"68606ceb989d5488fc7cf649f3d7c272ef055da1a93faecd55fe06f6967098ca",
	"44346bdeb7e052f6255048f0d9b42c425bab9c3dd24168212c3ecf1ebf34e6ae",
	"8e9cf6e1f366471f2ac7d2ee9b5e6266fda71f8f2e4109f2237ed5f8813fc718",
	"84bbeb8406d250951f8c1b3e86a7c010082921833dfd9555a2f909b1086eb4b8",
	"ee666f3eef0f7e2a9c222958c97eaf35f51ced393d714485ab09a069340fdf88",
	"c153d34a65c47b4a62c5cacf24010975d0356b2f32c8f5da530d338816ad5de6",
	"9fc5450109e1b779f6c7ae79d56c27635c8dd426c5a9d54e2578db989b8c3b4e",
	"d12bf3732ef4af5c22fa90356af8fc50fcb40f8f2ea5c8594737a3b3d5abdbd7",
	"11030b9289bba5af65260672ab6fee88b87420acef4a1789a2073b7ec2f2a09e",
	"69cb192b8444005c8c0ceb12c846860768188cda0aec27a9c8a55cdee2123632",
	"db444c15597b5f1a03d1f9edd16e4a9f43a667cc275175dfa2b704e3bb1a9b83",
	"3fb735061abc519dfe979e54c1ee5bfad0a9d858b3315bad34bde999efd724dd",
}

var hashes128 = []string{
	"9536f9b267655743dee97b8a670f9f53",
	"13bacfb85b48a1223c595f8c1e7e82cb",
	"d47a9b1645e2feae501cd5fe44ce6333",
	"1e2a79436a7796a3e9826bfedf07659f",
	"7640360ed3c4f3054dba79a21dda66b7",
	"d1207ac2bf5ac84fc9ef016da5a46a86",
	"3123987871e59305ece3125abfc0099a",
	"cf9e072ad522f2cda2d825218086731c",
	"95d22870392efe2846b12b6e8e84efbb",
	"7d63c30e2d51333f245601b038c0b93b",
	"ed608b98e13976bdf4bedc63fa35e443",
	"ed704b5cd1abf8e0dd67a6ac667a3fa5",
	"77dc70109827dc74c70fd26cba379ae5",
	"d2bf34508b07825ee934f33958f4560e",
	"a340baa7b8a93a6e658adef42e78eeb7",
	"b85c5ceaecbe9a251eac76f6932ba395",
	"246519722001f6e8e97a2183f5985e53",
	"5bce5aa0b7c6cac2ecf6406183cd779a",
	"13408f1647c02f6efd0047ad8344f695",
	"a63970f196760aa36cb965ab62f0e0fa",
	"bc26f48421dd99fd45e15e736d3e7dac",
	"4c6f70f9e3237cde918afb52d26f1823",
	"45ed610cfbc37db80c4bf0eef14ae8d6",
	"87c4c150705ea5078209ec008200539c",
	"54de21f5e0e6f2afe04daeb822b6931e",
	"9732a04e505064e19de3d542e7e71631",
	"d2bd27e95531d6957eef511c4ba64ad4",
	"7a36c9f70dcc7c3063b547101a5f6c35",
	"322007d1a44c4257bc7903b183305529",
	"dbcc9a09f412290ca2e0d53dfd142ddb",
	"df12ed43b8e53a56db20e0f83764002c",
	"d114cc11e7d5b33a360c45f18d4c7c6e",
	"c43b5e836af88620a8a71b1652cb8640",
	"9491c653e8867ed73c1b4ac6b5a9bb4d",
	"06d0e988df94ada6c6f9f36f588ab7c5",
	"561efad2480e93262c8eeaa3677615c4",
	"ba8ffc702e5adc93503045eca8702312",
	"5782be6ccdc78c8425285e85de8ccdc6",
	"aa1c4393e4c07b53ea6e2b5b1e970771",
	"42a229dc50e52271c51e8666023ebc1e",
	"53706110e919f84de7f8d6c7f0e7b831",
	"fc5ac8ee39cc1dd1424391323e2901bd",
	"bed27b62ff66cac2fbb68193c727106a",
	"cd5e689b96d0b9ea7e08dac36f7b211e",
	"0b4c7f604eba058d18e322c6e1baf173",
	"eb838227fdfad09a27f0f8413120675d",
	"3149cf9d19a7fd529e6154a8b4c3b3ad",
	"ca1e20126df930fd5fb7afe4422191e5",
	"b23398f910599f3c09b6549fa81bcb46",
	"27fb17c11b34fa5d8b5afe5ee3321ead",
	"0f665f5f04cf2d46b7fead1a1f328158",
	"8f068be73b3681f99f3b282e3c02bba5",
	"ba189bbd13808dcf4e002a4dd21660d5",
	"2732dcd1b16668ae6ab6a61595d0d62a",
	"d410ccdd059f0e02b472ec9ec54bdd3c",
	"b2eaa07b055b3a03a399971327f7e8c2",
	"2e8a225655e9f99b69c60dc8b4d8e566",
	"4eb55416c853f2152e67f8a224133cec",
	"49552403790d8de0505a8e317a443687",
	"7f2747cd41f56942752e868212c7d5ac",
	"02a28f10e193b430df7112d2d98cf759",
	"d4213404a9f1cf759017747cf5958270",
	"faa34884344f9c65e944882db8476d34",
	"ece382a8bd5018f1de5da44b72cea75b",
	"f1efa90d2547036841ecd3627fafbc36",
	"811ff8686d23a435ecbd0bdafcd27b1b",
	"b21beea9c7385f657a76558530438721",
	"9cb969da4f1b4fc5b13bf78fe366f0c4",
	"8850d16d7b614d3268ccfa009d33c7fc",
	"aa98a2b6176ea86415b9aff3268c6f6d",
	"ec3e1efa5ed195eff667e16b1af1e39e",
	"e40787dca57411d2630db2de699beb08",
	"554835890735babd06318de23d31e78a",
	"493957feecddc302ee2bb2086b6ebfd3",
	"f6069709ad5b0139163717e9ce1114ab",
	"ba5ed386098da284484b211555505a01",
	"9244c8dfad8cbb68c118fa51465b3ae4",
	"51e309a5008eb1f5185e5cc007cfb36f",
	"6ce9ff712121b4f6087955f4911eafd4",
	"59b51d8dcda031218ccdd7c760828155",
	"0012878767a3d4f1c8194458cf1f8832",
	"82900708afd5b6582dc16f008c655edd",
	"21302c7e39b5a4cdf1d6f86b4f00c9b4",
	"e894c7431591eab8d1ce0fe2aa1f01df",
	"b67e1c40ee9d988226d605621854d955",
	"6237bdafa34137cbbec6be43ea9bd22c",
	"4172a8e19b0dcb09b978bb9eff7af52b",
	"5714abb55bd4448a5a6ad09fbd872fdf",
	"7ce1700bef423e1f958a94a77a94d44a",
	"3742ec50cded528527775833453e0b26",
	"5d41b135724c7c9c689495324b162f18",
	"85c523333c6442c202e9e6e0f1185f93",
	"5c71f5222d40ff5d90e7570e71ab2d30",
	"6e18912e83d012efb4c66250ced6f0d9",
	"4add4448c2e35e0b138a0bac7b4b1775",
	"c0376c6bc5e7b8b9d2108ec25d2aab53",
	"f72261d5ed156765c977751c8a13fcc1",
	"cff4156c48614b6ceed3dd6b9058f17e",
	"36bfb513f76c15f514bcb593419835aa",
	"166bf48c6bffaf8291e6fdf63854bef4",
	"0b67d33f8b859c3157fbabd9e6e47ed0",
	"e4da659ca76c88e73a9f9f10f3d51789",
	"33c1ae2a86b3f51c0642e6ed5b5aa1f1",
	"27469b56aca2334449c1cf4970dcd969",
	"b7117b2e363378aa0901b0d6a9f6ddc0",
	"a9578233b09e5cd5231943fdb12cd90d",
	"486d7d75253598b716a068243c1c3e89",
	"66f6b02d682b78ffdc85e9ec86852489",
	"38a07b9a4b228fbcc305476e4d2e05d2",
	"aedb61c7970e7d05bf9002dae3c6858c",
	"c03ef441f7dd30fdb61ad2d4d8e4c7da",
	"7f45cc1eea9a00cb6aeb2dd748361190",
	"a59538b358459132e55160899e47bd65",
	"137010fef72364411820c3fbed15c8df",
	"d8362b93fc504500dbd33ac74e1b4d70",
	"a7e49f12c8f47e3b29cf8c0889b0a9c8",
	"072e94ffbfc684bd8ab2a1b9dade2fd5",
	"5ab438584bd2229e452052e002631a5f",
	"f233d14221097baef57d3ec205c9e086",
	"3a95db000c4a8ff98dc5c89631a7f162",
	"0544f18c2994ab4ddf1728f66041ff16",
	"0bc02116c60a3cc331928d6c9d3ba37e",
	"b189dca6cb5b813c74200834fba97f29",
	"ac8aaab075b4a5bc24419da239212650",
	"1e9f19323dc71c29ae99c479dc7e8df9",
	"12d944c3fa7caa1b3d62adfc492274dd",
	"b4c68f1fffe8f0030e9b18aad8c9dc96",
	"25887fab1422700d7fa3edc0b20206e2",
	"8c09f698d03eaf88abf69f8147865ef6",
	"5c363ae42a5bec26fbc5e996428d9bd7",
	"7fdfc2e854fbb3928150d5e3abcf56d6",
	"f0c944023f714df115f9e4f25bcdb89b",
	"6d19534b4c332741c8ddd79a9644de2d",
	"32595eb23764fbfc2ee7822649f74a12",
	"5a51391aab33c8d575019b6e76ae052a",
	"98b861ce2c620f10f913af5d704a5afd",
	"b7fe2fc8b77fb1ce434f8465c7ddf793",
	"0e8406e0cf8e9cc840668ece2a0fc64e",
	"b89922db99c58f6a128ccffe19b6ce60",
	"e1be9af665f0932b77d7f5631a511db7",
	"74b96f20f58de8dc9ff5e31f91828523",
	"36a4cfef5a2a7d8548db6710e50b3009",
	"007e95e8d3b91948a1dedb91f75de76b",
	"a87a702ce08f5745edf765bfcd5fbe0d",
	"847e69a388a749a9c507354d0dddfe09",
	"07176eefbc107a78f058f3d424ca6a54",
	"ad7e80682333b68296f6cb2b4a8e446d",
	"53c4aba43896ae422e5de5b9edbd46bf",
	"33bd6c20ca2a7ab916d6e98003c6c5f8",
	"060d088ea94aa093f9981a79df1dfcc8",
	"5617b214b9df08d4f11e58f5e76d9a56",
	"ca3a60ee85bd971e1daf9f7db059d909",
	"cd2b7754505d8c884eddf736f1ec613e",
	"f496163b252f1439e7e113ba2ecabd8e",
	"5719c7dcf9d9f756d6213354acb7d5cf",
	"6f7dd40b245c54411e7a9be83ae5701c",
	"c8994dd9fdeb077a45ea04a30358b637",
	"4b1184f1e35458c1c747817d527a252f",
	"fc7df674afeac7a3fd994183f4c67a74",
	"4f68e05ce4dcc533acf9c7c01d95711e",
	"d4ebc59e918400720035dfc88e0c486a",
	"d3105dd6fa123e543b0b3a6e0eeaea9e",
	"874196128ed443f5bdb2800ca048fcad",
	"01645f134978dc8f9cf0abc93b53780e",
	"5b8b64caa257873a0ffd47c981ef6c3f",
	"4ee208fc50ba0a6e65c5b58cec44c923",
	"53f409a52427b3b7ffabb057ca088428",
	"c1d6cd616f5341a93d921e356e5887a9",
	"e85c20fea67fa7320dc23379181183c8",
	"7912b6409489df001b7372bc94aebde7",
	"e559f761ec866a87f1f331767fafc60f",
	"20a6f5a36bc37043d977ed7708465ef8",
	"6a72f526965ab120826640dd784c6cc4",
	"bf486d92ad68e87c613689dd370d001b",
	"d339fd0eb35edf3abd6419c8d857acaf",
	"9521cd7f32306d969ddabc4e6a617f52",
	"a1cd9f3e81520842f3cf6cc301cb0021",
	"18e879b6f154492d593edd3f4554e237",
	"66e2329c1f5137589e051592587e521e",
	"e899566dd6c3e82cbc83958e69feb590",
	"8a4b41d7c47e4e80659d77b4e4bfc9ae",
	"f1944f6fcfc17803405a1101998c57dd",
	"f6bcec07567b4f72851b307139656b18",
	"22e7bb256918fe9924dce9093e2d8a27",
	"dd25b925815fe7b50b7079f5f65a3970",
	"0457f10f299acf0c230dd4007612e58f",
	"ecb420c19efd93814fae2964d69b54af",
	"14eb47b06dff685d88751c6e32789db4",
	"e8f072dbb50d1ab6654aa162604a892d",
	"69cff9c62092332f03a166c7b0034469",
	"d3619f98970b798ca32c6c14cd25af91",
	"2246d423774ee9d51a551e89c0539d9e",
	"75e5d1a1e374a04a699247dad827b6cf",
	"6d087dd1d4cd15bf47db07c7a96b1db8",
	"967e4c055ac51b4b2a3e506cebd5826f",
	"7417aa79247e473401bfa92a25b62e2a",
	"24f3f4956da34b5c533d9a551ccd7b16",
	"0c40382de693a5304e2331eb951cc962",
	"9436f949d51b347db5c8e6258dafaaac",
	"d2084297fe84c4ba6e04e4fb73d734fe",
	"42a6f8ff590af21b512e9e088257aa34",
	"c484ad06b1cdb3a54f3f6464a7a2a6fd",
	"1b8ac860f5ceb4365400a201ed2917aa",
	"c43eadabbe7b7473f3f837fc52650f54",
	"0e5d3205406126b1f838875deb150d6a",
	"6bf4946f8ec8a9c417f50cd1e67565be",
	"42f09a2522314799c95b3fc121a0e3e8",
	"06b8f1487f691a3f7c3f74e133d55870",
	"1a70a65fb4f314dcf6a31451a9d2704f",
	"7d4acdd0823279fd28a1e48b49a04669",
	"09545cc8822a5dfc93bbab708fd69174",
	"efc063db625013a83c9a426d39a9bddb",
	"213bbf89b3f5be0ffdb14854bbcb2588",
	"b69624d89fe2774df9a6f43695d755d4",
	"c0f9ff9ded82bd73c512e365a894774d",
	"d1b68507ed89c17ead6f69012982db71",
	"14cf16db04648978e35c44850855d1b0",
	"9f254d4eccab74cd91d694df863650a8",
	"8f8946e2967baa4a814d36ff01d20813",
	"6b9dc4d24ecba166cb2915d7a6cba43b",
	"eb35a80418a0042b850e294db7898d4d",
	"f55f925d280c637d54055c9df088ef5f",
	"f48427a04f67e33f3ba0a17f7c9704a7",
	"4a9f5bfcc0321aea2eced896cee65894",
	"8723a67d1a1df90f1cef96e6fe81e702",
	"c166c343ee25998f80bad4067960d3fd",
	"dab67288d16702e676a040fd42344d73",
	"c8e9e0d80841eb2c116dd14c180e006c",
	"92294f546bacf0dea9042c93ecba8b34",
	"013705b1502b37369ad22fe8237d444e",
	"9b97f8837d5f2ebab0768fc9a6446b93",
	"7e7e5236b05ec35f89edf8bf655498e7",
	"7be8f2362c174c776fb9432fe93bf259",
	"2422e80420276d2df5702c6470879b01",
	"df645795db778bcce23bbe819a76ba48",
	"3f97a4ac87dfc58761cda1782d749074",
	"50e3f45df21ebfa1b706b9c0a1c245a8",
	"7879541c7ff612c7ddf17cb8f7260183",
	"67f6542b903b7ba1945eba1a85ee6b1c",
	"b34b73d36ab6234b8d3f5494d251138e",
	"0aea139641fdba59ab1103479a96e05f",
	"02776815a87b8ba878453666d42afe3c",
	"5929ab0a90459ebac5a16e2fb37c847e",
	"c244def5b20ce0468f2b5012d04ac7fd",
	"12116add6fefce36ed8a0aeccce9b6d3",
	"3cd743841e9d8b878f34d91b793b4fad",
	"45e87510cf5705262185f46905fae35f",
	"276047016b0bfb501b2d4fc748165793",
	"ddd245df5a799417d350bd7f4e0b0b7e",
	"d34d917a54a2983f3fdbc4b14caae382",
	"7730fbc09d0c1fb1939a8fc436f6b995",
	"eb4899ef257a1711cc9270a19702e5b5",
	"8a30932014bce35bba620895d374df7a",
	"1924aabf9c50aa00bee5e1f95b5d9e12",
	"1758d6f8b982aec9fbe50f20e3082b46",
	"cd075928ab7e6883e697fe7fd3ac43ee",
}
