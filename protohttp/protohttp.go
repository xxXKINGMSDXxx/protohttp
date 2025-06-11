package protohttp

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"httpproto/http_pb"
	"io"

	"github.com/golang/protobuf/proto"
)

const minCompressLength = 100

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ResponseToProtobuf(status int, headers map[string]string, body []byte) (*http_pb.HttpResponse, error) {
	statusEnum, ok := httpStatusMap[status]
	if !ok {
		return nil, fmt.Errorf("unsupported HTTP status code: %d", status)
	}
	if headers == nil {
		return nil, errors.New("headers cannot be nil")
	}

	headerList := make([]*http_pb.Header, 0, len(headers))
	for k, v := range headers {
		headerList = append(headerList, &http_pb.Header{Key: proto.String(k), Value: proto.String(v)})
	}

	resp := &http_pb.HttpResponse{
		Status:  statusEnum,
		Headers: headerList,
	}
	if body != nil {
		resp.Body = &http_pb.Body{Content: body}
	}
	return resp, nil
}

func SerializeResponse(resp *http_pb.HttpResponse) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("response is nil")
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		return nil, err
	}
	if len(data) < minCompressLength {
		return proto.Marshal(&http_pb.Envelope{Encoding: http_pb.Envelope_UNCOMPRESSED.Enum(), Response: resp})
	}
	compressed, err := compress(data)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&http_pb.Envelope{Encoding: http_pb.Envelope_COMPRESSED.Enum(), CompressedData: compressed})
}

func DeserializeResponse(data []byte) (*http_pb.HttpResponse, error) {
	env := &http_pb.Envelope{}
	if err := proto.Unmarshal(data, env); err != nil {
		return nil, err
	}
	switch env.GetEncoding() {
	case http_pb.Envelope_UNCOMPRESSED:
		return env.GetResponse(), nil
	case http_pb.Envelope_COMPRESSED:
		decompressed, err := decompress(env.GetCompressedData())
		if err != nil {
			return nil, err
		}
		resp := &http_pb.HttpResponse{}
		if err := proto.Unmarshal(decompressed, resp); err != nil {
			return nil, err
		}
		return resp, nil
	default:
		return nil, fmt.Errorf("unknown encoding: %v", env.GetEncoding())
	}
}

func RequestToProtobuf(method, path string, headers map[string]string, body []byte) (*http_pb.HttpRequest, error) {
	methodEnum, ok := httpMethodMap[method]
	if !ok {
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}
	if headers == nil {
		return nil, errors.New("headers cannot be nil")
	}
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	headerList := make([]*http_pb.Header, 0, len(headers))
	for k, v := range headers {
		headerList = append(headerList, &http_pb.Header{Key: proto.String(k), Value: proto.String(v)})
	}

	req := &http_pb.HttpRequest{
		Method:  methodEnum,
		Path:    proto.String(path),
		Headers: headerList,
	}
	if body != nil {
		req.Body = &http_pb.Body{Content: body}
	}
	return req, nil
}

func SerializeRequest(req *http_pb.HttpRequest) ([]byte, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	if len(data) < minCompressLength {
		return proto.Marshal(&http_pb.Envelope{Encoding: http_pb.Envelope_UNCOMPRESSED.Enum(), Request: req})
	}
	compressed, err := compress(data)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&http_pb.Envelope{Encoding: http_pb.Envelope_COMPRESSED.Enum(), CompressedData: compressed})
}

func DeserializeRequest(data []byte) (*http_pb.HttpRequest, error) {
	env := &http_pb.Envelope{}
	if err := proto.Unmarshal(data, env); err != nil {
		return nil, err
	}
	switch env.GetEncoding() {
	case http_pb.Envelope_UNCOMPRESSED:
		return env.GetRequest(), nil
	case http_pb.Envelope_COMPRESSED:
		decompressed, err := decompress(env.GetCompressedData())
		if err != nil {
			return nil, err
		}
		req := &http_pb.HttpRequest{}
		if err := proto.Unmarshal(decompressed, req); err != nil {
			return nil, err
		}
		return req, nil
	default:
		return nil, fmt.Errorf("unknown encoding: %v", env.GetEncoding())
	}
}

var httpStatusMap = map[int]*http_pb.HttpStatus{
	100: http_pb.HttpStatus_CONTINUE.Enum(),
	101: http_pb.HttpStatus_SWITCHING_PROTOCOLS.Enum(),
	200: http_pb.HttpStatus_OK.Enum(),
	201: http_pb.HttpStatus_CREATED.Enum(),
	202: http_pb.HttpStatus_ACCEPTED.Enum(),
	204: http_pb.HttpStatus_NO_CONTENT.Enum(),
	301: http_pb.HttpStatus_MOVED_PERMANENTLY.Enum(),
	302: http_pb.HttpStatus_FOUND.Enum(),
	304: http_pb.HttpStatus_NOT_MODIFIED.Enum(),
	400: http_pb.HttpStatus_BAD_REQUEST.Enum(),
	401: http_pb.HttpStatus_UNAUTHORIZED.Enum(),
	402: http_pb.HttpStatus_PAYMENT_REQUIRED.Enum(),
	403: http_pb.HttpStatus_FORBIDDEN.Enum(),
	404: http_pb.HttpStatus_NOT_FOUND.Enum(),
	405: http_pb.HttpStatus_METHOD_NOT_ALLOWED.Enum(),
	406: http_pb.HttpStatus_NOT_ACCEPTABLE.Enum(),
	408: http_pb.HttpStatus_REQUEST_TIMEOUT.Enum(),
	409: http_pb.HttpStatus_CONFLICT.Enum(),
	410: http_pb.HttpStatus_GONE.Enum(),
	413: http_pb.HttpStatus_PAYLOAD_TOO_LARGE.Enum(),
	414: http_pb.HttpStatus_URI_TOO_LONG.Enum(),
	415: http_pb.HttpStatus_UNSUPPORTED_MEDIA_TYPE.Enum(),
	429: http_pb.HttpStatus_TOO_MANY_REQUESTS.Enum(),
	500: http_pb.HttpStatus_INTERNAL_SERVER_ERROR.Enum(),
	501: http_pb.HttpStatus_NOT_IMPLEMENTED.Enum(),
	502: http_pb.HttpStatus_BAD_GATEWAY.Enum(),
	503: http_pb.HttpStatus_SERVICE_UNAVAILABLE.Enum(),
	504: http_pb.HttpStatus_GATEWAY_TIMEOUT.Enum(),
	505: http_pb.HttpStatus_HTTP_VERSION_NOT_SUPPORTED.Enum(),
}

var httpMethodMap = map[string]*http_pb.HttpMethod{
	"GET":     http_pb.HttpMethod_GET.Enum(),
	"POST":    http_pb.HttpMethod_POST.Enum(),
	"PUT":     http_pb.HttpMethod_PUT.Enum(),
	"DELETE":  http_pb.HttpMethod_DELETE.Enum(),
	"PATCH":   http_pb.HttpMethod_PATCH.Enum(),
	"HEAD":    http_pb.HttpMethod_HEAD.Enum(),
	"OPTIONS": http_pb.HttpMethod_OPTIONS.Enum(),
}
