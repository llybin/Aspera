package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/template"

	"github.com/imroc/req"
)

const (
	getBytesURI = "https://wallet.burst.cryptoguru.org:8125/burst?requestType=getTransactionBytes&transaction="
	getTxURI    = "https://wallet.burst.cryptoguru.org:8125/burst?requestType=getTransaction&transaction="
)

type TransactionBytes struct {
	TransactionBytes string `json:"transactionBytes"`
}

var fileHeader = `package transaction

type ParseTest struct {
	header    Header
	txByteStr string
}

var ParseTransactionTests = []ParseTest{
`

var fileFooter = "}"

var parseTestTmpl = `ParseTest{
        header: Header{
                Type:                          {{.Type}},
                SubtypeAndVersion:             {{.SubtypeAndVersion}},
                Timestamp:                     {{.Timestamp}},
                Deadline:                      {{.Deadline}},
                SenderPublicKey:               []byte{ {{.SenderPublicKeyByteStr}} },
                RecipientID:                   {{.Recipient}},
                AmountNQT:                     {{.AmountNQT}},
                FeeNQT:                        {{.FeeNQT}},
                ReferencedTransactionFullHash: []byte{ {{.ReferencedTransactionFullHashByteStr}} },
                Signature:                     []byte{ {{.SignatureByteStr}} },
        },
        txByteStr: "{{.TxByteStr}}",
},
`

type Transaction struct {
	SenderPublicKey                      string `json:"senderPublicKey"`
	SenderPublicKeyByteStr               string `json:"-"`
	Signature                            string `json:"signature"`
	SignatureByteStr                     string `json:"-"`
	FeeNQT                               string `json:"feeNQT"`
	RequestProcessingTime                int    `json:"requestProcessingTime"`
	Type                                 int    `json:"type"`
	Confirmations                        int    `json:"confirmations"`
	FullHash                             string `json:"fullHash"`
	Version                              uint8  `json:"version"`
	SignatureHash                        string `json:"signatureHash"`
	SenderRS                             string `json:"senderRS"`
	Subtype                              uint8  `json:"subtype"`
	AmountNQT                            string `json:"amountNQT"`
	Sender                               string `json:"sender"`
	Block                                string `json:"block"`
	BlockTimestamp                       int    `json:"blockTimestamp"`
	Deadline                             int    `json:"deadline"`
	Transaction                          string `json:"transaction"`
	Timestamp                            int    `json:"timestamp"`
	Height                               int    `json:"height"`
	ReferencedTransactionFullHash        string `json:"referencedTransactionFullHash"`
	ReferencedTransactionFullHashByteStr string `json:"-"`
	Recipient                            uint64 `json:"recipient,string"`
	TxByteStr                            string `json:"-"`

	SubtypeAndVersion uint8

	Attachment interface{} `json:"attachment"`
}

func main() {
	f, err := os.Create("pkg/transaction/transaction_test_structs.go")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte(fileHeader))

	t, err := template.New("parse test template").Parse(parseTestTmpl)
	if err != nil {
		log.Fatal(err)
	}

	txIDs := []uint64{
		9529303547091196674,
		2567179027297123572,
		12563021031312925592,
		557538880898556285,
		1564237587434130001,
		15387360786593774841,
		11919316017481681188,
		12217945248258292403,
		15040591918593206314,
		10196271240168275850,
		17454674735070008619,

		16112657379421348746,
		3966659767133171969,
		6747040671494235679,
		17134635898950935218,
		13266965333363210135,
		81209830707133972,
		14565307473679759727,
		12842199825730749217,
		5153378449935477962,
		18325838682337833456,

		10026837077846575917,
		13306130168815928044,
		9529303547091196674,
		13868697402250917112,
		1077258772023233131,
		2323601749727489988,
		2445530380711524852,
		2465240108500726185,
		14738717527031843504,
		17736736152516527522,

		6622583241179933023,
		8108789634101562096,
		8723549642453320143,
		10776411678057817288,
		16196385164429935031,
		1784963395870167873,
		15937685718162689168,
		11901036560523968260,
		1610416157172001793,
		3251426045027230343,

		673530795527425458,
		12791182347560578640,
		11375670541237055652,
		7680550370207558962,
		15295227971848272658,
		18119229184320918134,
		13121531462532881593,
		14668748687827404894,
		12321670100346246080,
		16005748529364355925,

		2909453736022459933,
		16535832053716085132,
		13612139032234658489,
		11326518212174554290,
		4197694715829616929,
		8432224389058511248,
		1339605530172132612,
		11933830061953953509,
		756014781951608408,
		18379469307992717843,

		11829385583094655853,
		5772132433469037985,
		5797940506454661128,
		14288652381668190119,
		15889485693670820070,
		16427370942790530533,
		3197654589145922463,
		7114349718377395164,
		12657364929105380015,
		139630060266322000,

		16910137451424886744,
		11719668144177960815,
		14693620242351849219,
		1808041022271042016,
		6295299071064305375,
		5021229926059781750,
		14964916673564067161,
		9658174652416860999,
		17004748172314319909,
		2037735265604547924,

		287874344511560511,
		5568127458722620823,
		1018449388072224583,
		7881527805234160598,
		9702438265921857053,
		16079323455666762616,
		15631188299770605319,
		14472250979155237339,
		7541762227759514040,
		2539187201822986973,

		10925556180321046092,
		14451266577884543366,
		4397369085824640770,
		11033095764197080928,
		10939208830480024193,
		14005077035225519676,
		11839235073892302459,
		9657069711034203654,
		1184085479274584202,
		16696609534417625965,

		16278584235506052396,
		15072666732994271408,
		7142521794092945149,
		10465668294368196397,
		17692729464865584106,
		7418485618555623485,
		8712731867155519646,
		4468388687863944900,
		5083269758217375564,
		15551276155233618324,

		2345746398861397430,
		7096324336007218697,
		651526573899649510,
		1097808939115696735,
		17985312295376783999,
		7468535267301512895,
		5785076932494497461,
		9368537826769890728,
		9592208783854777205,
		14300411218192624460,

		4528685786491571648,
		15871788363065518964,
		15454095432612357237,
		3854645357595891418,
		6035617458952085994,
		17236825324546962897,
		7597516170530776953,
		7687205163750826870,
		4836309431406917522,
		9563169619186481874,

		17840054864280231575,
		7966859837131690269,
		7098876432424781020,
		14923629053432410004,
		12219406470993116856,
		740022619711853205,
		12092659273027629487,
		10407758240645343719,
		15431981157506386106,
		15465904042279203700,

		10545667340230627330,
		7731822882548751652,
		9472250891413828662,
		15954587087013074843,
		7054159474469276273,
		10365492847621727716,
		444519347741284919,
		13558608028846044533,
		7394366606010546296,
		2818373644747619373,

		2218761067536014861,
		10452497450240962053,
		3524426889480971752,
		13126799864978843109,
		17938823847930054188,
		8587962150927009207,
		8516394259998612102,
		9599001446308225731,
		6593815249931696964,
		12543633398518356794,

		2092108433483821133,
		13280547730028448909,
		9637514919996850899,
		15432339695536081898,
		6399724964882104517,
		11500759960845661496,
		10917394140008100342,
		12887550545785546146,
		15818117520938009175,
		10623842221768750106,

		12622964253012675657,
		17641249001523333480,
		11701076559072206327,
		6600698816967293092,
		18148803528224270024,
		17332582479356970806,
		6803065734153993695,
		11888050508687878959,
		6985904275465775259,
		10813103619628126283,

		1386958442866091979,
		3306012345502884736,
		11194322863207048292,
		15346396588418300491,
		9841489941469931493,
		17348022635526859794,
		299395649767828741,
		11943663038229771336,
		12155522700778574957,
		9684416936043483621,

		18323750784176662933,
		131308565171207160,
		2397541703904308897,
		1084327145434310779,
		13614587060597279858,
		13553762890387170542,
		2413267854240329373,
		12564142441156513615,
		6192406703792981940,
		6029166095840844139,

		9714355823629205772,
		6937252457464419013,
		13709196842371127179,
		3976030592487376259,
		3967851955116048111,
		14047717257355598245,

		11732085503602182596,
		18172949853770195508,

		16782829275752129712,
		13184226126879662090,
		11036460738217411734,
		3419610213944540596,
		11656522751314698381,
		15532927260031183208,
		5415046934104928462,

		13651664417187865973,
		14701801093277991628,
		17322299166289939582,
		8684483421994390782,
		11605472013584712748,
		17662299908984419112,
		11079590939262292356,
		6821970413897739901,
		18227622685994316776,
		2254786583603744796,

		4782004756683262167,
		1856858900781954020,
		15621412651120287780,
		3141525561741908975,
		18175703865326171721,
		17300884541990962365,
		14805830295725077458,
		14545803645200510433,
		13759681792925871338,
		1194045271593288128,

		13959718493540001839,
		17296434915078399891,
		14684929705063356684,
		14568813575652591198,
		11426767295900172124,
		15421927794069623929,
		15980300318440598608,
		10261352515687373745,
		4412162716568778357,
		2799325577224813821,

		16421399880013481766,
		2858752662119683455,
		9721932782903050917,
		12817375793870570853,
		1017797933567612814,
		15586913700653956050,
		3766853172264639544,
		4019669011507595237,
		7194336154395361171,
		5177833393476979204,

		9192603983782913483,
		4569261685979410777,
		5080866903808682902,
		7438451232804004876,
		575775830774922757,
		12104511505942878534,
		3064795205495094737,
		3036961561955890352,
		17526786022695695328,
		17645078820449732023,

		9192603983782913483,
		4569261685979410777,
		5080866903808682902,
		7438451232804004876,
		575775830774922757,
		12104511505942878534,
		3064795205495094737,
		3036961561955890352,
		17526786022695695328,
		17645078820449732023,

		16026363501325776268,
		1401655648928450405,
		1357667158480071432,
		5648464153296393663,
		1208401122190112763,

		4851252989191015378,
		17794579459843113812,
		7448786408250198337,
		8064036078157062016,
		6360213594564303682,
		2268955710733151026,
		6366252696374168690,
		10208116736897202708,
		7269161525823858368,
		6482957493043318799,

		8726617143761384940,
		1033955166307540957,
		7104910685588482610,
		5600786041813697757,
		11040151803952403463,
		13814892673077448720,
		14429525009898960257,
		3873071314407144875,
		14972567838121101789,
		4407089038656246116,
	}

	for _, id := range txIDs {
		printTest(t, id, buf)
	}

	buf.Write(([]byte(fileFooter)))
	f.Write(buf.Bytes())
}

func printTest(t *template.Template, txID uint64, buf io.Writer) {
	resp, err := req.Get(getBytesURI + strconv.FormatUint(txID, 10))
	if err != nil {
		log.Println(txID, err)
		return
	}
	var txBytes TransactionBytes
	err = json.Unmarshal(resp.Bytes(), &txBytes)
	if err != nil {
		log.Println(txID, err)
		return
	}

	resp, err = req.Get(getTxURI + strconv.FormatUint(txID, 10))
	if err != nil {
		log.Println(txID, err)
		return
	}
	var tx Transaction
	err = json.Unmarshal(resp.Bytes(), &tx)
	if err != nil {
		log.Println(txID, err)
		return
	}

	tx.TxByteStr = txBytes.TransactionBytes

	senderPublicKeyBytes, _ := hex.DecodeString(tx.SenderPublicKey)
	signatureBytes, _ := hex.DecodeString(tx.Signature)
	if len(signatureBytes) == 0 {
		signatureBytes = make([]byte, 64)
	}
	referencedTransactionFullHashBytes, _ := hex.DecodeString(tx.ReferencedTransactionFullHash)
	if len(referencedTransactionFullHashBytes) == 0 {
		referencedTransactionFullHashBytes = make([]byte, 32)
	}

	tx.SignatureByteStr = toByteStr(signatureBytes)
	tx.SenderPublicKeyByteStr = toByteStr(senderPublicKeyBytes)
	tx.ReferencedTransactionFullHashByteStr = toByteStr(referencedTransactionFullHashBytes)

	tx.SubtypeAndVersion = tx.Subtype&0x0F | (tx.Version<<4)&0xF0

	t.Execute(buf, tx)
}

func toByteStr(bs []byte) string {
	s := ""
	for i, b := range bs {
		if i == len(bs)-1 {
			s += fmt.Sprintf("%d", b)
		} else {
			s += fmt.Sprintf("%d, ", b)
		}
	}
	return s
}