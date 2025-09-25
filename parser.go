package cvimodelgo

import (
	"cvimodelgo/model"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"time"
)

var (
	ErrUnsupportedModel = errors.New("unsupported model file")
)

type ModelHeader struct {
	Magic    [8]byte
	BodySize uint32
	Major    byte
	Minor    byte
	Md5      [16]byte
	Chip     [16]byte
	Padding  [2]byte
}

type ModelInfo struct {
	Name        string
	Target      string
	BuildTime   time.Time
	InputQuant  string
	OutputQuant string
	Quant       string
}

func ParseModelFile(r io.Reader) (*ModelInfo, error) {
	// Read first 48 bytes header
	header := ModelHeader{}
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	if strings.ToLower(string(header.Magic[:])) != "cvimodel" {
		return nil, ErrUnsupportedModel
	}

	// Read rest as fb
	d, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	mdl := model.GetRootAsModel(d, 0)

	info := ModelInfo{
		Name:   string(mdl.Name()),
		Target: string(mdl.Target()),
	}
	info.BuildTime, _ = time.Parse("2006-01-02 15:04:05", string(mdl.BuildTime()))

	var p model.Program
	mdl.Programs(&p, 0)

	var t model.Tensor
	if p.TensorMap(&t, 0) {
		info.InputQuant = t.Dtype().String()
	}
	if p.TensorMap(&t, 1) {
		info.Quant = t.Dtype().String()
	}
	if p.TensorMap(&t, p.TensorMapLength()-1) {
		info.OutputQuant = t.Dtype().String()
	}

	return &info, nil
}
