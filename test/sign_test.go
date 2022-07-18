package main

import (
	"fmt"
	"github.com/dontpanicdao/caigo"
	"math/big"
	"testing"
)

type Message struct {
	Message string
}

func ToASCII(str string) string {
	runes := []rune(str)
	result := "0x"

	for i := 0; i < len(runes); i++ {
		result = result + fmt.Sprintf("%x", caigo.StrToBig(fmt.Sprint(int(runes[i]))))
	}
	return result
}

func (message Message) FmtDefinitionEncoding(field string) (fmtEnc []*big.Int) {
	if field == "message" {
		fmtEnc = append(fmtEnc, caigo.HexToBN(ToASCII(message.Message)))
	}
	return fmtEnc
}

func TypedData() (ttd caigo.TypedData) {
	exampleTypes := make(map[string]caigo.TypeDef)
	domDefs := []caigo.Definition{{"name", "felt"}, {"version", "felt"}, {"chainId", "felt"}}
	exampleTypes["StarkNetDomain"] = caigo.TypeDef{Definitions: domDefs}

	msgDefs := []caigo.Definition{{"message", "felt"}}
	exampleTypes["Message"] = caigo.TypeDef{Definitions: msgDefs}

	dm := caigo.Domain{
		Name:    "Example DApp",
		Version: "1",
		ChainId: 1, //
	}
	ttd, _ = caigo.NewTypedData(exampleTypes, "Message", dm)
	return ttd
}

func MsgHash() *big.Int {
	ttd := TypedData()

	msg := Message{
		Message: "test msg",
	}
	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}
	hash, _ := ttd.GetMessageHash(caigo.HexToBN("0x02039385c3fc65cfd45e06e78d257d52e3141e590f253c1c5be09bb9dad24b5c"), msg, curve)
	return hash
}

func TestASCII(t *testing.T) {
	str := "test msg"
	result := ToASCII(str)
	fmt.Println(result)
}

func TestTypeData(t *testing.T) {
	ttd := TypedData()

	msg := Message{
		Message: "test msg",
	}
	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}

	hash, err := ttd.GetMessageHash(caigo.HexToBN("0x02039385c3fc65cfd45e06e78d257d52e3141e590f253c1c5be09bb9dad24b5c"), msg, curve)
	hashArgentX := caigo.HexToBN("0x1216c177e0f98eaeda76a862c95a1b6028799b81b0b375f1b20e495063f0c80")

	fmt.Println("result *******************: ", hashArgentX)
	fmt.Println("hash: ", hash)
}

func TestVerify(t *testing.T) {
	ttd := TypedData()

	msg := Message{
		Message: "test msg",
	}
	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}

	hash, err := ttd.GetMessageHash(caigo.HexToBN("0x02039385c3fc65cfd45e06e78d257d52e3141e590f253c1c5be09bb9dad24b5c"), msg, curve)

	xStr := "1674666556728092608470336264650972261738561385251855821167604540120403694597"
	x := caigo.StrToBig(xStr)
	y := curve.GetYCoordinate(x)
	r := caigo.StrToBig("1000268350070957364344744286374819942330583452825450880273174518864520215865")
	s := caigo.StrToBig("2384134264167827274147140351381653026119943416245540545525304309238102630185")
	result := curve.Verify(hash, r, s, x, y)

	fmt.Println("result *******************: ", result)
	fmt.Println("hash: ", caigo.HexToBN("0x3a301c8bbe54830ceeb6caf7abe387749076e0a03efbcccc9348b0696fc8af8"))
}
func TestVerify1(t *testing.T) {

	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}

	hashArgentX := caigo.HexToBN("0x3a301c8bbe54830ceeb6caf7abe387749076e0a03efbcccc9348b0696fc8af8")

	xStr := "1674666556728092608470336264650972261738561385251855821167604540120403694597"
	x := caigo.StrToBig(xStr)
	y := curve.GetYCoordinate(x)
	r := caigo.StrToBig("3169396788236384945685856549533061799664702482332597452051824506506688700103")
	s := caigo.StrToBig("1816211193520515209228764487189767422614799904040204711934934720647499301994")
	result := curve.Verify(hashArgentX, r, s, x, y)

	fmt.Println("result *******************: ", result)
	fmt.Println("hash: ", caigo.HexToBN("0x3a301c8bbe54830ceeb6caf7abe387749076e0a03efbcccc9348b0696fc8af8"))
}

func TestSign(t *testing.T) {
	// NOTE: when not given local file path this pulls the curve data from Starkware github repo
	curve, err := caigo.SC(caigo.WithConstants("./pedersen_params.json"))
	if err != nil {
		panic(err.Error())
	}

	hash := MsgHash()

	priv, _ := curve.GetRandomPrivateKey()

	x, y, err := curve.PrivateToPoint(priv)
	if err != nil {
		panic(err.Error())
	}

	r, s, err := curve.Sign(hash, priv)
	if err != nil {
		panic(err.Error())
	}

	if curve.Verify(hash, r, s, x, y) {
		fmt.Println("signature is valid")
	} else {
		fmt.Println("signature is invalid")
	}
}
