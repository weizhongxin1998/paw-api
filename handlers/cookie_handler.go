package handlers

import "paw-api/services"

type CookieHandler struct {
	service *services.CookieService
}

func NewCookieHandler() *CookieHandler {
	return &CookieHandler{service: services.NewCookieService()}
}

func (h *CookieHandler) ListCookies(domain string) []services.CookieInfo {
	return h.service.ListCookies(domain)
}

func (h *CookieHandler) ClearCookies() {
	h.service.ClearCookies()
}
