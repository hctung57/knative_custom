/*
Copyright 2020 The Knative Authors

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
	"io"
	"log"
	"net/http"

	"knative.dev/serving/test"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("serving container received a request.")
	res, err := http.Get("http://127.0.0.1:8882")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, string(resp))
}

func main() {
	flag.Parse()
	log.Print("serving container started")
	test.ListenAndServeGracefully(":8881", handler)
}
