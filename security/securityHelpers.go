package security

import (
	ed "crypto/ed25519"
	"student20_pop/lib"
	"student20_pop/parser"
)

const MaxPropagationDelay = 600
const MaxClockDifference = 100

/*
	we check that Sign(sender||data) is the given signature
*/
func VerifySignature(publicKey string, data []byte, signature string) error {
	//check the size of the key as it will panic if we plug it in Verify
	if len(publicKey) != ed.PublicKeySize {
		return lib.ErrRequestDataInvalid
	}
	if ed.Verify([]byte(publicKey), data, []byte(signature)) {
		return nil
	}
	//invalid signature
	return lib.ErrRequestDataInvalid
}

/*
	handling of dynamic updates with object as item and not just string
	*publicKeys is already decoded
    *sender and signature are not already decoded
*/
func VerifyWitnessSignatures(authorizedWitnesses []string, witnessSignaturesEnc []string, sender string) error {
	senderDecoded, err := lib.Decode(sender)
	if err != nil {
		return lib.ErrEncodingFault
	}
	//TODO verify witnesses are in event's witness list (@ouriel)
	for i := 0; i < len(witnessSignaturesEnc); i++ {
		witnessSignatures, err := parser.ParseWitnessSignature(witnessSignaturesEnc[i])
		if err != nil {
			return err
		}
		//right now we apply the first option and publickeys is then usless here
		err = VerifySignature(witnessSignatures.Witness, senderDecoded, witnessSignatures.Signature)
		if err != nil {
			return err
		}
	}
	return nil
}