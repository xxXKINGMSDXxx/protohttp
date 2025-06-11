# ProtoHTTP

**ProtoHTTP** is a compact, language-agnostic binary protocol for serializing HTTP-style requests and responses using **Protocol Buffers** and optional **zlib compression**.

This repository provides a **Golang implementation** of the ProtoHTTP format.

---

## 🔧 Key Features

- ⚡ Efficient binary format for HTTP data (headers, body, status/methods)
- 🗜️ Automatic zlib compression for large payloads
- 📉 Reduces bandwidth compared to traditional plaintext HTTP
- 📦 Simple envelope structure (`Request` / `Response` + compression flag)
- 🚀 Fast parsing via generated Protobuf code

---

## 📁 Project Structure

```
.
├── go.mod / go.sum               # Go module files
├── http.proto                    # ProtoHTTP schema (language-agnostic)
├── example.go                    # Usage demo in Golang
├── http_pb/                      # Generated Protobuf Go code
│   └── http.pb.go
└── protohttp/                    # Golang implementation
    └── protohttp.go
```

---

## 🚀 Full Example (Golang)

### ➤ Serialize a Response

```go
resp, _ := protohttp.ResponseToProtobuf(200, map[string]string{
    "Content-Type": "application/json",
    "Cache-Control": "no-cache",
}, []byte(`{"message":"Hello, world!"}`))

data, _ := protohttp.SerializeResponse(resp)
decoded, _ := protohttp.DeserializeResponse(data)
```

#### ✅ Decoded Response Output

```
✅ Decoded Response:
Status: OK
Header: Content-Type = application/json
Header: Cache-Control = no-cache
Body: {"message":"Hello, world!"}
```

---

### ➤ Serialize a Request

```go
req, _ := protohttp.RequestToProtobuf("GET", "https://example.com/api", map[string]string{
    "User-Agent": "ProtoClient/1.0",
    "Accept": "application/json",
}, nil)

data, _ := protohttp.SerializeRequest(req)
decoded, _ := protohttp.DeserializeRequest(data)
```

#### ✅ Decoded Request Output

```
✅ Decoded Request:
Method: GET
Path: https://example.com/api
Header: Accept = application/json
Header: User-Agent = ProtoClient/1.0
```

---

## 🔍 ProtoHTTP vs Plain HTTP

| Feature                     | ProtoHTTP                   | Plain HTTP                  |
|----------------------------|-----------------------------|-----------------------------|
| Format                     | Protobuf + zlib (binary)    | Raw plaintext               |
| Size (response example)    | ~90 bytes compressed        | ~270 bytes (text)           |
| Compression                | ✅ Automatic (zlib)          | ❌ Manual if needed         |
| Header encoding            | Protobuf key-value pairs    | Flat strings + delimiters   |
| Parsing cost               | Low (binary decode)         | High (token/string parsing) |
| Transport-safe             | ✅ No encoding escaping Needed | ⚠️ Needs encoding escaping  |

---

## 📦 Compression Behavior

- Messages **< 100 bytes** → Uncompressed
- Messages **≥ 100 bytes** → zlib-compressed automatically

No manual compression logic needed — handled inside serialization.

---

## 🧪 Try the Example

```bash
go mod tidy
go run example.go
```

---

## 🌍 Language Support

ProtoHTTP is designed to be **language-agnostic** since its based on Google's Protocol Buffers. This repository provides:

- ✅ Golang implementation (currently)

The format is defined in [`http.proto`](./http.proto), which can be used to generate bindings in other languages.

---

## 📜 License

MIT License
