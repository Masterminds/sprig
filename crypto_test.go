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
	Username      string
	Password      string
	HashAlgorithm HashAlgorithm
	Valid         bool
}

func TestHtpasswd(t *testing.T) {
	expectations := []HtpasswdCred{
		{Username: "myUser", Password: "myPassword", HashAlgorithm: HashBCrypt, Valid: true},
		{Username: "special'o79Cv_*qFe,)<user", Password: "special<j7+3p#6-.Jx2U:m8G;kGypassword", HashAlgorithm: HashBCrypt, Valid: true},
		{Username: "wrongus:er", Password: "doesn'tmatter", HashAlgorithm: HashBCrypt, Valid: false}, // ':' isn't allowed in the username - https://tools.ietf.org/html/rfc2617#page-6
		{Username: "mySahUser", Password: "myShaPassword", HashAlgorithm: HashSHA, Valid: true},
		{Username: "myDefaultUser", Password: "defaulthashpass", Valid: true},
	}

	for _, credential := range expectations {
		out, err := runRaw(`{{htpasswd .Username .Password .HashAlgorithm}}`, credential)
		if err != nil {
			t.Error(err)
		}
		result := strings.Split(out, ":")
		if 0 != strings.Compare(credential.Username, result[0]) && credential.Valid {
			t.Error("Generated username did not match for:", credential.Username)
		}
		switch credential.HashAlgorithm {
		case HashSHA:
			if strings.TrimPrefix(result[1], "{SHA}") != hashSha(credential.Password) {
				t.Error("Generated hash is not the equivalent for password:", credential.Password)
			}
		default:
			if bcrypt_lib.CompareHashAndPassword([]byte(result[1]), []byte(credential.Password)) != nil && credential.Valid {
				t.Error("Generated hash is not the equivalent for password:", credential.Password)
			}
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

func TestDerivePublicKey(t *testing.T) {
	tpl := `{{genPrivateKey "rsa" | derivePublicKey}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "PUBLIC KEY") {
		t.Error("Expected PUBLIC KEY", out)
	}
	tpl = `{{genPrivateKey "dsa" | derivePublicKey}}`
	out, err = runRaw(tpl, nil)
	// x509.MarshalPKIXPublicKey() does not support DSA keys
	if err == nil || !strings.Contains(err.Error(), "x509: unsupported public key type") {
		t.Error("Expected error to contain 'x509: unsupported public key type'", err)
	}
	tpl = `{{genPrivateKey "ecdsa" | derivePublicKey}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "PUBLIC KEY") {
		t.Error("Expected PUBLIC KEY", out)
	}
	tpl = `{{genPrivateKey "ed25519" | derivePublicKey}}`
	out, err = runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(out, "PUBLIC KEY") {
		t.Error("Expected PUBLIC KEY", out)
	}
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	expectedPublicKey := `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAjrVUJpXK13XN7o+B1OdP
rOWt8xpDld3q9GX5f7HlxF98KDBrLzxLcI4nFE1XriEZSCitG7cSC9jvoiU4yA59
t+fIcmU5fLAwBmkNmWnTPD2YkH7GauenL2+LaGQ/6oc3/KqhHACQ7Sj+tzOwkMiv
hw7MMdrP1NEXsYQw1ht/o38EcdRf+G9w4d/YU7aKIxM7XX3evDFpada7RhBsMXoO
tHA/mE4KuztIFZ6e2McB+4fVNPsYN/k9f5ta/iCBdOkG1WdZvDj7KZfLUWno6emD
3oE6I1crrXeuz/tabjHuQoWhxCV2OZStDMdxFg1rAjZBsq9325kO3N6PiH+pyHrd
kRZvbQBiFjlJBa+/YMJzU3dDjpCxVBGlIXuT4/22JhdPBxvGwRx9ZLKup2qbfkCx
tquWMCQN+7SE3mNXxrGxBfMsFtCgVQvkVuDaGhiYQ98DgmR/sXSZ/0okWWIockoX
WOrnMrXzvhMkF4zsd1CqhF6ctN4SYADCiR2VmeN2brzB3JZHh97DHWDlkWmDALSy
Izka65Tg7xaEYvAluaKAkbjMYBP9t9rSE3Z5w4BkwTypAnyJkOd4NwzGon3OhEjD
zHXuCEHwBPZgEcNFbYS1kezDB7QkuzbyTZvDhY+jOnJZuD6JxDLox20LBEGB1dFH
cn2WwZpyAtzWCqPvtFqKpA0CAwEAAQ==
-----END PUBLIC KEY-----
`
	tpl = `{{derivePublicKey .}}`
	out, err = runRaw(tpl, testPrivateKey)
	if err != nil {
		t.Error(err)
	}
	if out != expectedPublicKey {
		t.Error("Got incorrect public key", out)
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
