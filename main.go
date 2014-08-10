package main

import (
	"./eventlog"
	"./idfactory"
	"./lock"
	"./storage"

	"./service"

	"flag"
	"net/http"
)

var (
	listen = flag.String("listen", "localhost:8080", "The adress to listen on [ip:port]")

	// Backends
	switchEventLog    = flag.String("eventlog", "none", "The EventLog backend to use - none, log or amqp")
	switchIdFactory   = flag.String("idfactory", "uuid", "The IdFactory backend to use - uuid or seq")
	switchStorage     = flag.String("storage", "memory", "The Storage backend to use - memory or etcd")
	switchLockManager = flag.String("lockmanager", "memory", "The LockManager backend to use - memory")

	// EventLog
	/// Amqp
	eventlogAmqpAddress = flag.String("eventlog-amqp-address", "amqp://guest:guest@localhost", "The adress url to connect to.")

	// IdFactory
	/// Sequence
	idfactorySequenceFormat = flag.String("idfactory-seq-format", "game%d", "The format to use when creating the game id from the sequence number")

	// Storage
	/// Etcd Storage
	storageEtcdPeer   = flag.String("storage-etcd-peer", "localhost:4001", "The address to find the etcd peer.")
	storageEtcdPrefix = flag.String("storage-etcd-prefix", "moinz.de/tictactoe", "The prefix to store the data in.")
	storageEtcdTtl    = flag.Uint64("storage-etcd-ttl", 365*24*60*60, "How long to keep data in Etcd. Defaults to 1 year.")

	// Lock Manager
	/// memory
	lockmanagerMemoryVerbose = flag.Bool("lockmanager-memory-verbose", false, "Print every lock/unlock to stderr.")
)

func EventLog() service.EventLog {
	switch *switchEventLog {
	case "none":
		return eventlog.NewMockEventLog()
	case "log":
		return eventlog.NewLog()
	case "amqp":
		return eventlog.NewCoresAmqpEventLog(*eventlogAmqpAddress)
	default:
		panic("Unknown eventlog value: " + *switchEventLog)
	}
}

func IdFactory() service.IdFactory {
	switch *switchIdFactory {
	case "uuid":
		return idfactory.UUIDNextId
	case "seq":
		fallthrough
	case "sequence":
		return idfactory.NewSequenceFactory(*idfactorySequenceFormat)
	default:
		panic("Unknown idfactory value: " + *switchIdFactory)
	}
}

func Storage() service.Storage {
	switch *switchStorage {
	case "memory":
		return storage.New()
	case "etcd":
		return storage.NewEtcdStorage(
			*storageEtcdPeer,
			*storageEtcdPrefix,
			*storageEtcdTtl,
		)
	default:
		panic("Unknown storage value: " + *switchStorage)
	}
}

func LockManager() service.LockManager {
	switch *switchLockManager {
	case "memory":
		return lock.New(*lockmanagerMemoryVerbose)
	default:
		panic("Unknown lockmanager value: " + *switchLockManager)
	}
}

func main() {
	flag.Parse()

	eventLog := EventLog()
	idFactory := IdFactory()
	gameStorage := Storage()
	lockManager := LockManager()

	ticTacToeService := service.TicTacToeService{idFactory, gameStorage, eventLog, lockManager}

	handlers := Handlers{&ticTacToeService}
	http.HandleFunc("/game/new", handlers.NewGameHandler)
	http.HandleFunc("/game/move", handlers.MoveHandler)
	http.HandleFunc("/game/get", handlers.GetGameHandler)

	if err := http.ListenAndServe(*listen, nil); err != nil {
		panic(err)
	}
}
