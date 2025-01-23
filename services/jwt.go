package services

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var TestPublicKey ecdsa.PublicKey

// Claims documentation
// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-the-id-token.html

type Identities struct {
	dateCreated  int
	userId       int
	providerName string
	providerType string
	issuer       string
	primary      bool
}

type Payload struct {
	Sub            string
	Email_verified string
	Identities     string
	Email          string
	Username       string
	Exp            int
	Iss            string
}

type Header struct {
	Typ    string
	Kid    string
	Alg    string
	Iss    string
	Client string
	Signer string
	Exp    int
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegmentAWS(seg string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(seg)
}

func VerifyAWS(m *jwt.SigningMethodECDSA, signingString string, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = DecodeSegmentAWS(signature); err != nil {
		return err
	}

	// Get the key
	var ecdsaKey *ecdsa.PublicKey
	switch k := key.(type) {
	case *ecdsa.PublicKey:
		ecdsaKey = k
	default:
		return jwt.ErrInvalidKeyType
	}

	if len(sig) != 2*m.KeySize {
		return jwt.ErrECDSAVerification
	}

	r := big.NewInt(0).SetBytes(sig[:m.KeySize])
	s := big.NewInt(0).SetBytes(sig[m.KeySize:])

	// Create hasher
	if !m.Hash.Available() {
		return jwt.ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Verify the signature
	if verifystatus := ecdsa.Verify(ecdsaKey, hasher.Sum(nil), r, s); verifystatus {
		return nil
	}

	return jwt.ErrECDSAVerification
}

func ParsePayload(rawPayload string) (Payload, error) {
	var payload Payload
	p, err := base64.URLEncoding.DecodeString(rawPayload)
	if err != nil {
		return payload, fmt.Errorf("error decoding payload: %s", err)
	}
	// fmt.Printf("decoded payload: %s\n", string(p))
	err = json.Unmarshal(p, &payload)
	if err != nil {
		return payload, fmt.Errorf("err unmarshalling payload: %s", err)
	}

	return payload, nil
}

func GetPublicKey(kid string) (any, error) {
	env := os.Getenv("env")

	if env == "dev" {
		return &TestPublicKey, nil
	}
	if env == "prod" {
		publicKeyUrl := fmt.Sprintf("https://public-keys.auth.elb.us-west-1.amazonaws.com/%s", kid)
		res, err := http.Get(publicKeyUrl)
		if err != nil {
			return "", fmt.Errorf("error getting public key: %s", err)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("error reading body from public key request: %s", err)
		}
		res.Body.Close()

		block, _ := pem.Decode(body)
		if block == nil {
			return "", fmt.Errorf("error: failed to parse PEM block containing public key")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("error parsing public key block.Bytes: %s", err)
		}
		return pub, nil
	}

	return nil, fmt.Errorf("environment not specified")
}

func ParseAndVerifyJWT(amznData string) (Payload, error) {
	parts := strings.Split(amznData, ".")
	rawHeader := parts[0]
	rawPayload := parts[1]
	rawSignature := parts[2]
	claims := Payload{}

	// fmt.Printf("rawHeader: %s\n", rawHeader)
	// fmt.Printf("rawPayload: %s\n", rawPayload)
	// fmt.Printf("rawSignature: %s\n", rawSignature)
	//
	h, err := base64.URLEncoding.DecodeString(rawHeader)
	if err != nil {
		return claims, fmt.Errorf("error decoding header: %s", err)
	}
	// fmt.Printf("header: %s\n", string(h))

	var header jwt.MapClaims

	err = json.Unmarshal(h, &header)
	if err != nil {
		return claims, fmt.Errorf("err unmarshalling header: %s", err)
	}

	kid, ok := header["kid"].(string)
	if !ok {
		return claims, fmt.Errorf("kid not of type string")
	}

	pub, err := GetPublicKey(kid)
	if err != nil {
		return Payload{}, err
	}

	if !ok {
		return claims, fmt.Errorf("error extracting 'kid' from header: %s", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return claims, fmt.Errorf("error pub is not of type ecdsa.PublicKey: %s", err)
	}

	signingString := fmt.Sprintf("%s.%s", rawHeader, rawPayload)
	signingMethod := jwt.SigningMethodES256

	err = VerifyAWS(signingMethod, signingString, rawSignature, ecdsaPub)
	if err != nil {
		return claims, fmt.Errorf("error verifying jwt: %s", err)
	} else {
		payload, err := ParsePayload(rawPayload)
		if err != nil {
			return claims, err
		}
		return payload, nil
	}
}

type TestCustomClaims struct {
	Foo string `json:"foo"`
	Kid string `json:"kid"`
	jwt.RegisteredClaims
}

func CreateTestJWT(cognitoId string) (string, error) {
	claims := TestCustomClaims{
		"bar",
		"test",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   cognitoId,
			Audience:  []string{"test-server"},
			ID:        "1",
		},
	}

	TestPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", err
	}
	TestPublicKey = TestPrivateKey.PublicKey
	token := AWSNewWithClaims(jwt.SigningMethodES256, claims)
	return token.AWSSignedString(TestPrivateKey)
}

func (t *AWSToken) AWSSignedString(key interface{}) (string, error) {
	sstr, err := t.AWSSigningString()
	if err != nil {
		return "", err
	}

	sig, err := t.Method.Sign(sstr, key)
	if err != nil {
		return "", err
	}

	return sstr + "." + t.AWSEncodeSegment([]byte(sig)), nil
}

type TokenOption func(*AWSToken)

func AWSNewWithClaims(method jwt.SigningMethod, claims jwt.Claims, opts ...TokenOption) *AWSToken {
	return &AWSToken{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			"kid": "test",
		},
		Claims: claims,
		Method: method,
	}
}

/* An extension of the jwt.Token type that can hold Kid in its header in order to match AWS tokens */
type AWSToken struct {
	Raw       string                 // Raw contains the raw token.  Populated when you [Parse] a token
	Method    jwt.SigningMethod      // Method is the signing method used or to be used
	Header    map[string]interface{} // Header is the first segment of the token in decoded form
	Claims    jwt.Claims             // Claims is the second segment of the token in decoded form
	Signature []byte                 // Signature is the third segment of the token in decoded form.  Populated when you Parse a token
	Valid     bool                   // Valid specifies if the token is valid.  Populated when you Parse/Verify a token
	kid       string                 /* Added additional field */
}

func (t *AWSToken) AWSSigningString() (string, error) {
	h, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}

	c, err := json.Marshal(t.Claims)
	if err != nil {
		return "", err
	}

	return t.AWSEncodeSegment(h) + "." + t.AWSEncodeSegment(c), nil
}

func (*AWSToken) AWSEncodeSegment(seg []byte) string {
	return base64.URLEncoding.EncodeToString(seg)
}
