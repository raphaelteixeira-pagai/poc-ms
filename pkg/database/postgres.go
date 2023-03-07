package database

import (
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/golangsugar/chatty"
	"os"
	"sync"
)

type Postgres struct {
	mu        sync.Mutex
	closed    bool
	resources chan *dbr.Session
	db        *dbr.Connection
}

// NewPostgres returns a new instance of DBPool
func NewPostgres(size uint) (DBPool, error) {
	if size <= 0 {
		return nil, errors.New("size too small")
	}

	return &Postgres{
		mu:        sync.Mutex{},
		closed:    false,
		resources: make(chan *dbr.Session, size),
		db:        connectPostgres(size),
	}, nil
}

// Acquire returns a new DB session to process queries
func (p *Postgres) Acquire() (*dbr.Session, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			return nil, errors.New("pool closed")
		}
		return r, nil
	default:
		return p.factorySession(), nil
	}
}

// Release try to put connection back to the pool
// if not able to put back, connection is closed
func (p *Postgres) Release(resource *dbr.Session) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		_ = resource.Close()
		return
	}
	select {
	case p.resources <- resource:
	default:
		_ = resource.Close()
	}
}

// Close finish the resources pool
func (p *Postgres) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	p.closed = true

	close(p.resources)

	for r := range p.resources {
		_ = r.Close()
	}
}

// factorySession produces new sessions to be used into pool
func (p *Postgres) factorySession() *dbr.Session {
	return p.db.NewSession(nil)
}

// Ping try to ping database
func (p *Postgres) Ping() bool {
	chatty.Info("trying to ping database...")
	err := p.db.Ping()
	if err != nil {
		chatty.Errorf("failed to ping db, err: %v", err)
		return false
	}

	return true
}

// connectPostgres try to connectPostgres DB with environment variable
func connectPostgres(size uint) *dbr.Connection {
	connStr := os.Getenv("DB_URI")
	if connStr == "" {
		chatty.Fatal("empty db connection string")
	}

	db, err := dbr.Open("postgres", connStr, nil)
	if err != nil {
		chatty.Fatalf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(int(size))
	db.SetMaxIdleConns(int(size))
	return db
}
