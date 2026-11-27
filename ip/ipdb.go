// Copyright 2026 Casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ip

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"strings"
)

// ipip.net .ipdb 格式解析器（ipipfree.ipdb 与商业版通用）。
//
// 布局：4 字节元数据长度 + 元数据 JSON + 二叉前缀树（每节点 8 字节，两个大端 uint32）
// + 数据区（每条记录 2 字节长度 + 制表符分隔的字段）。叶节点值为 nodeCount + 数据区字节偏移。

type ipdbHeader struct {
	IPVersion uint16         `json:"ip_version"`
	Languages map[string]int `json:"languages"`
	NodeCount int            `json:"node_count"`
	TotalSize int            `json:"total_size"`
	Fields    []string       `json:"fields"`
}

// IPDBLocator 解析 .ipdb 格式的 IP 库
type IPDBLocator struct {
	data      []byte
	nodeCount int
	v4Offset  int
	langOff   int
	fieldCnt  int
}

// isIPDB ipdb 文件以 4 字节元数据长度开头，紧接着是 JSON（形如 {"build":...}）；
// 17mon .dat 的前 4 字节是文本区偏移，偏移之后是索引字节，不会出现 '{"'。
func isIPDB(data []byte) bool {
	return len(data) > 5 && bytes.HasPrefix(data[4:], []byte(`{"`))
}

func NewIPDBLocator(data []byte) (*IPDBLocator, error) {
	if len(data) < 4 {
		return nil, errors.New("ipdb too short")
	}
	metaLen := int(binary.BigEndian.Uint32(data[:4]))
	if metaLen <= 0 || 4+metaLen > len(data) {
		return nil, errors.New("ipdb has an invalid meta length")
	}

	var header ipdbHeader
	if err := json.Unmarshal(data[4:4+metaLen], &header); err != nil {
		return nil, err
	}
	if header.NodeCount <= 0 || len(header.Fields) == 0 {
		return nil, errors.New("ipdb meta is incomplete")
	}
	if header.TotalSize <= 0 || 4+metaLen+header.TotalSize > len(data) {
		return nil, errors.New("ipdb file is truncated")
	}
	if header.IPVersion&1 == 0 {
		return nil, errors.New("ipdb has no ipv4 data")
	}

	loc := &IPDBLocator{
		data:      data[4+metaLen:],
		nodeCount: header.NodeCount,
		langOff:   header.Languages["CN"],
		fieldCnt:  len(header.Fields),
	}

	// IPv4 子树入口：先沿 80 个 0 位、再沿 16 个 1 位向下走
	node := 0
	for i := 0; i < 96 && node >= 0 && node < header.NodeCount; i++ {
		index := 0
		if i >= 80 {
			index = 1
		}
		node = loc.readNode(node, index)
	}
	if node < 0 || node >= header.NodeCount {
		return nil, errors.New("ipdb has no valid ipv4 subtree")
	}
	loc.v4Offset = node
	return loc, nil
}

func (loc *IPDBLocator) readNode(node, index int) int {
	off := node*8 + index*4
	if node < 0 || off+4 > len(loc.data) {
		return -1
	}
	return int(binary.BigEndian.Uint32(loc.data[off : off+4]))
}

// resolve 叶节点值 → 数据区记录内容
func (loc *IPDBLocator) resolve(node int) ([]byte, bool) {
	off := loc.nodeCount*8 + node - loc.nodeCount
	if off < 0 || off+2 > len(loc.data) {
		return nil, false
	}
	size := int(binary.BigEndian.Uint16(loc.data[off : off+2]))
	if off+2+size > len(loc.data) {
		return nil, false
	}
	return loc.data[off+2 : off+2+size], true
}

// Find locationInfo by ip string
func (loc *IPDBLocator) Find(ipstr string) (*LocationInfo, error) {
	ip := net.ParseIP(ipstr).To4()
	if ip == nil {
		return nil, ErrInvalidIp
	}
	info := loc.FindByUint(binary.BigEndian.Uint32(ip))
	if info == nil {
		return nil, ErrInvalidIp
	}
	return info, nil
}

// Find locationInfo by uint32, nil when the IP is not covered by the database
func (loc *IPDBLocator) FindByUint(ip uint32) *LocationInfo {
	node := loc.v4Offset
	for i := 0; i < 32; i++ {
		if node < 0 || node >= loc.nodeCount {
			return nil
		}
		node = loc.readNode(node, int(ip>>uint(31-i)&1))
	}
	if node <= loc.nodeCount {
		return nil
	}

	raw, ok := loc.resolve(node)
	if !ok {
		return nil
	}
	fields := strings.Split(string(raw), "\t")
	if loc.langOff < 0 || loc.langOff+loc.fieldCnt > len(fields) || loc.fieldCnt < 3 {
		return nil
	}
	fields = fields[loc.langOff : loc.langOff+loc.fieldCnt]
	return &LocationInfo{
		Country: fields[0],
		Region:  fields[1],
		City:    fields[2],
	}
}
