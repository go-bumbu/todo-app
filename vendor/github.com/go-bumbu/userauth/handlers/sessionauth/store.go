package sessionauth

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"net/http"

	"github.com/gorilla/sessions"
)

// NewFsStore is a convenience function to generate a new  File system store
// is uses a secure cookie to keep the session id
func NewFsStore(path string, HashKey, BlockKey []byte) (*sessions.FilesystemStore, error) {
	hashL := len(HashKey)
	if hashL != 32 && hashL != 64 {
		return nil, fmt.Errorf("HashKey lenght should be 32 or 64 bytes")
	}
	blockKeyL := len(BlockKey)
	if blockKeyL != 16 && blockKeyL != 24 && blockKeyL != 32 {
		return nil, fmt.Errorf("blockKey lenght should be 16, 24 or 32 bytes")
	}
	fsStore := sessions.NewFilesystemStore(path, HashKey, BlockKey)
	// fsStore.MaxAge() TODO set max age of store

	return fsStore, nil
}

func NewCookieStore(HashKey, BlockKey []byte) (*sessions.CookieStore, error) {
	hashL := len(HashKey)
	if hashL != 32 && hashL != 64 {
		return nil, fmt.Errorf("HashKey lenght should be 32 or 64 bytes")
	}
	blockKeyL := len(BlockKey)
	if blockKeyL != 16 && blockKeyL != 24 && blockKeyL != 32 {
		return nil, fmt.Errorf("blockKey lenght should be 16, 24 or 32 bytes")
	}

	keyPair := make([][]byte, 2)
	keyPair = append(keyPair, HashKey, BlockKey)

	cs := &sessions.CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPair...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Secure:   true,
		},
	}
	cs.MaxAge(cs.Options.MaxAge)
	return cs, nil
}
