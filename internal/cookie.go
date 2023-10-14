package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type Cookie struct {
	Domain               string `json:"domain"`
	IsIncludedSubDomains bool   `json:"isIncludedSubDomains"`
	Path                 string `json:"path"`
	HttpsOnly            bool   `json:"httpsOnly"`
	ExpiresAt            uint   `json:"expiresAt"`
	Name                 string `json:"name"`
	Value                string `json:"value"`
	Flags                uint32 `json:"flags"`
}

func NewCookie(domain, path, name, value string, expiresAt float64, flags uint32) *Cookie {
	return &Cookie{
		Domain:               domain,
		Path:                 path,
		Name:                 name,
		Value:                value,
		IsIncludedSubDomains: domain[0] == '.',
		HttpsOnly:            flags == 1 || flags == 5,
		ExpiresAt:            uint(expiresAt),
		Flags:                flags,
	}
}

type BinaryCookie struct {
	file *os.File
}

func NewBinaryCookie(path string) (*BinaryCookie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &BinaryCookie{
		file: file,
	}, nil
}

func (b *BinaryCookie) IsBinaryCookie() bool {
	magic_number := make([]byte, 4)
	b.file.Read(magic_number)
	return string(magic_number) == "cook"
}

func (b *BinaryCookie) ReadNumberOfPages() uint32 {
	var n uint32
	binary.Read(b.file, binary.BigEndian, &n)
	return n
}

func (b *BinaryCookie) ReadEachPageSize(n uint32) []uint32 {
	pageSizes := []uint32{}
	for i := uint32(0); i < n; i++ {
		var pageSize uint32
		binary.Read(b.file, binary.BigEndian, &pageSize)
		pageSizes = append(pageSizes, pageSize)
	}
	return pageSizes
}

func (b *BinaryCookie) String(pageSizes []uint32, format string) (string, error) {
	pages := [][]byte{}
	for _, v := range pageSizes {
		pageBuf := make([]byte, v)
		b.file.Read(pageBuf)
		pages = append(pages, pageBuf)
	}
	c := []*Cookie{}
	for _, v := range pages {

		buf := bytes.NewReader(v)

		var pageHeader uint32
		binary.Read(buf, binary.BigEndian, &pageHeader)
		if pageHeader != 0x100 {
			return "", errors.New("the page header has unexpected")
		}

		var numCookies uint32
		binary.Read(buf, binary.LittleEndian, &numCookies)
		cookieOffsets := []int64{}
		for i := uint32(0); i < numCookies; i++ {
			var cookieOffset int32
			binary.Read(buf, binary.LittleEndian, &cookieOffset)
			cookieOffsets = append(
				cookieOffsets,
				int64(cookieOffset),
			)
		}

		var pageFooter uint32
		binary.Read(buf, binary.LittleEndian, &pageFooter)
		if pageFooter != 0x0 {
			return "", errors.New("the page footer has unexpected")
		}

		for _, v := range cookieOffsets {

			buf.Seek(v, io.SeekStart)

			var cookieSize uint32
			binary.Read(buf, binary.LittleEndian, &cookieSize)

			cookie := make([]byte, cookieSize)
			buf.Read(cookie)
			bufCookie := bytes.NewReader(cookie)

			// Skip 4 bytes that is unknown
			bufCookie.Seek(4, io.SeekCurrent)

			var flags uint32
			binary.Read(bufCookie, binary.LittleEndian, &flags)

			// Skip 4 bytes that is unknown
			bufCookie.Seek(4, io.SeekCurrent)

			var urlOffset uint32
			binary.Read(bufCookie, binary.LittleEndian, &urlOffset)

			var nameOffset uint32
			binary.Read(bufCookie, binary.LittleEndian, &nameOffset)

			var pathOffset uint32
			binary.Read(bufCookie, binary.LittleEndian, &pathOffset)

			var valueOffset uint32
			binary.Read(bufCookie, binary.LittleEndian, &valueOffset)

			// Skip 8 bytes that is the end of cookie
			bufCookie.Seek(8, io.SeekCurrent)

			var expiry_date_epoch float64
			binary.Read(bufCookie, binary.LittleEndian, &expiry_date_epoch)
			expiry_date_epoch += 978307200

			var create_date_epoch float64
			binary.Read(bufCookie, binary.LittleEndian, &create_date_epoch)
			create_date_epoch += 978307200

			domain := ReadString(bufCookie)
			name := ReadString(bufCookie)
			path := ReadString(bufCookie)
			value := ReadString(bufCookie)

			c = append(c, NewCookie(domain, path, name, value, expiry_date_epoch, flags))
		}
	}
	if val, ok := writer[format]; ok {
		return val(c)
	} else {
		return "", errors.New("invalid output format")
	}
}
