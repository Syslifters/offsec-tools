package basedata

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Common) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Collector":
			z.Collector, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Collector")
				return
			}
		case "Version":
			z.Version, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "Commit":
			z.Commit, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Commit")
				return
			}
		case "Collected":
			z.Collected, err = dc.ReadTime()
			if err != nil {
				err = msgp.WrapError(err, "Collected")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Common) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Collector"
	err = en.Append(0x84, 0xa9, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72)
	if err != nil {
		return
	}
	err = en.WriteString(z.Collector)
	if err != nil {
		err = msgp.WrapError(err, "Collector")
		return
	}
	// write "Version"
	err = en.Append(0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Version)
	if err != nil {
		err = msgp.WrapError(err, "Version")
		return
	}
	// write "Commit"
	err = en.Append(0xa6, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Commit)
	if err != nil {
		err = msgp.WrapError(err, "Commit")
		return
	}
	// write "Collected"
	err = en.Append(0xa9, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteTime(z.Collected)
	if err != nil {
		err = msgp.WrapError(err, "Collected")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Common) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Collector"
	o = append(o, 0x84, 0xa9, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72)
	o = msgp.AppendString(o, z.Collector)
	// string "Version"
	o = append(o, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Version)
	// string "Commit"
	o = append(o, 0xa6, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74)
	o = msgp.AppendString(o, z.Commit)
	// string "Collected"
	o = append(o, 0xa9, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64)
	o = msgp.AppendTime(o, z.Collected)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Common) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Collector":
			z.Collector, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Collector")
				return
			}
		case "Version":
			z.Version, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "Commit":
			z.Commit, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Commit")
				return
			}
		case "Collected":
			z.Collected, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Collected")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Common) Msgsize() (s int) {
	s = 1 + 10 + msgp.StringPrefixSize + len(z.Collector) + 8 + msgp.StringPrefixSize + len(z.Version) + 7 + msgp.StringPrefixSize + len(z.Commit) + 10 + msgp.TimeSize
	return
}