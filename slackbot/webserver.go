/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

var ()

type Server struct {
	client *client.Client
}

// serveStatus returns "pass", "running", or "fail".
func (s *Server) serveScale(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got request: %#v\n", r)
	// command := r.URL.Query().Get("command")
	r.ParseForm()
	switch r.FormValue("command") {
	case "/scale":
		scaleString := r.FormValue("text")
		scale, err := strconv.Atoi(scaleString)
		if err != nil {
			fmt.Fprintf(w, "couldn't parse %q: %v", scaleString, err)
			return
		}
		d, err := s.client.Extensions().Deployments("default").Get("slackscale")
		if err != nil {
			fmt.Fprintf(w, "couldn't get deployment: %v", err)
			return
		}
		d.Spec.Replicas = int32(scale)
		_, err = s.client.Extensions().Deployments("default").Update(d)
		if err != nil {
			fmt.Fprintf(w, "couldn't write deployment: %v", err)
			return
		}
		fmt.Fprintf(w, "scaled to %v", scale)
	case "/getscale":
		d, err := s.client.Extensions().Deployments("default").Get("slackscale")
		if err != nil {
			fmt.Fprintf(w, "couldn't get deployment: %v", err)
			return
		}
		fmt.Fprintf(w, "Currently there are %v replicas available. (%v unavailable)", d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
	default:
		fmt.Fprintf(w, "don't know command %v", r.FormValue("command"))
	}
}

func main() {
	flag.Parse()

	config, err := restclient.InClusterConfig()
	if err != nil {
		log.Fatalf("Unable to create config; error: %v\n", err)
	}
	config.ContentType = "application/vnd.kubernetes.protobuf"
	client, err := client.New(config)
	if err != nil {
		log.Fatalf("Unable to create client; error: %v\n", err)
	}

	s := &Server{client: client}

	http.HandleFunc("/scale", s.serveScale)

	go log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))

	select {}
}
