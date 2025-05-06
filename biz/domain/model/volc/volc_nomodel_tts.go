package volc

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/xh-polaris/psych-senior/biz/domain/model"
	"github.com/xh-polaris/psych-senior/biz/infrastructure/util"
	"io"
	"net/http"
	"sync"
)

var _ model.TtsApp = (*VcNoModelTtsApp)(nil)

type VcNoModelTtsApp struct {
	// ws 连接
	ws     *websocket.Conn
	mu     sync.Mutex
	closed bool

	appKey    string
	accessKey string
	speaker   string
	cluster   string
	url       string
	lang      string
	opt       string
	params    map[string]map[string]interface{}

	// seq 发送的消息序列号
	seq int
	// connId 连接id, 标识一次连接
	connId string
	// logId 服务端返回的logId, 用于定位问题
	logId string
	// sessionId
	sessionId string
	// header 是请求头, 携带鉴权信息
	header http.Header
}

// NewVcNoModelTtsApp 构造一个新的
func NewVcNoModelTtsApp(appKey, accessKey, speaker, cluster, url string) *VcNoModelTtsApp {
	connId := uuid.New().String()
	logId := genLogID()
	sessionId := uuid.New().String()
	app := &VcNoModelTtsApp{
		ws:        nil,
		appKey:    appKey,
		accessKey: accessKey,
		url:       url,
		speaker:   speaker,
		cluster:   cluster,
		opt:       optSubmit,
		connId:    connId,
		logId:     logId,
		sessionId: sessionId,
		seq:       1,
		mu:        sync.Mutex{},
	}
	app.buildHTTPHeader()
	return app
}

// Dial 建立ws连接
func (app *VcNoModelTtsApp) Dial() error {
	conn, r, err := websocket.DefaultDialer.DialContext(context.Background(), app.url, app.header)
	if err != nil {
		if r != nil {
			body, parseErr := io.ReadAll(r.Body)
			if parseErr != nil {
				parseErr = fmt.Errorf("parse response body failed: %w", parseErr)
				body = []byte(parseErr.Error())
			}
			err = fmt.Errorf("[code=%s] [body=%s] %w", r.Status, body, err)
		}
	}
	app.ws = conn
	return err
}

// Start 应用层协议握手
func (app *VcNoModelTtsApp) Start() (err error) {
	reqID := app.sessionId
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//平台上查看具体appid
	params["app"]["appid"] = app.appKey
	params["app"]["token"] = "access_token"
	//平台上查看具体集群名称
	params["app"]["cluster"] = app.cluster
	params["user"] = make(map[string]interface{})
	params["user"]["uid"] = app.connId
	params["audio"] = make(map[string]interface{})
	params["audio"]["language"] = app.lang
	params["audio"]["voice_type"] = app.speaker
	params["audio"]["encoding"] = "pcm"
	params["audio"]["rate"] = 24000
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = reqID
	params["request"]["text"] = ""
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = app.opt
	app.params = params
	return nil
}

func (app *VcNoModelTtsApp) Send(text string) (err error) {
	app.params["request"]["text"] = text
	input, err := json.Marshal(app.params)
	if err != nil {
		return err
	}
	input, err = util.GzipCompress(input)
	if err != nil {
		return err
	}
	payloadSize := len(input)
	payloadArr := make([]byte, 4)
	binary.BigEndian.PutUint32(payloadArr, uint32(payloadSize))
	clientRequest := make([]byte, len(defaultHeader))
	copy(clientRequest, defaultHeader)
	clientRequest = append(clientRequest, payloadArr...)
	clientRequest = append(clientRequest, input...)
	err = app.ws.WriteMessage(websocket.BinaryMessage, clientRequest)
	if err != nil {
		return err
	}
	return nil
}

func (app *VcNoModelTtsApp) Receive() []byte {
	if app.ws == nil || app.closed {
		return nil
	}
	_, msg, err := app.ws.ReadMessage()
	if err != nil {
		glog.Errorf("Receive message error: %v", err)
		return nil
	}
	resp, err := parseResponse(msg)
	if err != nil {
		glog.Errorf("Receive response error: %v", err)
		return nil
	}
	return resp.Audio
}

// Close 关闭连接释放资源
func (app *VcNoModelTtsApp) Close() (err error) {
	app.closed = true
	if app.ws == nil {
		return
	}
	return app.ws.Close()
}

// parseResponse 暂时只保留了有用部分, 之后需要再根据协议进行补充
func parseResponse(res []byte) (resp synResp, err error) {
	//protoVersion := res[0] >> 4
	headSize := res[0] & 0x0f
	messageType := res[1] >> 4
	messageTypeSpecificFlags := res[1] & 0x0f
	//serializationMethod := res[2] >> 4
	messageCompression := res[2] & 0x0f
	//reserve := res[3]
	//headerExtensions := res[4 : headSize*4]
	payload := res[headSize*4:]

	//fmt.Printf("            Protocol version: %x - version %d\n",
	//	protoVersion, protoVersion)
	//fmt.Printf("                 Header size: %x - %d bytes\n",
	//	headSize, headSize*4)
	//fmt.Printf("                Message type: %x - %s\n", messageType,
	//	enumMessageType[messageType])
	//fmt.Printf(" Message type specific flags: %x - %s\n", messageTypeSpecificFlags,
	//	enumMessageTypeSpecificFlags[messageTypeSpecificFlags])
	//fmt.Printf("Message serialization method: %x - %s\n",
	//	serializationMethod, enumMessageSerializationMethods[serializationMethod])
	//fmt.Printf("         Message compression: %x - %s\n",
	//	messageCompression, enumMessageCompression[messageCompression])
	//fmt.Printf("                    Reserved: %d\n", reserve)
	//if headSize != 1 {
	//	fmt.Printf("           Header extensions: %s\n",
	//		headerExtensions)
	//}
	// audio-only server response
	if messageType == 0xb {
		// no sequence number as ACK
		if messageTypeSpecificFlags == 0 {
			//fmt.Println("                Payload size: 0")
		} else {
			sequenceNumber := int32(binary.BigEndian.Uint32(payload[0:4]))
			//payloadSize := int32(binary.BigEndian.Uint32(payload[4:8]))
			payload = payload[8:]
			resp.Audio = append(resp.Audio, payload...)
			//fmt.Printf("             Sequence number: %d\n",
			//	sequenceNumber)
			//fmt.Printf("                Payload size: %d\n", payloadSize)
			if sequenceNumber < 0 {
				resp.IsLast = true
			}
		}
	} else if messageType == 0xf {
		//code := int32(binary.BigEndian.Uint32(payload[0:4]))
		errMsg := payload[8:]
		if messageCompression == 1 {
			errMsg, _ = util.GzipDecompress(errMsg)
		}
		//fmt.Printf("                  Error code: %d\n", code)
		//fmt.Printf("                   Error msg: %s\n", string(errMsg))
		err = errors.New(string(errMsg))
		return
	} else if messageType == 0xc {
		//msgSize = int32(binary.BigEndian.Uint32(payload[0:4]))
		payload = payload[4:]
		if messageCompression == 1 {
			payload, _ = util.GzipDecompress(payload)
		}
		fmt.Printf("            Frontend message: %s\n", string(payload))
	} else {
		//fmt.Printf("          wrong message type:%d\n", messageType)
		//err = errors.New("wrong message type")
		return
	}
	return
}

// buildHTTPHeader 构造鉴权请求头
func (app *VcNoModelTtsApp) buildHTTPHeader() {
	app.header = http.Header{"Authorization": []string{fmt.Sprintf("Bearer;%s", app.accessKey)}}
}

// version: b0001 (4 bits)
// header size: b0001 (4 bits)
// message type: b0001 (Full client request) (4bits)
// message type specific flags: b0000 (none) (4bits)
// message serialization method: b0001 (JSON) (4 bits)
// message compression: b0001 (gzip) (4bits)
// reserved data: 0x00 (1 byte)
var defaultHeader = []byte{0x11, 0x10, 0x11, 0x00}

const (
	optQuery  string = "query"
	optSubmit string = "submit"
)

var (
	enumMessageType = map[byte]string{
		11: "audio-only server response",
		12: "frontend server response",
		15: "error message from server",
	}
	enumMessageTypeSpecificFlags = map[byte]string{
		0: "no sequence number",
		1: "sequence number > 0",
		2: "last message from server (seq < 0)",
		3: "sequence number < 0",
	}
	enumMessageSerializationMethods = map[byte]string{
		0:  "no serialization",
		1:  "JSON",
		15: "custom type",
	}
	enumMessageCompression = map[byte]string{
		0:  "no compression",
		1:  "gzip",
		15: "custom compression method",
	}
)

type synResp struct {
	Audio  []byte
	IsLast bool
}
