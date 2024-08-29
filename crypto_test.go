package sprig

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	bcrypt_lib "golang.org/x/crypto/bcrypt"
)

const (
	beginCertificate = "-----BEGIN CERTIFICATE-----"
	endCertificate   = "-----END CERTIFICATE-----"
)

var (
	// fastCertKeyAlgos is the list of private key algorithms that are supported for certificate use, and
	// are fast to generate.
	fastCertKeyAlgos = []string{
		"ecdsa",
		"ed25519",
	}
)

func TestSha512Sum(t *testing.T) {
	tpl := `{{"abc" | sha512sum}}`
	if err := runt(tpl, "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f"); err != nil {
		t.Error(err)
	}
}

func TestSha256Sum(t *testing.T) {
	tpl := `{{"abc" | sha256sum}}`
	if err := runt(tpl, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"); err != nil {
		t.Error(err)
	}
}

func TestSha1Sum(t *testing.T) {
	tpl := `{{"abc" | sha1sum}}`
	if err := runt(tpl, "a9993e364706816aba3e25717850c26c9cd0d89d"); err != nil {
		t.Error(err)
	}
}

func TestAdler32Sum(t *testing.T) {
	tpl := `{{"abc" | adler32sum}}`
	if err := runt(tpl, "38600999"); err != nil {
		t.Error(err)
	}
}

func TestBcrypt(t *testing.T) {
	out, err := runRaw(`{{"abc" | bcrypt}}`, nil)
	if err != nil {
		t.Error(err)
	}
	if bcrypt_lib.CompareHashAndPassword([]byte(out), []byte("abc")) != nil {
		t.Error("Generated hash is not the equivalent for password:", "abc")
	}
}

type HtpasswdCred struct {
	Username string
	Password string
	Valid    bool
}

func TestHtpasswd(t *testing.T) {
	expectations := []HtpasswdCred{
		{Username: "myUser", Password: "myPassword", Valid: true},
		{Username: "special'o79Cv_*qFe,)<user", Password: "special<j7+3p#6-.Jx2U:m8G;kGypassword", Valid: true},
		{Username: "wrongus:er", Password: "doesn'tmatter", Valid: false}, // ':' isn't allowed in the username - https://tools.ietf.org/html/rfc2617#page-6
	}

	for _, credential := range expectations {
		out, err := runRaw(`{{htpasswd .Username .Password}}`, credential)
		if err != nil {
			t.Error(err)
		}
		result := strings.Split(out, ":")
		if 0 != strings.Compare(credential.Username, result[0]) && credential.Valid {
			t.Error("Generated username did not match for:", credential.Username)
		}
		if bcrypt_lib.CompareHashAndPassword([]byte(result[1]), []byte(credential.Password)) != nil && credential.Valid {
			t.Error("Generated hash is not the equivalent for password:", credential.Password)
		}
	}
}

func TestDerivePassword(t *testing.T) {
	expectations := map[string]string{
		`{{derivePassword 1 "long" "password" "user" "example.com"}}`:    "ZedaFaxcZaso9*",
		`{{derivePassword 2 "long" "password" "user" "example.com"}}`:    "Fovi2@JifpTupx",
		`{{derivePassword 1 "maximum" "password" "user" "example.com"}}`: "pf4zS1LjCg&LjhsZ7T2~",
		`{{derivePassword 1 "medium" "password" "user" "example.com"}}`:  "ZedJuz8$",
		`{{derivePassword 1 "basic" "password" "user" "example.com"}}`:   "pIS54PLs",
		`{{derivePassword 1 "short" "password" "user" "example.com"}}`:   "Zed5",
		`{{derivePassword 1 "pin" "password" "user" "example.com"}}`:     "6685",
	}

	for tpl, result := range expectations {
		out, err := runRaw(tpl, nil)
		if err != nil {
			t.Error(err)
		}
		if 0 != strings.Compare(out, result) {
			t.Error("Generated password does not match for", tpl)
		}
	}
}

// NOTE(bacongobbler): this test is really _slow_ because of how long it takes to compute
// and generate a new crypto key.
func TestGenPrivateKey(t *testing.T) {
	// test that calling by default generates an RSA private key
	tpl := `{{genPrivateKey ""}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "RSA PRIVATE KEY") {
		t.Error("Expected RSA PRIVATE KEY")
	}
	// test all acceptable arguments
	tpl = `{{genPrivateKey "rsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "RSA PRIVATE KEY") {
		t.Error("Expected RSA PRIVATE KEY")
	}
	tpl = `{{genPrivateKey "dsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "DSA PRIVATE KEY") {
		t.Error("Expected DSA PRIVATE KEY")
	}
	tpl = `{{genPrivateKey "ecdsa"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "EC PRIVATE KEY") {
		t.Error("Expected EC PRIVATE KEY")
	}
	tpl = `{{genPrivateKey "ed25519"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "PRIVATE KEY") {
		t.Error("Expected PRIVATE KEY")
	}
	// test bad
	tpl = `{{genPrivateKey "bad"}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if out != "Unknown type bad" {
		t.Error("Expected type 'bad' to be an unknown crypto algorithm")
	}
	// ensure that we can base64 encode the string
	tpl = `{{genPrivateKey "rsa" | b64enc}}`
	_, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestRandBytes(t *testing.T) {
	tpl := `{{randBytes 12}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	bytes, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		t.Error(err)
	}
	if len(bytes) != 12 {
		t.Error("Expected 12 base64-encoded bytes")
	}

	out2, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	if out == out2 {
		t.Error("Expected subsequent randBytes to be different")
	}
}

func TestUUIDGeneration(t *testing.T) {
	tpl := `{{uuidv4}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	if len(out) != 36 {
		t.Error("Expected UUID of length 36")
	}

	out2, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	if out == out2 {
		t.Error("Expected subsequent UUID generations to be different")
	}
}

func TestBuildCustomCert(t *testing.T) {
	ca, _ := generateCertificateAuthority("example.com", 365)
	tpl := fmt.Sprintf(
		`{{- $ca := buildCustomCert "%s" "%s"}}
{{- $ca.Cert }}`,
		base64.StdEncoding.EncodeToString([]byte(ca.Cert)),
		base64.StdEncoding.EncodeToString([]byte(ca.Key)),
	)
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	tpl2 := fmt.Sprintf(
		`{{- $ca := buildCustomCert "%s" "%s"}}
{{- $ca.Cert }}`,
		base64.StdEncoding.EncodeToString([]byte("fail")),
		base64.StdEncoding.EncodeToString([]byte(ca.Key)),
	)
	out2, _ := runRaw(tpl2, nil)

	assert.Equal(t, out, ca.Cert)
	assert.NotEqual(t, out2, ca.Cert)
}

func TestGenCA(t *testing.T) {
	testGenCA(t, nil)
}

func TestGenCAWithKey(t *testing.T) {
	for _, keyAlgo := range fastCertKeyAlgos {
		t.Run(keyAlgo, func(t *testing.T) {
			testGenCA(t, &keyAlgo)
		})
	}
}

func testGenCA(t *testing.T, keyAlgo *string) {
	const cn = "foo-ca"

	var genCAExpr string
	if keyAlgo == nil {
		genCAExpr = "genCA"
	} else {
		genCAExpr = fmt.Sprintf(`genPrivateKey "%s" | genCAWithKey`, *keyAlgo)
	}

	tpl := fmt.Sprintf(
		`{{- $ca := %s "%s" 365 }}
{{ $ca.Cert }}
`,
		genCAExpr,
		cn,
	)
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, out, beginCertificate)
	assert.Contains(t, out, endCertificate)

	decodedCert, _ := pem.Decode([]byte(out))
	assert.Nil(t, err)
	cert, err := x509.ParseCertificate(decodedCert.Bytes)
	assert.Nil(t, err)

	assert.Equal(t, cn, cert.Subject.CommonName)
	assert.True(t, cert.IsCA)
}

func TestGenSelfSignedCert(t *testing.T) {
	testGenSelfSignedCert(t, nil)
}

func TestGenSelfSignedCertWithKey(t *testing.T) {
	for _, keyAlgo := range fastCertKeyAlgos {
		t.Run(keyAlgo, func(t *testing.T) {
			testGenSelfSignedCert(t, &keyAlgo)
		})
	}
}

func testGenSelfSignedCert(t *testing.T, keyAlgo *string) {
	const (
		cn   = "foo.com"
		ip1  = "10.0.0.1"
		ip2  = "10.0.0.2"
		dns1 = "bar.com"
		dns2 = "bat.com"
	)

	var genSelfSignedCertExpr string
	if keyAlgo == nil {
		genSelfSignedCertExpr = "genSelfSignedCert"
	} else {
		genSelfSignedCertExpr = fmt.Sprintf(`genPrivateKey "%s" | genSelfSignedCertWithKey`, *keyAlgo)
	}

	tpl := fmt.Sprintf(
		`{{- $cert := %s "%s" (list "%s" "%s") (list "%s" "%s") 365 }}
{{ $cert.Cert }}`,
		genSelfSignedCertExpr,
		cn,
		ip1,
		ip2,
		dns1,
		dns2,
	)

	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, out, beginCertificate)
	assert.Contains(t, out, endCertificate)

	decodedCert, _ := pem.Decode([]byte(out))
	assert.Nil(t, err)
	cert, err := x509.ParseCertificate(decodedCert.Bytes)
	assert.Nil(t, err)

	assert.Equal(t, cn, cert.Subject.CommonName)
	assert.Equal(t, 1, cert.SerialNumber.Sign())
	assert.Equal(t, 2, len(cert.IPAddresses))
	assert.Equal(t, ip1, cert.IPAddresses[0].String())
	assert.Equal(t, ip2, cert.IPAddresses[1].String())
	assert.Contains(t, cert.DNSNames, dns1)
	assert.Contains(t, cert.DNSNames, dns2)
	assert.False(t, cert.IsCA)
}

func TestGenSignedCert(t *testing.T) {
	testGenSignedCert(t, nil, nil)
}

func TestGenSignedCertWithKey(t *testing.T) {
	for _, caKeyAlgo := range fastCertKeyAlgos {
		for _, certKeyAlgo := range fastCertKeyAlgos {
			t.Run(fmt.Sprintf("%s-%s", caKeyAlgo, certKeyAlgo), func(t *testing.T) {
				testGenSignedCert(t, &caKeyAlgo, &certKeyAlgo)
			})
		}
	}
}

func testGenSignedCert(t *testing.T, caKeyAlgo, certKeyAlgo *string) {
	const (
		cn   = "foo.com"
		ip1  = "10.0.0.1"
		ip2  = "10.0.0.2"
		dns1 = "bar.com"
		dns2 = "bat.com"
	)

	var genCAExpr, genSignedCertExpr string
	if caKeyAlgo == nil {
		genCAExpr = "genCA"
	} else {
		genCAExpr = fmt.Sprintf(`genPrivateKey "%s" | genCAWithKey`, *caKeyAlgo)
	}
	if certKeyAlgo == nil {
		genSignedCertExpr = "genSignedCert"
	} else {
		genSignedCertExpr = fmt.Sprintf(`genPrivateKey "%s" | genSignedCertWithKey`, *certKeyAlgo)
	}

	tpl := fmt.Sprintf(
		`{{- $ca := %s "foo" 365 }}
{{- $cert := %s "%s" (list "%s" "%s") (list "%s" "%s") 365 $ca }}
{{ $cert.Cert }}
`,
		genCAExpr,
		genSignedCertExpr,
		cn,
		ip1,
		ip2,
		dns1,
		dns2,
	)
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}

	assert.Contains(t, out, beginCertificate)
	assert.Contains(t, out, endCertificate)

	decodedCert, _ := pem.Decode([]byte(out))
	assert.Nil(t, err)
	cert, err := x509.ParseCertificate(decodedCert.Bytes)
	assert.Nil(t, err)

	assert.Equal(t, cn, cert.Subject.CommonName)
	assert.Equal(t, 1, cert.SerialNumber.Sign())
	assert.Equal(t, 2, len(cert.IPAddresses))
	assert.Equal(t, ip1, cert.IPAddresses[0].String())
	assert.Equal(t, ip2, cert.IPAddresses[1].String())
	assert.Contains(t, cert.DNSNames, dns1)
	assert.Contains(t, cert.DNSNames, dns2)
	assert.False(t, cert.IsCA)
}

func TestEncryptDecryptAES(t *testing.T) {
	tpl := `{{"plaintext" | encryptAES "secretkey" | decryptAES "secretkey"}}`
	if err := runt(tpl, "plaintext"); err != nil {
		t.Error(err)
	}
}
