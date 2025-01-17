package message

import (
	"encoding/hex"

	"google.golang.org/protobuf/proto"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/client/pb/msg"
)

/*
var imgOld = []byte{
	0x15, 0x36, 0x20, 0x39, 0x32, 0x6B, 0x41, 0x31, 0x00, 0x38, 0x37, 0x32, 0x66, 0x30, 0x36, 0x36, 0x30, 0x33, 0x61, 0x65, 0x31, 0x30, 0x33, 0x62, 0x37, 0x20, 0x20, 0x20, 0x20, 0x20,
	0x20, 0x35, 0x30, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x7B, 0x30, 0x31, 0x45, 0x39, 0x34, 0x35, 0x31, 0x42, 0x2D, 0x37, 0x30, 0x45, 0x44,
	0x2D, 0x45, 0x41, 0x45, 0x33, 0x2D, 0x42, 0x33, 0x37, 0x43, 0x2D, 0x31, 0x30, 0x31, 0x46, 0x31, 0x45, 0x45, 0x42, 0x46, 0x35, 0x42, 0x35, 0x7D, 0x2E, 0x70, 0x6E, 0x67, 0x41,
}
*/

func (e *TextElement) Pack() (r []*msg.Elem) {
	r = append(r, &msg.Elem{
		Text: &msg.Text{
			Str: &e.Content,
		},
	})
	return
}

func (e *FaceElement) Pack() (r []*msg.Elem) {
	r = []*msg.Elem{}
	if e.Index >= 260 {
		elem := &msg.MsgElemInfoServtype33{
			Index:  proto.Uint32(uint32(e.Index)),
			Text:   []byte("/" + e.Name),
			Compat: []byte("/" + e.Name),
		}
		b, _ := proto.Marshal(elem)
		r = append(r, &msg.Elem{
			CommonElem: &msg.CommonElem{
				ServiceType:  proto.Int32(33),
				PbElem:       b,
				BusinessType: proto.Int32(1),
			},
		})
	} else {
		r = append(r, &msg.Elem{
			Face: &msg.Face{
				Index: &e.Index,
				Old:   binary.ToBytes(int16(0x1445 - 4 + e.Index)),
				Buf:   []byte{0x00, 0x01, 0x00, 0x04, 0x52, 0xCC, 0xF5, 0xD0},
			},
		})
	}
	return
}

func (e *AtElement) Pack() (r []*msg.Elem) {
	r = []*msg.Elem{}
	if e.Guild {
		pb, _ := proto.Marshal(&msg.TextResvAttr{AtType: proto.Uint32(2), AtMemberTinyid: proto.Uint64(uint64(e.Target))})
		r = append(r, &msg.Elem{
			Text: &msg.Text{
				Str:       &e.Display,
				PbReserve: pb,
			},
		})
	} else {
		r = append(r, &msg.Elem{
			Text: &msg.Text{
				Str: &e.Display,
				Attr6Buf: binary.NewWriterF(func(w *binary.Writer) {
					w.WriteUInt16(1)
					w.WriteUInt16(0)
					w.WriteUInt16(uint16(len([]rune(e.Display))))
					w.WriteByte(func() byte {
						if e.Target == 0 {
							return 1
						}
						return 0
					}())
					w.WriteUInt32(uint32(e.Target))
					w.WriteUInt16(0)
				}),
			},
		})
	}
	r = append(r, &msg.Elem{Text: &msg.Text{Str: proto.String(" ")}})
	return
}

func (e *ServiceElement) Pack() (r []*msg.Elem) {
	r = []*msg.Elem{}
	// id =35 已移至 ForwardElement
	if e.Id == 1 {
		r = append(r, &msg.Elem{
			Text: &msg.Text{Str: &e.ResId},
		})
	}
	r = append(r, &msg.Elem{
		RichMsg: &msg.RichMsg{
			Template1: append([]byte{1}, binary.ZlibCompress([]byte(e.Content))...),
			ServiceId: &e.Id,
		},
	})
	return
}

func (e *LightAppElement) Pack() (r []*msg.Elem) {
	r = []*msg.Elem{}
	r = append(r, &msg.Elem{
		LightApp: &msg.LightAppElem{
			Data: append([]byte{1}, binary.ZlibCompress([]byte(e.Content))...),
			// MsgResid: []byte{1},
		},
	})
	return
}

func (e *ShortVideoElement) Pack() (r []*msg.Elem) {
	r = append(r, &msg.Elem{
		Text: &msg.Text{
			Str: proto.String("你的QQ暂不支持查看视频短片，请期待后续版本。"),
		},
	})
	r = append(r, &msg.Elem{
		VideoFile: &msg.VideoFile{
			FileUuid:               e.Uuid,
			FileMd5:                e.Md5,
			FileName:               []byte(hex.EncodeToString(e.Md5) + ".mp4"),
			FileFormat:             proto.Int32(3),
			FileTime:               proto.Int32(10),
			FileSize:               proto.Int32(e.Size),
			ThumbWidth:             proto.Int32(1280),
			ThumbHeight:            proto.Int32(720),
			ThumbFileMd5:           e.ThumbMd5,
			ThumbFileSize:          proto.Int32(e.ThumbSize),
			BusiType:               proto.Int32(0),
			FromChatType:           proto.Int32(-1),
			ToChatType:             proto.Int32(-1),
			BoolSupportProgressive: proto.Bool(true),
			FileWidth:              proto.Int32(1280),
			FileHeight:             proto.Int32(720),
		},
	})
	return
}
