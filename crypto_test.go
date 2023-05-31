package sprig

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
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

type testVector struct {
	password string
	salt     string
	iter     int
	output   []byte
}

// Test vectors from RFC 6070, http://tools.ietf.org/html/rfc6070
var sha1TestVectors = []testVector{
	{
		"password",
		"salt",
		1,
		[]byte{
			0x0c, 0x60, 0xc8, 0x0f, 0x96, 0x1f, 0x0e, 0x71,
			0xf3, 0xa9, 0xb5, 0x24, 0xaf, 0x60, 0x12, 0x06,
			0x2f, 0xe0, 0x37, 0xa6,
		},
	},
	{
		"password",
		"salt",
		2,
		[]byte{
			0xea, 0x6c, 0x01, 0x4d, 0xc7, 0x2d, 0x6f, 0x8c,
			0xcd, 0x1e, 0xd9, 0x2a, 0xce, 0x1d, 0x41, 0xf0,
			0xd8, 0xde, 0x89, 0x57,
		},
	},
	{
		"password",
		"salt",
		4096,
		[]byte{
			0x4b, 0x00, 0x79, 0x01, 0xb7, 0x65, 0x48, 0x9a,
			0xbe, 0xad, 0x49, 0xd9, 0x26, 0xf7, 0x21, 0xd0,
			0x65, 0xa4, 0x29, 0xc1,
		},
	},
	{
		"passwordPASSWORDpassword",
		"saltSALTsaltSALTsaltSALTsaltSALTsalt",
		4096,
		[]byte{
			0x3d, 0x2e, 0xec, 0x4f, 0xe4, 0x1c, 0x84, 0x9b,
			0x80, 0xc8, 0xd8, 0x36, 0x62, 0xc0, 0xe4, 0x4a,
			0x8b, 0x29, 0x1a, 0x96, 0x4c, 0xf2, 0xf0, 0x70,
			0x38,
		},
	},
	{
		"pass\000word",
		"sa\000lt",
		4096,
		[]byte{
			0x56, 0xfa, 0x6a, 0xa7, 0x55, 0x48, 0x09, 0x9d,
			0xcc, 0x37, 0xd7, 0xf0, 0x34, 0x25, 0xe0, 0xc3,
		},
	},
}

// Test vectors from
// http://stackoverflow.com/questions/5130513/pbkdf2-hmac-sha2-test-vectors
var sha256TestVectors = []testVector{
	{
		"password",
		"salt",
		1,
		[]byte{
			0x12, 0x0f, 0xb6, 0xcf, 0xfc, 0xf8, 0xb3, 0x2c,
			0x43, 0xe7, 0x22, 0x52, 0x56, 0xc4, 0xf8, 0x37,
			0xa8, 0x65, 0x48, 0xc9,
		},
	},
	{
		"password",
		"salt",
		2,
		[]byte{
			0xae, 0x4d, 0x0c, 0x95, 0xaf, 0x6b, 0x46, 0xd3,
			0x2d, 0x0a, 0xdf, 0xf9, 0x28, 0xf0, 0x6d, 0xd0,
			0x2a, 0x30, 0x3f, 0x8e,
		},
	},
	{
		"password",
		"salt",
		4096,
		[]byte{
			0xc5, 0xe4, 0x78, 0xd5, 0x92, 0x88, 0xc8, 0x41,
			0xaa, 0x53, 0x0d, 0xb6, 0x84, 0x5c, 0x4c, 0x8d,
			0x96, 0x28, 0x93, 0xa0,
		},
	},
	{
		"passwordPASSWORDpassword",
		"saltSALTsaltSALTsaltSALTsaltSALTsalt",
		4096,
		[]byte{
			0x34, 0x8c, 0x89, 0xdb, 0xcb, 0xd3, 0x2b, 0x2f,
			0x32, 0xd8, 0x14, 0xb8, 0x11, 0x6e, 0x84, 0xcf,
			0x2b, 0x17, 0x34, 0x7e, 0xbc, 0x18, 0x00, 0x18,
			0x1c,
		},
	},
	{
		"pass\000word",
		"sa\000lt",
		4096,
		[]byte{
			0x89, 0xb6, 0x9d, 0x05, 0x16, 0xf8, 0x29, 0x89,
			0x3c, 0x69, 0x62, 0x26, 0x65, 0x0a, 0x86, 0x87,
		},
	},
}

// Test vectors from
// https://github.com/Anti-weakpasswords/PBKDF2-Test-Vectors/releases/tag/1.0
var sha224TestVectors = []testVector{
	{
		"passDATAb00AB7YxDTTlRH2dqxD",
		"saltKEYbcTcXHCBxtjD2PnBh44A",
		1,
		decodeHexString("86AB2F3D0CB39839B46DA2DD8F210915D79AD2E6F2093D155D75C8D9"),
	},
	{
		"passDATAb00AB7YxDTTlRH2dqxD",
		"saltKEYbcTcXHCBxtjD2PnBh44A",
		100000,
		decodeHexString("0ADF2D99E7FF8DBC6B1DF4382D32959021BFDACB99B796BF9089D0E3"),
	},
}

// Test vectors from
// https://github.com/Anti-weakpasswords/PBKDF2-Test-Vectors/releases/tag/1.0
var sha384TestVectors = []testVector{
	{
		"passDATAb00AB7YxDTTlRH2dqxDx19GDxDV1zFMz7E6QVqK",
		"saltKEYbcTcXHCBxtjD2PnBh44AIQ6XUOCESOhXpEp3HrcG",
		1,
		decodeHexString("0644A3489B088AD85A0E42BE3E7F82500EC18936699151A2C90497151BAC7BB69300386A5E798795BE3CEF0A3C803227"),
	},
	{
		"passDATAb00AB7YxDTTlRH2dqxDx19GDxDV1zFMz7E6QVqK",
		"saltKEYbcTcXHCBxtjD2PnBh44AIQ6XUOCESOhXpEp3HrcG",
		100000,
		decodeHexString("BF625685B48FE6F187A1780C5CB8E1E4A7B0DBD6F551827F7B2B598735EAC158D77AFD3602383D9A685D87F8B089AF30"),
	},
}

// Test vectors from
// https://github.com/Anti-weakpasswords/PBKDF2-Test-Vectors/releases/tag/1.0
var sha512TestVectors = []testVector{
	{
		"passDATAb00AB7YxDTT",
		"saltKEYbcTcXHCBxtjD",
		1,
		decodeHexString("CBE6088AD4359AF42E603C2A33760EF9D4017A7B2AAD10AF46F992C660A0B461ECB0DC2A79C2570941BEA6A08D15D6887E79F32B132E1C134E9525EEDDD744FA"),
	},
	{
		"passDATAb00AB7YxDTT",
		"saltKEYbcTcXHCBxtjD",
		100000,
		decodeHexString("ACCDCD8798AE5CD85804739015EF2A11E32591B7B7D16F76819B30B0D49D80E1ABEA6C9822B80A1FDFE421E26F5603ECA8A47A64C9A004FB5AF8229F762FF41F"),
	},
}

func decodeHexString(decode string) []byte {
	if value, err := hex.DecodeString(decode); err == nil {
		return value
	}
	return []byte(`failed to decode hex string`)
}

func testHash(t *testing.T, hashName string, vectors []testVector) {
	for i, v := range vectors {
		tpl := fmt.Sprintf("{{ pbkdf2hash \"%s\" \"%s\" %d %d \"%s\" }}", v.password, v.salt, v.iter, len(v.output), hashName)
		o, err := runBytes(tpl, map[string]string{})
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(o, v.output) {
			t.Errorf("%s %d: expected %x, got %x", hashName, i, v.output, o)
		}
	}
}

func TestPbkdf2hash(t *testing.T) {
	testHash(t, "sha1", sha1TestVectors)
	testHash(t, "sha224", sha224TestVectors)
	testHash(t, "sha256", sha256TestVectors)
	testHash(t, "sha384", sha384TestVectors)
	testHash(t, "sha512", sha512TestVectors)
}
