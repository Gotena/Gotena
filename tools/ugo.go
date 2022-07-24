package tools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/JoshuaDoes/gotena/utils"
	crunch "github.com/superwhiskers/crunch/v3"
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

func FileMagic() []byte { return []byte{'U', 'G', 'A', 'R'} }

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
	Index1 int          `json:"number,omitempty"`
	Index2 int          `json:"number2,omitempty"`
	Index3 int          `json:"number3,omitempty"`
	Index4 int          `json:"number4,omitempty"`
}

func PackUgoJson(path string) ([]byte, error) {
	ugoJSON, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("[UGO] Error reading UGO JSON:", err)
		return nil, err
	}

	ugo, err := ParseUgo(ugoJSON)
	if err != nil {
		fmt.Println("[UGO] Error parsing UGO JSON:", err)
		return nil, err
	}

	ugoBytes, err := ugo.Pack()
	if err != nil {
		fmt.Println("[UGO] Error packing UGO JSON:", err)
		return nil, err
	}

	return ugoBytes, nil
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
	next := make([]byte, 0)
	w := crunch.NewBuffer(make([]byte, 0))

	tableOfContents := crunch.NewBuffer(make([]byte, 0))
	extraData := crunch.NewBuffer(make([]byte, 0))

	next = []byte(fmt.Sprintf("0\t%d\n", u.Layout))
	tableOfContents.Grow(int64(len(next)))
	tableOfContents.WriteBytesNext(next)
	labels := []string{}

	for _, name := range u.Title.Names {
		labels = append(labels, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(name)))
	}

	next = []byte(fmt.Sprintf("1\t0\t%s\t%s\t%s\t%s\t%s\n", labels[0], labels[1], labels[2], labels[3], labels[4]))
	tableOfContents.Grow(int64(len(next)))
	tableOfContents.WriteBytesNext(next)

	sections := uint32(0)
	if u.Assets != nil {
		sections++

		for _, asset := range u.Assets {
			switch asset.Type {
			case TypeButton:
				writeButton(tableOfContents, extraData, asset)
			case TypeCategory:
				writeCategory(tableOfContents, asset)
			}
		}

		for len(tableOfContents.Bytes())%4 != 0 {
			alignment := make([]byte, 4-len(tableOfContents.Bytes())%4)
			tableOfContents.Grow(int64(len(alignment)))
			tableOfContents.WriteBytesNext(alignment)
		}
	}

	if len(extraData.Bytes()) > 0 {
		sections++
		for len(extraData.Bytes())%4 != 0 {
			alignment := make([]byte, 4-len(extraData.Bytes())%4)
			extraData.Grow(int64(len(alignment)))
			extraData.WriteBytesNext(alignment)
		}
	}

	w.Grow(4 + 4)
	next = FileMagic()
	w.WriteBytesNext(next)
	w.WriteU32LENext([]uint32{sections})

	if sections > 0 {
		w.Grow(4)
		w.WriteU32LENext([]uint32{uint32(tableOfContents.ByteCapacity())})
	}
	if sections > 1 {
		w.Grow(4)
		w.WriteU32LENext([]uint32{uint32(extraData.ByteCapacity())})
	}
	w.Grow(int64(tableOfContents.ByteCapacity()))
	w.WriteBytesNext(tableOfContents.Bytes())
	w.Grow(int64(extraData.ByteCapacity()))
	w.WriteBytesNext(extraData.Bytes())

	return w.Bytes(), nil
}

func writeButton(tableContents, extraData *crunch.Buffer, asset UgoAsset) {
	next := []byte(fmt.Sprintf("4\t%s\t%d\t%s\t%d\t%d\t%d\n", asset.Url, asset.Index1, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(asset.Name)), asset.Index2, asset.Index3, asset.Index4))

	if asset.Index1 == 0 {
		next = []byte(fmt.Sprintf("4\t%s\t%d\t%s\n", asset.Url, asset.Index1, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(asset.Name))))
	}

	tableContents.Grow(int64(len(next)))
	tableContents.WriteBytesNext(next)

	if strings.HasSuffix(asset.Url, ".ppm") {
		ppmFile, err := os.ReadFile("services/web/routes/res/bokeh.ppm")
		if err != nil {
			return
		}

		extraData.Grow(1696)
		extraData.WriteBytesNext(ppmFile[:1696])
	}

	if asset.Image != "" {
		imageFile, err := os.ReadFile(asset.Image)
		if err != nil {
			return
		}

		extraData.Grow(int64(len(imageFile)))
		extraData.WriteBytesNext(imageFile)
	}

}

func writeCategory(w *crunch.Buffer, asset UgoAsset) {
	next := []byte(fmt.Sprintf("2\t%s\t%s\t%d\n", asset.Url, base64.StdEncoding.EncodeToString(utils.WriteUTF16String(asset.Name)), asset.Index1))
	w.Grow(int64(len(next)))
	w.WriteBytesNext(next)
}
