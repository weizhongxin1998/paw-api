package services

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

var (
	sharedJar   http.CookieJar
	sharedJarMu sync.Mutex
)

func getSharedJar() http.CookieJar {
	sharedJarMu.Lock()
	defer sharedJarMu.Unlock()
	if sharedJar == nil {
		sharedJar, _ = cookiejar.New(nil)
	}
	return sharedJar
}

type CookieInfo struct {
	Domain   string    `json:"domain"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HttpOnly bool      `json:"http_only"`
}

type CookieService struct {
	jar http.CookieJar
}

func NewCookieService() *CookieService {
	return &CookieService{jar: getSharedJar()}
}

func (s *CookieService) ListCookies(domain string) []CookieInfo {
	if s.jar == nil {
		return nil
	}
	if domain == "" {
		return nil
	}
	u, err := url.Parse(domain)
	if err != nil {
		return nil
	}
	cookies := s.jar.Cookies(u)
	var result []CookieInfo
	for _, c := range cookies {
		expires := time.Time{}
		if !c.Expires.IsZero() {
			expires = c.Expires
		}
		result = append(result, CookieInfo{
			Domain:   c.Domain,
			Name:     c.Name,
			Value:    c.Value,
			Path:     c.Path,
			Expires:  expires,
			Secure:   c.Secure,
			HttpOnly: c.HttpOnly,
		})
	}
	return result
}

func (s *CookieService) GetAllCookies() []CookieInfo {
	return nil
}

func (s *CookieService) ClearCookies() {
	jar, _ := cookiejar.New(nil)
	s.jar = jar
}
