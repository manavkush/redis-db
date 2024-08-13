package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Config struct {
	ListenAddr string
}

const defaultListenAddr = ":6379"

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitCh    chan *Peer
	aof       *Aof
}

// Creates a new server from a given config
func NewServer(conf Config) *Server {
	if len(conf.ListenAddr) == 0 {
		conf.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:    conf,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan *Peer),
	}
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case peer := <-s.quitCh:
			if _, ok := s.peers[peer]; ok {
				delete(s.peers, peer)
			}
		}
	}
}

// Starts the server to listen on the listenAddress
func (s *Server) Start() error {
	slog.Info("Starting server.", "listenAddr", s.ListenAddr)
	// log.Printf("Starting to listen on %v\n", s.ListenAddr)
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("Couldn't start the server. Err: ", err)
		return err
	}

	s.ln = ln

	// start the channel loop in the background. This is used for adding/removing peers
	go s.loop()

	// start the accept loop
	return s.acceptLoop()
}

// acceptLoop listens for incoming connections, accepts them and starts a go routine for each of them
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error.", "err", err)
			continue
		}

		// handle the new connection in a separate loop so that we can return to accepting.
		// this allows multiple connections to be served concurrently
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerCh <- peer

	slog.Info("new peer connected.", "remoteAddr", conn.RemoteAddr())

	// block by running a read loop
	peer.readLoop(s.aof)

	s.quitCh <- peer
}

// migrate the old data from the append only log
func (s *Server) migrate() error {
	s.aof.Read(func(value Value) {
		if value.typ != TYPE_ARRAY || len(value.array) == 0 {
			return
		}
		command := COMMAND(strings.ToUpper(value.array[0].bulk))
		commandArgs := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			slog.Error("Invalid command type found.", "command", command)
		}
		handler(commandArgs)
	})
	return nil
}

func main() {
	server := NewServer(Config{})

	// Create an Aof
	aof, err := NewAof("redis.aof")
	if err != nil {
		slog.Error("Migration failed. Unable to initialize AOF", "err", err)
		os.Exit(1)
	}
	defer aof.Close()

	server.aof = aof // add the Aof to the server

	if err := server.migrate(); err != nil { // Migrate from the Aof
		slog.Error("Unable to migrate.", "err", err)
	}

	log.Fatal(server.Start()) // Start the server
}
