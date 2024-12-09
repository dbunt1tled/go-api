package hasher

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"go_echo/internal/util/rand"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
	p                      = &Params{ //nolint:gochecknoglobals // based on argon defaults
		memory:      64 * 1024, //nolint:mnd // based on argon defaults
		iterations:  3,         //nolint:mnd // based on argon defaults
		parallelism: 2,         //nolint:mnd // based on argon defaults
		saltLength:  16,        //nolint:mnd // based on argon defaults
		keyLength:   32,        //nolint:mnd // based on argon defaults
	}
)

func HashArgon(password string) (string, error) {
	salt, err := rand.Bytes(16) //nolint:mnd //standard salt size
	if err != nil {
		return "", err
	}

	argonHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(argonHash)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.memory,
		p.iterations,
		p.parallelism,
		b64Salt,
		b64Hash,
	)
	return encodedHash, nil
}

func CompareArgon(password string, encodedHash string) (bool, error) {
	p, salt, hash, err := decodeHashArgon(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHashArgon(encodedHash string) (*Params, []byte, []byte, error) {
	var (
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

	p = &Params{} // TODO check if is normal
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt)) //nolint:gosec // salt length is fixed

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash)) //nolint:gosec // hash length is fixed

	return p, salt, hash, nil
}

func UUIDVv7() (uuid.UUID, error) {
	return uuid.NewV7()
}

func UUIDVv4() (uuid.UUID, error) {
	return uuid.NewRandom()
}

func XXHash64(s string) string {
	return strconv.FormatUint(xxhash.Sum64String(s), 10)
}
