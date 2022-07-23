package tools

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/JoshuaDoes/gotena/utils"
)

type UGO struct {
	Layout int        `json:"layout"`
	Title  UgoTitle   `json:"title"`
	Assets []UgoAsset `json:"contents,omitempty"`
}

type UgoTitle struct {
	Names [5]string `json:"names"`
	Index int       `json:"index"`
}

type UgoLayout int
type UgoAssetType int

func FileMagic() [4]byte { return [4]byte{'U', 'G', 'A', 'R'} }

const (
	TypeButton UgoAssetType = iota
	TypeCategory
	TypePost
)

type UgoAsset struct {
	Type   UgoAssetType `json:"type,omitempty"`
	Name   string       `json:"name,omitempty"`
	Url    string       `json:"url,omitempty"`
	Image  string       `json:"image,omitempty"`
	Number int          `json:"number,omitempty"`
}

func ParseUgo(data []byte) (*UGO, error) {
	var ugo UGO
	err := json.Unmarshal(data, &ugo)
	if err != nil {
		return nil, err
	}

	return &ugo, nil
}

func (u *UGO) Pack() ([]byte, error) {
	w := bytes.NewBuffer([]byte{})

	sections := uint32(0)

	tableOfContents := bytes.NewBuffer([]byte{})

	binary.Write(tableOfContents, binary.LittleEndian, []byte(fmt.Sprintf("0\t%d\n", u.Layout)))
	labels := []string{}

	for _, name := range u.Title.Names {
		labels = append(labels, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(name)))
	}

	binary.Write(tableOfContents, binary.LittleEndian, []byte(fmt.Sprintf("1\t0\t%s\t%s\t%s\t%s\t%s\n", labels[0], labels[1], labels[2], labels[3], labels[4])))

	if u.Assets != nil {
		sections++

		for _, asset := range u.Assets {
			switch asset.Type {
			case TypeButton:
				writeButton(tableOfContents, asset)
			case TypeCategory:
				writeCategory(tableOfContents, asset)
			}
		}

		for len(tableOfContents.Bytes())%4 != 0 {
			alignment := make([]byte, 4-len(tableOfContents.Bytes())%4)
			tableOfContents.Write(alignment)
		}
	}

	binary.Write(w, binary.LittleEndian, FileMagic())
	binary.Write(w, binary.LittleEndian, sections)

	binary.Write(w, binary.LittleEndian, uint32(len(tableOfContents.Bytes())))
	binary.Write(w, binary.LittleEndian, tableOfContents.Bytes())

	return w.Bytes(), nil
}

func writeButton(w io.Writer, asset UgoAsset) {
	binary.Write(w, binary.LittleEndian, []byte(fmt.Sprintf("4\t%s\t%d\t%s\n", asset.Url, asset.Number, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(asset.Name)))))
}

func writeCategory(w io.Writer, asset UgoAsset) {
	binary.Write(w, binary.LittleEndian, []byte(fmt.Sprintf("2\t%s\t%s\t%d\n", asset.Url, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(asset.Name)), asset.Number)))
}
