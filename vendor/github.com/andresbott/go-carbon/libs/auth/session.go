package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

type SessionMgr struct {
	store sessions.Store

	sessionDur    time.Duration
	minWriteSpace time.Duration
	maxSessionDur time.Duration
}

type SessionCfg struct {
	Store sessions.Store

	SessionDur    time.Duration // normal session duration, can be renewed on subsequent requests
	MinWriteSpace time.Duration // time between the last session update, used to not overload the session store
	MaxSessionDur time.Duration // force a logout after this time

}

func init() {
	// this is needed so that cookiestore can manage the session data
	gob.Register(SessionData{})
}

func NewSessionMgr(cfg SessionCfg) (*SessionMgr, error) {

	if cfg.SessionDur == 0 {
		cfg.SessionDur = time.Hour * 1
	}
	if cfg.MaxSessionDur == 0 {
		cfg.MaxSessionDur = time.Hour * 24
	}
	if cfg.MinWriteSpace == 0 {
		cfg.MinWriteSpace = time.Minute * 2
	}

	c := SessionMgr{
		sessionDur:    cfg.SessionDur,
		minWriteSpace: cfg.MinWriteSpace,
		maxSessionDur: cfg.MaxSessionDur,
		store:         cfg.Store,
	}
	return &c, nil
}

// CookieStore is a convenience function to generate a new secure cookiestore
// based on the securecookie.New doc:
//
// hashKey is required, used to authenticate values using HMAC. Create it using
// GenerateRandomKey(). It is recommended to use a key with 32 or 64 bytes.
//
// blockKey is optional, used to encrypt values. Create it using
// GenerateRandomKey(). The key length must correspond to the key size
// of the encryption algorithm. For AES, used by default, valid lengths are
// 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
// The default encoder used for cookie serialization is encoding/gob.
func CookieStore(HashKey, BlockKey []byte) (*sessions.CookieStore, error) {
	hashL := len(HashKey)
	if hashL != 32 && hashL != 64 {
		return nil, fmt.Errorf("HashKey lenght should be 32 or 64 bytes")
	}
	blockKeyL := len(BlockKey)
	if blockKeyL != 16 && blockKeyL != 24 && blockKeyL != 32 {
		return nil, fmt.Errorf("blockKey lenght should be 16, 24 or 32 bytes")
	}
	return sessions.NewCookieStore(HashKey, BlockKey), nil
}

// FsStore is a convenience function to generate a new  File system store
// is uses a secure cookie to keep the session id (?)
func FsStore(path string, HashKey, BlockKey []byte) (*sessions.FilesystemStore, error) {
	hashL := len(HashKey)
	if hashL != 32 && hashL != 64 {
		return nil, fmt.Errorf("HashKey lenght should be 32 or 64 bytes")
	}
	blockKeyL := len(BlockKey)
	if blockKeyL != 16 && blockKeyL != 24 && blockKeyL != 32 {
		return nil, fmt.Errorf("blockKey lenght should be 16, 24 or 32 bytes")
	}
	return sessions.NewFilesystemStore(path, HashKey, BlockKey), nil
}

type UserData struct {
	UserId          string // ID or username
	DeviceID        string // hold information about the device
	IsAuthenticated bool
}
type SessionData struct {
	UserData

	// expiration of the session, e.g. 2 days, after a login is required, this value can be updated by "keep me logged in"
	Expiration time.Time
	// force re-auth, max time a session is valid, even if keep logged in is in place.
	ForceReAuth time.Time
	LastUpdate  time.Time
}

func (d *SessionData) Process(extend time.Duration) {
	// check expiration
	if d.Expiration.Before(time.Now()) {
		d.IsAuthenticated = false
	}
	// check hard expiration
	if d.ForceReAuth.Before(time.Now()) {
		d.IsAuthenticated = false
	}
	// extend normal expiration
	if d.IsAuthenticated && extend > 0 {
		d.Expiration = d.Expiration.Add(extend)
	}
}

const (
	SessionName    = "_c_auth"
	sessionDataKey = "data"
)

// Login is a convenience function to write a new logged-in session for a specific user id and write it
func (auth *SessionMgr) Login(r *http.Request, w http.ResponseWriter, user string) error {
	authData := SessionData{
		UserData: UserData{
			UserId: user,

			IsAuthenticated: true,
		},
		Expiration:  time.Now().Add(auth.sessionDur),
		ForceReAuth: time.Now().Add(auth.maxSessionDur),
	}
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		return err
	}
	return auth.write(r, w, session, authData)
}

// Logout is a convenience function to logout the current user
func (auth *SessionMgr) Logout(r *http.Request, w http.ResponseWriter) error {
	authData := SessionData{
		UserData: UserData{
			IsAuthenticated: false,
		},
	}
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		return err
	}
	return auth.write(r, w, session, authData)
}

// the Write public function is commented out for now, until it might be needed to not blow the API
// Use Login() instead
//func (auth *SessionMgr) Write(r *http.Request, w http.ResponseWriter, data SessionData) ( error) {
//	session, err := auth.store.Get(r, SessionName)
//	if err != nil {
//		return err
//	}
//	return auth.write(r, w, session, data)
//}

func (auth *SessionMgr) write(r *http.Request, w http.ResponseWriter, session *sessions.Session, data SessionData) error {
	now := time.Now()
	if data.LastUpdate.Add(auth.minWriteSpace).After(now) {
		return nil
	}
	data.LastUpdate = now

	session.Values[sessionDataKey] = data
	err := session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func (auth *SessionMgr) Read(r *http.Request) (SessionData, error) {
	data, _, err := auth.read(r)
	return data, err
}
func (auth *SessionMgr) read(r *http.Request) (SessionData, *sessions.Session, error) {
	session, err := auth.store.Get(r, SessionName)
	if err != nil {
		// TODO find better solution to handle sessions when the FS store is gone but the client still has a session
		return SessionData{}, nil, err
	}

	key := session.Values[sessionDataKey]
	if key == nil {
		return SessionData{}, nil, err
	}
	authData := key.(SessionData)
	authData.Process(auth.sessionDur)
	return authData, session, err
}

// ReadUpdate is used to read the session, and update the session expiry timestamp
// it only extends the session if enough time has passed since the last write to not overload
// the session store on many requests.
// it returns the session data if the user is logged in
func (auth *SessionMgr) ReadUpdate(r *http.Request, w http.ResponseWriter) (SessionData, error) {
	data, session, err := auth.read(r)
	if err != nil {
		return SessionData{}, err
	}

	if data.IsAuthenticated {
		err = auth.write(r, w, session, data)
		if err != nil {
			return SessionData{}, err
		}
		return data, nil
	}
	return SessionData{}, nil
}

const UserIdKey = "loggedInUserID"
const UserIsLoggedInKey = "isUserLoggedIn"

// Middleware is a simple session auth middleware that will only allow access if the user is logged in
// this can be used as simple implementations or as inspiration to customize an authentication middleware
func (auth *SessionMgr) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := auth.ReadUpdate(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if data.IsAuthenticated {
			// add user ID into request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIdKey, data.UserId)
			ctx = context.WithValue(ctx, UserIsLoggedInKey, true)

			req := r.WithContext(ctx)
			*r = *req

			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

type loginData struct {
	User     string `json:"user"`
	Pw       string `json:"password"`
	Redirect string
}

// FormAuthHandler is a simple session auth handler that will respond to a form POST request and login a user
// this can be used as simple implementations or as inspiration to customize an authentication middleware
func FormAuthHandler(session *SessionMgr, user UserLogin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		var payload loginData
		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&payload, r.PostForm)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if user.AllowLogin(payload.User, payload.Pw) {
			err = session.Login(r, w, payload.User)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

// JsonAuthHandler is a simple session auth handler that will respond to a Json POST request and login a user
// this can be used as simple implementations or as inspiration to customize an authentication middleware
func JsonAuthHandler(session *SessionMgr, user UserLogin) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload loginData

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.AllowLogin(payload.User, payload.Pw) {
			err = session.Login(r, w, payload.User)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
