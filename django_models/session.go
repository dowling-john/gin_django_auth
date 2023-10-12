package django_models

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"math/rand"
	"time"

	"encoding/base64"

	"github.com/uptrace/bun"
	"github.com/enricofoltran/signing"
)

const (
	Seperator string = ":"
	Salt string = "django.core.signing"
)

var validDigits = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var key string = "django-insecure-2z74-f1&^1xqaiw4!9^32^c*(9zr-zs1y5w2j9dlftb4@rz_3d"


type (
	Session struct {
		bun.BaseModel `bun:"table:django_session"`

		SessionKey  string    `bun:"session_key,pk"`
		SessionData string    `bun:"session_data"`
		ExpireDate  time.Time `bun:"expire_date"`
	}
)

func (s *Session) CreateKey() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = validDigits[rand.Intn(len(validDigits))]
	}
	return string(b)
}

func (s *Session) CompressObject(uncompressedObject []byte) ([]byte) {
	compressionBuffer := &bytes.Buffer{}
	compressionWriter := zlib.NewWriter(compressionBuffer)
	compressionWriter.Write(uncompressedObject)
	compressionWriter.Close()
	return compressionBuffer.Bytes()
}

func (s *Session) EncodeObject(objectToEncode []byte) (string) {
	return fmt.Sprintf(".%s", base64.RawURLEncoding.EncodeToString(objectToEncode))
}

func (s *Session) SignObject(data []byte) string {
	encodedObject := s.EncodeObject(s.CompressObject(data))
	timeStampSigner, _ := signing.NewTimestampSigner(key, Seperator, Salt)
	timeStampedSignedString := timeStampSigner.Sign(encodedObject)
	signer, _ := signing.NewSigner(key, Seperator, Salt)
	return signer.Sign(timeStampedSignedString)
}
