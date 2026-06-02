package mockserver

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Order struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	Items     []OrderItem    `json:"items"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
}

type OrderItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Server struct {
	port   int
	server *http.Server
	mu     sync.RWMutex
	users  []User
	orders []Order
	nextID int
}

func New() *Server {
	s := &Server{
		users:  seedUsers(),
		orders: seedOrders(),
		nextID: 100,
	}
	return s
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("HEAD /api/health", s.handleHealth)
	mux.HandleFunc("OPTIONS /api/users", s.handleOptions)
	mux.HandleFunc("GET /api/users", s.handleListUsers)
	mux.HandleFunc("POST /api/users", s.handleCreateUser)
	mux.HandleFunc("GET /api/users/{id}", s.handleGetUser)
	mux.HandleFunc("PUT /api/users/{id}", s.handleUpdateUser)
	mux.HandleFunc("PATCH /api/users/{id}", s.handlePatchUser)
	mux.HandleFunc("DELETE /api/users/{id}", s.handleDeleteUser)
	mux.HandleFunc("GET /api/orders", s.handleListOrders)
	mux.HandleFunc("POST /api/orders", s.handleCreateOrder)
	mux.HandleFunc("GET /api/orders/{id}", s.handleGetOrder)
	mux.HandleFunc("DELETE /api/orders/{id}", s.handleDeleteOrder)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	s.port = listener.Addr().(*net.TCPAddr).Port

	s.server = &http.Server{
		Handler: corsMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		s.server.Shutdown(context.Background())
	}()

	return s.server.Serve(listener)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "HEAD" {
		w.WriteHeader(200)
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.WriteHeader(204)
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	users := s.users
	s.mu.RUnlock()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	start := (page - 1) * size
	if start > len(users) {
		writeJSON(w, 200, []User{})
		return
	}
	end := start + size
	if end > len(users) {
		end = len(users)
	}
	writeJSON(w, 200, users[start:end])
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid user id")
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			writeJSON(w, 200, u)
			return
		}
	}
	writeError(w, 404, "user not found")
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid json body")
		return
	}
	if input.Name == "" {
		writeError(w, 400, "name is required")
		return
	}

	s.mu.Lock()
	s.nextID++
	user := User{
		ID:        s.nextID,
		Name:      input.Name,
		Email:     input.Email,
		Age:       input.Age,
		Role:      input.Role,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if user.Role == "" {
		user.Role = "user"
	}
	s.users = append(s.users, user)
	s.mu.Unlock()

	writeJSON(w, 201, user)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid user id")
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid json body")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.users {
		if u.ID == id {
			if input.Name != "" {
				s.users[i].Name = input.Name
			}
			if input.Email != "" {
				s.users[i].Email = input.Email
			}
			if input.Age > 0 {
				s.users[i].Age = input.Age
			}
			if input.Role != "" {
				s.users[i].Role = input.Role
			}
			writeJSON(w, 200, s.users[i])
			return
		}
	}
	writeError(w, 404, "user not found")
}

func (s *Server) handlePatchUser(w http.ResponseWriter, r *http.Request) {
	s.handleUpdateUser(w, r)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid user id")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.users {
		if u.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			w.WriteHeader(204)
			return
		}
	}
	writeError(w, 404, "user not found")
}

func (s *Server) handleListOrders(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	orders := s.orders
	s.mu.RUnlock()
	writeJSON(w, 200, orders)
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid order id")
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, o := range s.orders {
		if o.ID == id {
			writeJSON(w, 200, o)
			return
		}
	}
	writeError(w, 404, "order not found")
}

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID int         `json:"user_id"`
		Items  []OrderItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, 400, "invalid json body")
		return
	}

	s.mu.Lock()
	s.nextID++
	total := 0.0
	for _, item := range input.Items {
		total += item.Price * float64(item.Quantity)
	}
	order := Order{
		ID:        s.nextID,
		UserID:    input.UserID,
		Items:     input.Items,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.orders = append(s.orders, order)
	s.mu.Unlock()

	writeJSON(w, 201, order)
}

func (s *Server) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, 400, "invalid order id")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i, o := range s.orders {
		if o.ID == id {
			s.orders = append(s.orders[:i], s.orders[i+1:]...)
			w.WriteHeader(204)
			return
		}
	}
	writeError(w, 404, "order not found")
}

func seedUsers() []User {
	return []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Age: 28, Role: "admin", CreatedAt: "2026-01-15T10:00:00Z"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Age: 32, Role: "user", CreatedAt: "2026-02-20T14:30:00Z"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 24, Role: "user", CreatedAt: "2026-03-10T09:15:00Z"},
		{ID: 4, Name: "Diana", Email: "diana@example.com", Age: 35, Role: "moderator", CreatedAt: "2026-04-05T16:45:00Z"},
		{ID: 5, Name: "Eve", Email: "eve@example.com", Age: 29, Role: "user", CreatedAt: "2026-05-01T08:00:00Z"},
	}
}

func seedOrders() []Order {
	return []Order{
		{ID: 1, UserID: 1, Items: []OrderItem{{Name: "Widget A", Price: 19.99, Quantity: 2}, {Name: "Widget B", Price: 9.99, Quantity: 1}}, Total: 49.97, Status: "completed", CreatedAt: "2026-05-15T11:00:00Z"},
		{ID: 2, UserID: 2, Items: []OrderItem{{Name: "Gadget X", Price: 99.99, Quantity: 1}}, Total: 99.99, Status: "pending", CreatedAt: "2026-05-20T15:30:00Z"},
		{ID: 3, UserID: 3, Items: []OrderItem{{Name: "Service Y", Price: 29.99, Quantity: 3}}, Total: 89.97, Status: "shipped", CreatedAt: "2026-05-25T09:00:00Z"},
	}
}
