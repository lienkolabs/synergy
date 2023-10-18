package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/breeze/util"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

const cookieName = "synergySession"

type filePasswordManager struct {
	mu        sync.Mutex
	file      os.File
	hashes    []crypto.Hash
	passwords map[crypto.Token]int
}

func (f *filePasswordManager) Check(user crypto.Token, password crypto.Hash) bool {
	if n, ok := f.passwords[user]; ok {
		if n >= len(f.hashes) {
			log.Printf("unexpected error in file password manager")
			return false
		}
		return password.Equal(f.hashes[n])
	}
	return false
}

func (f *filePasswordManager) Set(user crypto.Token, password crypto.Hash) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	data := append(user[:], password[:]...)
	n, ok := f.passwords[user]
	if ok {
		if n > len(f.hashes) {
			log.Printf("unexpected error in file password manager")
			return false
		}
		if n, err := f.file.WriteAt(data, int64(n)*2*crypto.Size); n != 64 || err != nil {
			log.Printf("unexpected error in file password manager: %v", err)
			return false
		}
		f.hashes[n] = password
	}
	f.file.Seek(0, 2)
	if n, err := f.file.Write(data); n != 64 || err != nil {
		log.Printf("unexpected error in file password manager: %v", err)
		return false
	}
	f.hashes = append(f.hashes, password)
	f.passwords[user] = len(f.hashes) - 1
	return true
}

func NewFilePasswordManager(filename string) PasswordManager {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("could not open password manager file: %v", err)
	}
	stat, err := file.Stat()
	if err != nil {
		log.Fatalf("could not stat password manager file: %v", err)
	}
	size := stat.Size()
	if size%64 != 0 {
		log.Fatal("corrupted password manager file: incompatible size")
	}
	manager := filePasswordManager{
		file:      *file,
		hashes:    make([]crypto.Hash, size/64),
		passwords: make(map[crypto.Token]int),
	}
	entry := make([]byte, 64)
	for n := 0; n < int(size)/64; n++ {
		if n, err := file.Read(entry); n != 64 {
			log.Fatalf("corrupted password manager file: %v", err)
		}
		var token crypto.Token
		copy(token[:], entry[:32])
		copy(manager.hashes[n][:], entry[32:])
		manager.passwords[token] = n
	}
	return &manager
}

type PasswordManager interface {
	Check(user crypto.Token, password crypto.Hash) bool
	Set(user crypto.Token, password crypto.Hash) bool
}

const cookieLifeItemSeconds = 60 * 60 * 24 * 7 // 1 week

func newCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    url.QueryEscape(value),
		MaxAge:   cookieLifeItemSeconds,
		Secure:   true,
		HttpOnly: true,
	}
}

type AttorneyGeneral struct {
	epoch        uint64
	pk           crypto.PrivateKey
	credentials  PasswordManager
	wallet       crypto.PrivateKey
	pending      map[crypto.Hash]actions.Action
	gateway      social.Gatewayer
	state        *state.State
	templates    *template.Template
	indexer      *index.Index
	session      map[string]crypto.Token
	sessionend   map[uint64][]string
	genesisTime  time.Time
	ephemeralprv crypto.PrivateKey
	ephemeralpub crypto.Token
}

func (a *AttorneyGeneral) CreateSession(handle string, password string) string {
	token, ok := a.state.MembersIndex[handle]
	if !ok {
		return ""
	}
	if password != "1234" {
		return ""
	}
	seed := make([]byte, 32)
	if n, err := rand.Read(seed); n != 32 || err != nil {
		log.Printf("unexpected error in cookie generation:%v", err)
		return ""
	}
	cookie := hex.EncodeToString(seed)
	a.session[cookie] = token
	epoch := a.state.Epoch + cookieLifeItemSeconds
	if ends, ok := a.sessionend[epoch]; ok {
		a.sessionend[epoch] = append(ends, cookie)
	} else {
		a.sessionend[epoch] = []string{cookie}
	}
	return cookie
}

func (a *AttorneyGeneral) Author(r *http.Request) crypto.Token {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Println("cookie", err)
		return crypto.ZeroToken
	}
	if token, ok := a.session[cookie.Value]; ok {
		fmt.Println("token", token)
		return token
	}
	fmt.Println("zero token")
	return crypto.ZeroToken
}

func (a *AttorneyGeneral) Handle(r *http.Request) string {
	author := a.Author(r)
	handle := a.state.Members[crypto.HashToken(author)]
	return handle
}

func (a *AttorneyGeneral) Send(all []actions.Action, author crypto.Token) {
	for _, action := range all {
		dressed := a.DressAction(action, author)
		a.gateway.Action(dressed)
	}
}

// Dress a giving action with current epoch, attorney´s author
// attorneys signature, attorneys wallet and wallet signature
func (a *AttorneyGeneral) DressAction(action actions.Action, author crypto.Token) []byte {
	bytes := action.Serialize()
	dress := []byte{0}
	util.PutUint64(a.epoch, &dress)
	util.PutToken(author, &dress)
	dress = append(dress, 0, 1, 0, 0) // axe void synergy
	dress = append(dress, bytes[8+crypto.TokenSize:]...)
	util.PutToken(a.pk.PublicKey(), &dress)
	signature := a.pk.Sign(dress)
	util.PutSignature(signature, &dress)
	util.PutToken(a.wallet.PublicKey(), &dress)
	util.PutUint64(0, &dress) // zero fee
	signature = a.wallet.Sign(dress)
	util.PutSignature(signature, &dress)
	return dress
}

func (a *AttorneyGeneral) Confirmed(hash crypto.Hash) {
	delete(a.pending, hash)
}
