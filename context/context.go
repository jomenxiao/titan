package context

import (
	"context"
	"net"
	"sync"
	"time"

	"gitlab.meitu.com/platform/thanos/db"
)

const (
	DefaultNamespace = "default"
)

// Command releated context
type Command struct {
	Name string
	Args []string
}

// Client is the runtime context of a client
type Client struct {
	DB            *db.DB
	Authenticated bool   // Client has be authenticated
	Namespace     string // Namespace of database
	RemoteAddr    string // Client remote address
	ID            int64  // Client uniq ID
	Name          string // Name is set by client setname
	Created       time.Time
	Updated       time.Time
	LastCmd       string
	SkipN         int // Skip N following commands, (-1 for skipping all commands)
	Close         func() error

	// When client is in multi...exec block, the Txn is assigned and Multi is set to be true
	// Before exec, all command called will be queued in Commands
	Txn      *db.Transaction // Txn is set when client is in transaction which is triggered by watch command
	Multi    bool
	Commands []*Command

	Done chan struct{}
}

func NewClient(id int64, conn net.Conn) *Client {
	now := time.Now()
	cli := &Client{
		ID:            id,
		Created:       now,
		Updated:       now,
		Namespace:     DefaultNamespace,
		RemoteAddr:    conn.RemoteAddr().String(),
		Authenticated: false,
		Multi:         false,
		Done:          make(chan struct{}),
		Close:         conn.Close,
	}
	return cli
}

// Server is the runtime context of the server
type Server struct {
	RequirePass string
	Store       *db.RedisStore
	Monitors    sync.Map
	Clients     sync.Map
	Pause       time.Duration // elapse to pause all clients
}

// Context combines the client and server context
type Context struct {
	context.Context
	Client *Client
	Server *Server
}

// New a context
func New(c *Client, s *Server) *Context {
	return &Context{Context: context.Background(), Client: c, Server: s}
}

// CancelFunc tells an operation to abandon its work
type CancelFunc context.CancelFunc

// WithCancel returns a copy of parent with a new Done channel
func WithCancel(parent *Context) (*Context, CancelFunc) {
	ctx := *parent
	child, cancel := context.WithCancel(parent.Context)
	ctx.Context = child
	return &ctx, CancelFunc(cancel)
}

// WithDeadline returns a copy of the parent context with the deadline adjusted to be no later than d
func WithDeadline(parent *Context, d time.Time) (*Context, CancelFunc) {
	ctx := *parent
	child, cancel := context.WithDeadline(parent.Context, d)
	ctx.Context = child
	return &ctx, CancelFunc(cancel)
}

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
func WithTimeout(parent *Context, timeout time.Duration) (*Context, CancelFunc) {
	ctx := *parent
	child, cancel := context.WithTimeout(parent.Context, timeout)
	ctx.Context = child
	return &ctx, CancelFunc(cancel)
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent *Context, key, val interface{}) *Context {
	ctx := *parent
	ctx.Context = context.WithValue(parent.Context, key, val)
	return &ctx
}
