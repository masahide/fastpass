package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/pprof"
	"strings"

	"github.com/masahide/fastpass/status"
)

type ManageCmd uint

const (
	PostNode ManageCmd = 1 << iota
	MoveNode
	DelNode
	GetNodeList
	GetNodes

	PPROF_FILE = "/tmp/gordb.pprof"
)

type ManageRequest struct {
	Cmd   ManageCmd
	Name  string
	Path  string
	ResCh chan ManageResponse
	From  string
	To    string
}

type ManageResponse struct {
	Body interface{}
	Err  error
}

func (d *Daemon) ManageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST": // load
		name := strings.TrimRight(path.Base(r.URL.Path), "/")
		dirpath := r.PostForm.Get("path")
		if name != "" {
			name = strings.TrimRight(path.Base(dirpath), "/")
		}
		err := d.BroadcastManageReq(ManageRequest{Cmd: PostNode, Path: dirpath, Name: name})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
		}
	case "PUT": // move
		name := strings.TrimRight(path.Base(r.URL.Path), "/")
		switch name {
		case "move":
			from := r.PostForm.Get("from")
			to := r.PostForm.Get("to")
			err := d.BroadcastManageReq(ManageRequest{Cmd: MoveNode, From: from, To: to})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
			}
		case "start_pprof":
			f, err := os.Create(PPROF_FILE)
			if err != nil {
				log.Print(err)
				break
			}
			pprof.StartCPUProfile(f)
			log.Printf("StartCPUProfile(%s)", PPROF_FILE)
		case "end_pprof":
			pprof.StopCPUProfile()
			log.Print("StopCPUProfile")
		}
	case "DELETE": //delete
		name := strings.TrimRight(path.Base(r.URL.Path), "/")
		err := d.BroadcastManageReq(ManageRequest{Cmd: DelNode, Name: name})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
		}
	case "GET":
		switch r.URL.Path {
		case "/list":
			res := d.sendManageReq(ManageRequest{Cmd: GetNodeList})
			if res.Err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, res.Err)
			}
			json.NewEncoder(w).Encode(res.Body)
		case "/status":
			stat := status.Get()
			json.NewEncoder(w).Encode(stat)
		}
	}
}

func (d *Daemon) BroadcastManageReq(req ManageRequest) error {
	for _, q := range d.MngQ {
		mr := req
		mr.ResCh = make(chan ManageResponse, 1)
		q <- mr
		res := <-mr.ResCh
		if res.Err != nil {
			return res.Err
		}
	}
	return nil
}

func (d *Daemon) sendManageReq(req ManageRequest) ManageResponse {
	q := d.MngQ[0]
	mr := req
	mr.ResCh = make(chan ManageResponse, 1)
	q <- mr
	res := <-mr.ResCh
	return res
}

func (d *Daemon) manageWork(req ManageRequest) ManageResponse {
	res := ManageResponse{}
	switch req.Cmd {
	case GetNodeList:
		return res
	}
	return ManageResponse{Err: fmt.Errorf("Unkown ManageRequest:%v", req.Cmd)}
}
