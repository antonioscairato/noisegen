package noise

import (
	"encoding/binary"
	"testing"
)

func samples(t *testing.T, g *Generator, frames int) []int16 {
	t.Helper()
	buf := make([]byte, frames*ChannelCount*BitDepth)
	n, err := g.Read(buf)
	if err != nil {
		t.Fatalf("Read retornou erro: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("Read preencheu %d bytes, esperado %d", n, len(buf))
	}
	out := make([]int16, len(buf)/2)
	for i := range out {
		out[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return out
}

func TestReadFillsWholeBuffer(t *testing.T) {
	for _, nt := range []Type{White, Pink, Brown} {
		g := New(Config{Type: nt, Volume: 0.5, Tone: 0.02})
		samples(t, g, 1024)
	}
}

func TestReadIgnoresPartialFrame(t *testing.T) {
	g := New(DefaultConfig())
	buf := make([]byte, 10) // 2 frames completos + 2 bytes sobrando
	n, _ := g.Read(buf)
	if n != 8 {
		t.Fatalf("n = %d, esperado 8", n)
	}
}

func TestZeroVolumeIsSilence(t *testing.T) {
	for _, nt := range []Type{White, Pink, Brown} {
		g := New(Config{Type: nt, Volume: 0, Tone: 0.02})
		for _, s := range samples(t, g, 512) {
			if s != 0 {
				t.Fatalf("%v: amostra %d com volume 0", nt, s)
			}
		}
	}
}

func TestNonZeroOutput(t *testing.T) {
	for _, nt := range []Type{White, Pink, Brown} {
		g := New(Config{Type: nt, Volume: 1, Tone: 0.05})
		var sum int64
		for _, s := range samples(t, g, 4096) {
			if s < 0 {
				sum -= int64(s)
			} else {
				sum += int64(s)
			}
		}
		if sum == 0 {
			t.Fatalf("%v: saída toda zerada", nt)
		}
	}
}

func TestSettersClamp(t *testing.T) {
	g := New(DefaultConfig())

	g.SetVolume(5)
	if g.Volume() != MaxVolume {
		t.Errorf("Volume = %v, esperado %v", g.Volume(), MaxVolume)
	}
	g.SetVolume(-1)
	if g.Volume() != MinVolume {
		t.Errorf("Volume = %v, esperado %v", g.Volume(), MinVolume)
	}
	g.SetTone(10)
	if g.Tone() != MaxTone {
		t.Errorf("Tone = %v, esperado %v", g.Tone(), MaxTone)
	}
	g.SetTone(0)
	if g.Tone() != MinTone {
		t.Errorf("Tone = %v, esperado %v", g.Tone(), MinTone)
	}
}

func TestInvalidTypeFallsBackToDefault(t *testing.T) {
	g := New(Config{Type: Type(42), Volume: 0.5, Tone: 0.02})
	if g.Type() != DefaultType {
		t.Fatalf("Type = %v, esperado %v", g.Type(), DefaultType)
	}
}

func TestSetType(t *testing.T) {
	g := New(DefaultConfig())
	g.SetType(White)
	if g.Type() != White {
		t.Fatalf("Type = %v, esperado %v", g.Type(), White)
	}
}

func TestConcurrentSetters(t *testing.T) {
	g := New(DefaultConfig())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1000; i++ {
			g.SetVolume(float64(i%100) / 100)
			g.SetTone(0.005 + float64(i%10)/100)
			g.SetType(Type(i % 3))
		}
	}()
	buf := make([]byte, 4096)
	for i := 0; i < 200; i++ {
		g.Read(buf)
	}
	<-done
}
