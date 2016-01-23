package daemon

import (
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/netutil"
)

type Query struct {
	Options `json:"options"`
}

type Querys []Query

type QuerysRequest struct {
	Querys
	Name   string
	ResChs []chan Response
	EndCh  chan bool
}

type Request struct {
	Name  string
	ResCh chan Response
	EndCh chan bool
}

type Response struct {
	Err error
}

type Daemon struct {
	Config

	QuerysQ     chan QuerysRequest
	Queue       chan Request
	PoolCounter chan bool
	MaxWorker   chan int
	MngQ        []chan ManageRequest
	WorkCount   chan bool
}

func NewDaemon(conf Config) *Daemon {
	return &Daemon{
		Config:    conf,
		Queue:     make(chan Request, conf.WorkerLimit),
		MngQ:      make([]chan ManageRequest, conf.WorkerLimit),
		QuerysQ:   make(chan QuerysRequest, conf.WorkerLimit),
		WorkCount: make(chan bool, conf.WorkerLimit),
	}
}

func (d *Daemon) Serve(ctx context.Context) error {

	go d.QueryHandleWorker(ctx)
	for i := 0; i < d.WorkerLimit; i++ {
		d.MngQ[i] = make(chan ManageRequest, 1)
		go d.Worker(ctx, d.MngQ[i])
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/query/", d.JsonHandler)
	mux.HandleFunc("/json/", d.JsonHandler)
	s := &http.Server{
		Addr:           d.Listen,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	listener, err := net.Listen("tcp", d.Listen)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("listen: %s", d.Listen)
	//err = s.ListenAndServe()
	err = s.Serve(netutil.LimitListener(listener, d.ListenLimit))
	if err != nil {
		log.Printf("ListenAndServe err:%s", err)
	}
	return err
}

func (d *Daemon) UtilServe() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", d.ManageHandler)
	s := &http.Server{
		Addr:           d.ManageListen,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Managelisten: %s", d.ManageListen)
	err := s.ListenAndServe()
	if err != nil {
		log.Printf("ListenAndServe err:%s", err)
	}
}
