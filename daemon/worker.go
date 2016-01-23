package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"
)

type Options struct {
	OrderDesc  bool   `json:"order_by,omitempty"`
	OrderBy    string `json:"order_desc,omitempty"`
	QueryCache bool   `json:"query_cache,omitempty"`
}

type Worker struct {
	*Daemon
}

func NewWorker(d *Daemon) *Worker {
	return &Worker{
		Daemon: d,
	}

}

func (d *Daemon) JsonHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	name := r.PostForm.Get("name")
	if name == "" {
		name = strings.TrimRight(path.Base(r.URL.Path), "/")
	}
	defer r.Body.Close()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var querys Querys
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&querys)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("json.Decode err:%s", err))
		return
	}
	elapsendJsonDecode := time.Now().Sub(startTime)
	elapsendQuery := time.Now().Sub(startTime) - elapsendJsonDecode
	json.NewEncoder(w).Encode(nil)
	elapsedAll := time.Now().Sub(startTime)
	if d.LogLevel > 0 {
		log.Printf("elapsed:%s, json decode:%s, query:%s, json encode:%s", elapsedAll, elapsendJsonDecode, elapsendQuery, elapsedAll-elapsendQuery-elapsendJsonDecode)
	}
	return

}

func (d *Daemon) QueryHandleWorker(ctx context.Context) {
	for {
		select {
		case qs := <-d.QuerysQ:
			for i := 0; i < len(qs.Querys); i++ {
				select {
				case d.WorkCount <- true:
				case <-ctx.Done():
					return
				}
			}
			for i := 0; i < len(qs.Querys); i++ {
				select {
				case d.Queue <- Request{Name: qs.Name, ResCh: qs.ResChs[i], EndCh: qs.EndCh}:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}

}

func (d *Daemon) QueryStreams(name string, querys Querys) (endCh chan bool, err error) {
	resChs := make([]chan Response, len(querys))
	endCh = make(chan bool)
	for i := 0; i < len(querys); i++ {
		resChs[i] = make(chan Response, 1)
	}
	d.QuerysQ <- QuerysRequest{Querys: querys, Name: name, ResChs: resChs, EndCh: endCh}
	/*
		for i, query := range querys {
			resChs[i] = make(chan Response, 1)
			d.Queue <- Request{Query: query.Stream, Name: name, ResCh: resChs[i], EndCh: endCh}
		}
	*/
	for i := 0; i < len(querys); i++ {
		res := <-resChs[i]
		if res.Err != nil {
			return nil, res.Err
		}
	}
	return endCh, nil

}

func (d *Daemon) Worker(ctx context.Context, ManageCh chan ManageRequest) {
	worker := NewWorker(d)
	for {
		select {
		case req := <-d.Queue:
			res := worker.work(req)
			select {
			case req.ResCh <- res:
			case <-ctx.Done():
				return
			}
			if res.Err != nil {
				log.Printf("work err: %s", res.Err)
			}
			select {
			case <-req.EndCh:
			case <-ctx.Done():
				return
			}
			<-d.WorkCount
		case req := <-ManageCh:
			res := worker.manageWork(req)
			req.ResCh <- res
			if res.Err != nil {
				log.Printf("manageWork err: %s", res.Err)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (d *Worker) work(req Request) Response {
	res := Response{}
	return res
}
