package main

import (
	"fmt"
	"httpproto/http_pb"
	"httpproto/protohttp"
)

func main() {
	// --- RESPONSE EXAMPLE ---
	status := 200
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Cache-Control": "no-cache",
	}
	body := []byte(`{"message": "Hello, World!"}`)

	resp, err := protohttp.ResponseToProtobuf(status, headers, body)
	if err != nil {
		fmt.Println("❌ Failed to create response proto:", err)
		return
	}

	encodedResp, err := protohttp.SerializeResponse(resp)
	if err != nil {
		fmt.Println("❌ Failed to serialize response:", err)
		return
	}

	decodedResp, err := protohttp.DeserializeResponse(encodedResp)
	if err != nil {
		fmt.Println("❌ Failed to deserialize response:", err)
		return
	}

	fmt.Println("✅ Decoded Response:")
	printResponse(decodedResp)

	// --- REQUEST EXAMPLE ---
	reqHeaders := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "ProtoClient/1.0",
	}

	req, err := protohttp.RequestToProtobuf("GET", "https://example.com/api", reqHeaders, nil)
	if err != nil {
		fmt.Println("❌ Failed to create request proto:", err)
		return
	}

	encodedReq, err := protohttp.SerializeRequest(req)
	if err != nil {
		fmt.Println("❌ Failed to serialize request:", err)
		return
	}

	decodedReq, err := protohttp.DeserializeRequest(encodedReq)
	if err != nil {
		fmt.Println("❌ Failed to deserialize request:", err)
		return
	}

	fmt.Println("✅ Decoded Request:")
	printRequest(decodedReq)
}

func printResponse(resp *http_pb.HttpResponse) {
	fmt.Printf("Status: %v\n", resp.GetStatus())
	for _, h := range resp.GetHeaders() {
		fmt.Printf("Header: %s = %s\n", h.GetKey(), h.GetValue())
	}
	if resp.GetBody() != nil {
		fmt.Printf("Body: %s\n", string(resp.GetBody().GetContent()))
	}
}

func printRequest(req *http_pb.HttpRequest) {
	fmt.Printf("Method: %v\n", req.GetMethod())
	fmt.Printf("Path: %s\n", req.GetPath())
	for _, h := range req.GetHeaders() {
		fmt.Printf("Header: %s = %s\n", h.GetKey(), h.GetValue())
	}
	if req.GetBody() != nil {
		fmt.Printf("Body: %s\n", string(req.GetBody().GetContent()))
	}
}
