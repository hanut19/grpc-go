/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func formatIP(addr string) (addrIP string, ok bool) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	return "[" + addr + "]", true
}

func newformatIP(addr string) (addrIP string, ok bool) {
	// Split the address into IP and zone parts if it contains a zone identifier
	ipStr := addr
	if strings.Contains(addr, "%") {
		ipStr = addr[:strings.Index(addr, "%")]
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", false
	}
	if ip.To4() != nil {
		return addr, true
	}
	// IPv6 addresses need brackets
	if strings.Contains(addr, "%") {
		return "[" + ipStr + addr[strings.Index(addr, "%"):] + "]", true
	}
	return "[" + addr + "]", true
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	fmt.Println("+++++++++++++++++++++++++++++++++++")

	testCases := []string{
		"[fe80::1ff:fe23:4567:890a%25eth2]:8080",
		"[fe80::1ff:fe23:4567:890a%25eth2]",
		"fe80::1ff:fe23:4567:890a%25eth2",
		"fe80::1ff:fe23:4567:890a%eth2",
		"fe80::1ff:fe23:4567:890",
		//	"dns://[fe80::1ff:fe23:4567:890a%eth2]:8080",
		//	"[fe80::1ff:fe23:4567:890a%eth2]:8080",
	}

	fmt.Println("\n--------------Ipv6 test using net.ParseIP()---------------------")
	for _, ipv6 := range testCases {
		ip := net.ParseIP(ipv6)
		if ip == nil {
			fmt.Println("Failed using net.ParseIP(): ", ipv6)
		} else {
			fmt.Println("Successfully using net.ParseIP(): ", ip)
		}
	}
	fmt.Println("\n--------------Ipv6 test using formatIP()---------------------")
	for _, ipv6 := range testCases {
		ip, ok := formatIP(ipv6)
		if ok {
			fmt.Println("Successfully using formatIP(): ", ip)
		} else {
			fmt.Println("Failed using formatIP(): ", ipv6)
		}
	}
	fmt.Println("\n--------------Ipv6 test using newformatIP()---------------------")
	for _, ipv6 := range testCases {
		ip, ok := newformatIP(ipv6)
		if ok {
			fmt.Println("Successfully using newformatIP(): ", ip)
		} else {
			fmt.Println("Failed using newformatIP(): ", ipv6)
		}
	}
}
