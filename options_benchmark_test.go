package logger

import (
	"bytes"
	"io"
	"testing"
	"time"
)

var (
	input = []byte(`Caliban: I must eat my dinner.
	This island's mine, by Sycorax my mother,
	Which thou takest from me. When thou camest first,
	Thou strokedst me and madest much of me, wouldst give me
	Water with berries in't, and teach me how
	To name the bigger light, and how the less,
	That burn by day and night: and then I loved thee
	And show'd thee all the qualities o' the isle,
	The fresh springs, brine-pits, barren place and fertile:
	Cursed be I that did so! All the charms
	Of Sycorax, toads, beetles, bats, light on you!
	For I am all the subjects that you have,
	Which first was mine own king: and here you sty me
	In this hard rock, whiles you do keep from me
	The rest o' the island.
Prospero: Thou most lying slave,
	Whom stripes may move, not kindness! I have used thee,
	Filth as thou art, with human care, and lodged thee
	In mine own cell, till thou didst seek to violate
	The honour of my child.

Caliban: O ho, O ho! would't had been done!
	Thou didst prevent me; I had peopled else
	This isle with Calibans. 
	Caliban: I must eat my dinner.
	This island's mine, by Sycorax my mother,
	Which thou takest from me. When thou camest first,
	Thou strokedst me and madest much of me, wouldst give me
	Water with berries in't, and teach me how
	To name the bigger light, and how the less,
	That burn by day and night: and then I loved thee
	And show'd thee all the qualities o' the isle,
	The fresh springs, brine-pits, barren place and fertile:
	Cursed be I that did so! All the charms
	Of Sycorax, toads, beetles, bats, light on you!
	For I am all the subjects that you have,
	Which first was mine own king: and here you sty me
	In this hard rock, whiles you do keep from me
	The rest o' the island.
Prospero: Thou most lying slave,
	Whom stripes may move, not kindness! I have used thee,
	Filth as thou art, with human care, and lodged thee
	In mine own cell, till thou didst seek to violate
	The honour of my child.

Caliban: O ho, O ho! would't had been done!
	Thou didst prevent me; I had peopled else
	This isle with Calibans. 
	Caliban: I must eat my dinner.
	This island's mine, by Sycorax my mother,
	Which thou takest from me. When thou camest first,
	Thou strokedst me and madest much of me, wouldst give me
	Water with berries in't, and teach me how
	To name the bigger light, and how the less,
	That burn by day and night: and then I loved thee
	And show'd thee all the qualities o' the isle,
	The fresh springs, brine-pits, barren place and fertile:
	Cursed be I that did so! All the charms
	Of Sycorax, toads, beetles, bats, light on you!
	For I am all the subjects that you have,
	Which first was mine own king: and here you sty me
	In this hard rock, whiles you do keep from me
	The rest o' the island.
Prospero: Thou most lying slave,
	Whom stripes may move, not kindness! I have used thee,
	Filth as thou art, with human care, and lodged thee
	In mine own cell, till thou didst seek to violate
	The honour of my child.

`)
)

func Benchmark_writeByNewLineSync(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		writers := []io.Writer{
			&slowWriter{latency: time.Second}, &slowWriter{latency: time.Millisecond * 500}, &slowWriter{latency: 3 * time.Second},
		}
		r := bytes.NewReader(input)
		debugger := &testDebugger{noPrint: true}
		err := writeByNewLineSync(debugger, r, writers...)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_writeByNewLine(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		writers := []io.Writer{
			&slowWriter{latency: time.Second}, &slowWriter{latency: time.Millisecond * 500}, &slowWriter{latency: 3 * time.Second},
		}
		r := bytes.NewReader(input)
		debugger := &testDebugger{noPrint: true}
		err := writeByNewLine(debugger, r, writers...)
		if err != nil {
			b.Error(err)
		}
	}
}
