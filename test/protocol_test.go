package test

import (
	"agent/src/g"
	"log"
	"testing"
)

func TestProtocol(t *testing.T) {
	p := g.Protocol{}

	pkt := g.NewPkt()
	pkt.Data = []byte("hi,world")
	d, _ := pkt.Serialize()
	log.Println(pkt.Len())
	var data []byte
	data = append(data, d...)
	data = append(data, d...)
	for {
		log.Println(len(data))
		a, rawData, err := p.UnPack(data, false)
		data = data[a:]
		if err != nil {
			t.Error(err)
			return
		}
		log.Println(len(rawData))
		tp, err := p.Decode(rawData)
		if err != nil {
			t.Error(err)
			return
		}
		pkt = tp.(*g.Packet)
		log.Println(string(pkt.Data))
	}

}
