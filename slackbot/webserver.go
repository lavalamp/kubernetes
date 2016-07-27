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
	"flag"
	"fmt"
	"log"
	"net/http"

	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

var ()

type Server struct {
	client *client.Client
}

// serveStatus returns "pass", "running", or "fail".
func (s *Server) serveScale(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request: %#v", r)
	fmt.Fprintf(w, "not implemented yet")
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
