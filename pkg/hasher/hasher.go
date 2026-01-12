package hasher

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jxskiss/base62"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type ArgonConfig struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type JWTManager struct {
	algorithm     string
	publicKey     interface{}
	privateKey    interface{}
	signingMethod jwt.SigningMethod
}

type Hasher struct {
	jwtManager *JWTManager
	argon      *ArgonConfig
}

type DecodeOpt func(*decopts)

type decopts struct {
	subject *string
	expire  bool
}

func defaultDecOpts() *decopts {
	return &decopts{
		subject: nil,
		expire:  true,
	}
}

func WithSubject(subject string) DecodeOpt {
	return func(o *decopts) {
		o.subject = &subject
	}
}

func WithExpire(expire bool) DecodeOpt {
	return func(o *decopts) {
		o.expire = expire
	}
}

func NewHasher(jwtAlgorithm string, jwtPublicKey string, jwtPrivateKey string) (*Hasher, error) {
	jm, err := getJWTManager(jwtAlgorithm, jwtPublicKey, jwtPrivateKey)
	if err != nil {
		return nil, err
	}
	return &Hasher{
		jwtManager: jm,
		argon: &ArgonConfig{
			memory:      64 * 1024, //nolint:mnd // based on argon defaults
			iterations:  3,         //nolint:mnd // based on argon defaults
			parallelism: 2,         //nolint:mnd // based on argon defaults
			saltLength:  16,        //nolint:mnd // based on argon defaults
			keyLength:   32,        //nolint:mnd // based on argon defaults
		},
	}, nil
}

func (h *Hasher) HashArgon(password string) (string, error) {
	salt, err := h.RandomBytes(16) //nolint:mnd //standard salt size
	if err != nil {
		return "", err
	}

	argonHash := argon2.IDKey(
		[]byte(password),
		salt,
		h.argon.iterations,
		h.argon.memory,
		h.argon.parallelism,
		h.argon.keyLength,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(argonHash)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.argon.memory,
		h.argon.iterations,
		h.argon.parallelism,
		b64Salt,
		b64Hash,
	)
	return encodedHash, nil
}

func (h *Hasher) CompareArgon(password string, encodedHash string) (bool, error) {
	a, salt, hash, err := h.decodeHashArgon(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		a.iterations,
		a.memory,
		a.parallelism,
		a.keyLength,
	)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (h *Hasher) UUIDVv7() (uuid.UUID, error) {
	return uuid.NewV7()
}

func (h *Hasher) UUIDVv4() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func (h *Hasher) decodeHashArgon(encodedHash string) (*ArgonConfig, []byte, []byte, error) {
	var (
		a    *ArgonConfig
		err  error
		salt []byte
		hash []byte
	)
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 { //nolint:mnd // number argon parameters
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	a = &ArgonConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &a.memory, &a.iterations, &a.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	a.saltLength = uint32(len(salt)) //nolint:gosec // salt length is fixed

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	a.keyLength = uint32(len(hash)) //nolint:gosec // hash length is fixed

	return a, salt, hash, nil
}

func (h *Hasher) RandomBytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (h *Hasher) RandomString(length uint) (string, error) {
	var realLength uint
	if length%2 == 0 {
		realLength = length / 2 //nolint:mnd // 1 symbol = 2 bytes
	} else {
		realLength = (length / 2) + 1 //nolint:mnd // 1 symbol = 2 bytes
	}
	bytes, err := h.RandomBytes(realLength)
	if err != nil {
		return "", err
	}
	return h.HexString(bytes)[:length], nil
}

func (h *Hasher) HexString(bytes []byte) string {
	return base62.EncodeToString(bytes)
}

func (h *Hasher) HexEncodeString(str string) string {
	return h.HexString([]byte(str))
}

func (h *Hasher) HexDecodeString(hex string) (string, error) {
	b, err := base62.DecodeString(hex)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (h *Hasher) NewHexID() (string, error) {
	ms := time.Now().UnixMilli()
	tsBytes := make([]byte, 4)                               //nolint:mnd // reserved 4 bytes
	binary.BigEndian.PutUint32(tsBytes, uint32(ms&0xFFFFFF)) //nolint:mnd // timestamp

	entropy, err := h.RandomBytes(7) //nolint:mnd // generate 7 bytes
	if err != nil {
		return "", err
	}
	combined := append(tsBytes, entropy...)

	return base62.EncodeToString(combined), nil
}

// Encode time.Now().Add(time.Second * time.Duration(exp)).Unix().
func (h *Hasher) EncodeJWT(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims(payload)
	token := jwt.NewWithClaims(h.jwtManager.signingMethod, claims)
	privateKey := h.jwtManager.privateKey
	return token.SignedString(privateKey)
}

func (h *Hasher) DecodeJWT(token string, opts ...DecodeOpt) (map[string]interface{}, error) {
	var (
		claims jwt.MapClaims
		sub    string
		err    error
	)

	o := defaultDecOpts()
	for _, opt := range opts {
		opt(o)
	}

	tokenData, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != h.jwtManager.signingMethod {
			return nil, errors.New("invalid token algorithm")
		}
		return h.jwtManager.publicKey, nil
	})

	if tokenData == nil || !tokenData.Valid || err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) && !o.expire {
			return claims, nil
		}
		return nil, errors.New("invalid token")
	}

	if o.subject != nil {
		sub, err = tokenData.Claims.GetSubject()

		if err != nil {
			return nil, errors.New("invalid token subject error")
		}

		if *o.subject != sub {
			return nil, errors.New("invalid token subject")
		}
	}

	return claims, nil
}

func getJWTManager(jwtAlgorithm string, jwtPublicKey string, jwtPrivateKey string) (*JWTManager, error) {
	var err error
	jm := &JWTManager{
		algorithm: jwtAlgorithm,
	}
	switch jwtAlgorithm {
	case "RS512":
		jm.signingMethod = jwt.SigningMethodRS512
		jm.publicKey, jm.privateKey, err = loadRSAKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "RS256":
		jm.signingMethod = jwt.SigningMethodRS256
		jm.publicKey, jm.privateKey, err = loadRSAKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "ES512":
		jm.signingMethod = jwt.SigningMethodES512
		jm.publicKey, jm.privateKey, err = loadECKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "ES256":
		jm.signingMethod = jwt.SigningMethodES256
		jm.publicKey, jm.privateKey, err = loadECKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid algorithm")
	}
	return jm, nil
}

func loadRSAKeys(jwtPublicKey string, jwtPrivateKey string) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	var (
		privateKeyBytes []byte
		publicKeyBytes  []byte
		publicKey       *rsa.PublicKey
		privateKey      *rsa.PrivateKey
		err             error
	)
	privateKeyBytes, err = base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

func loadECKeys(jwtPublicKey string, jwtPrivateKey string) (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	var (
		privateKeyBytes []byte
		publicKeyBytes  []byte
		publicKey       *ecdsa.PublicKey
		privateKey      *ecdsa.PrivateKey
		err             error
	)
	privateKeyBytes, err = base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}
