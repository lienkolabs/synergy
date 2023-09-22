package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/lienkolabs/breeze/crypto"
	"github.com/lienkolabs/synergy/social"
	"github.com/lienkolabs/synergy/social/actions"
	"github.com/lienkolabs/synergy/social/index"
	"github.com/lienkolabs/synergy/social/state"
)

const cookieName = "synergyUser"

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
		Name:     "synergySession",
		Value:    url.QueryEscape(value),
		MaxAge:   cookieLifeItemSeconds,
		Secure:   true,
		HttpOnly: true,
	}
}

type AttorneyGeneral struct {
	pk          crypto.PrivateKey
	credentials PasswordManager
	wallet      crypto.PrivateKey
	pending     map[crypto.Hash]actions.Action
	gateway     social.Gatewayer
	state       *state.State
	templates   *template.Template
	indexer     *index.Index
	mux         *http.ServeMux
	session     map[string]crypto.Token
	sessionend  map[uint64][]string
}

func (a *AttorneyGeneral) CreateSession(handle string, password string) string {
	token, ok := a.state.MembersIndex[handle]
	if !ok {
		return ""
	}
	seed := make([]byte, 32)
	if n, err := rand.Read(seed); n != 32 || err != nil {
		log.Printf("unexpected error in cookie generation:%v", err)
		return ""
	}
	a.session[hex.EncodeToString(seed)] = token
	cookie := hex.EncodeToString(seed)
	epoch := a.state.Epoch + cookieLifeItemSeconds
	if ends, ok := a.sessionend[epoch]; ok {
		a.sessionend[epoch] = append(ends, cookie)
	} else {
		a.sessionend[epoch] = []string{cookie}
	}
	return hex.EncodeToString(seed)
}

func (a *AttorneyGeneral) ListenAndServe(port int) chan error {
	check := make(chan error, 2)
	a.templates = template.New("root")
	files := make([]string, len(templateFiles))
	for n, file := range templateFiles {
		files[n] = fmt.Sprintf("./api/templates/%v.html", file)
	}
	t, err := template.ParseFiles(files...)
	if err != nil {
		check <- err
		return check
	}
	a.templates = t
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./api/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs)) //
	//mux.HandleFunc("/api", a.ApiHandler)
	go func() {
		srv := &http.Server{
			Addr:         fmt.Sprintf(":%v", port),
			Handler:      mux,
			WriteTimeout: 2 * time.Second,
		}
		srv.ListenAndServe()
		check <- err
	}()
	return check
}

func NewAttorneyGeneral(pk crypto.PrivateKey, port int, gateway social.Gatewayer, indexer *index.Index, credentials PasswordManager) *AttorneyGeneral {
	attorney := AttorneyGeneral{
		pk:          pk,
		wallet:      pk,
		pending:     make(map[crypto.Hash]actions.Action),
		gateway:     gateway,
		state:       gateway.State(),
		indexer:     indexer,
		credentials: credentials,
	}
	//attorney.HandleWithHash("editview.html", "/editview/", EditDetailFromState2)
	return &attorney
}

func (a *AttorneyGeneral) Author(r *http.Request) crypto.Token {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return crypto.ZeroToken
	}
	if token, ok := a.session[cookie.Value]; ok {
		return token
	}
	return crypto.ZeroToken
}

type ViewerWithHash func(*state.State, *index.Index, crypto.Hash, crypto.Token) any

type ViewerWithString func(*state.State, *index.Index, crypto.Token, string) any

type Viewer func(*state.State, *index.Index, crypto.Token) any

func (a *AttorneyGeneral) HandleWithHash(template, path string, viewer ViewerWithHash) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		hash := getHash(r.URL.Path, path)
		author := a.Author(r)
		view := viewer(a.state, a.indexer, hash, author)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}

func (a *AttorneyGeneral) HandleWithString(template, path string, viewer ViewerWithString) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		item := strings.Replace(r.URL.Path, path, "", 1)
		author := a.Author(r)
		view := viewer(a.state, a.indexer, author, item)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}

func (a *AttorneyGeneral) HandleSimple(template, path string, viewer Viewer) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		author := a.Author(r)
		view := viewer(a.state, a.indexer, author)
		if err := a.templates.ExecuteTemplate(w, template, view); err != nil {
			log.Println(err)
		}
	}
	a.mux.HandleFunc(path, handler)
}
