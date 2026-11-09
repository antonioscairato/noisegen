// Package noise gera ruído sintético (branco, rosa e marrom) como PCM
// estéreo de 16 bits little-endian a 44100 Hz.
//
// O Generator implementa io.Reader, então pode ser plugado diretamente em
// qualquer saída de áudio que consuma PCM bruto, como o oto:
//
//	gen := noise.New(noise.Config{Type: noise.Brown, Volume: 0.5, Tone: 0.02})
//	player := ctx.NewPlayer(gen)
//
// Volume, tom e tipo podem ser alterados de qualquer goroutine enquanto o
// áudio está tocando. Read, porém, deve ser chamado por uma única goroutine.
package noise

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

// Parâmetros do formato de áudio produzido por Read.
const (
	SampleRate   = 44100
	ChannelCount = 2
	BitDepth     = 2 // bytes por amostra (int16)
)

// Limites e valores padrão dos controles.
const (
	MinVolume = 0.0
	MaxVolume = 1.0
	MinTone   = 0.005
	MaxTone   = 0.2

	DefaultVolume = 0.5
	DefaultTone   = 0.02
	DefaultType   = Brown
)

// Type identifica o tipo de ruído.
type Type int32

const (
	White Type = 0
	Pink  Type = 1
	Brown Type = 2
)

// String devolve o nome do tipo em inglês.
func (t Type) String() string {
	switch t {
	case White:
		return "white"
	case Pink:
		return "pink"
	case Brown:
		return "brown"
	default:
		return "unknown"
	}
}

// Valid informa se t é um tipo de ruído conhecido.
func (t Type) Valid() bool {
	return t == White || t == Pink || t == Brown
}

// Config reúne os parâmetros iniciais do gerador.
type Config struct {
	Type   Type
	Volume float64 // 0.0 a 1.0
	Tone   float64 // 0.005 a 0.2
}

// DefaultConfig devolve a configuração padrão.
func DefaultConfig() Config {
	return Config{Type: DefaultType, Volume: DefaultVolume, Tone: DefaultTone}
}

// Generator produz amostras de ruído. Use New para criar um.
type Generator struct {
	noiseType atomic.Int32
	volume    atomic.Uint64 // bits de um float64
	tone      atomic.Uint64 // bits de um float64

	rng        *rand.Rand
	lastBrownL float64
	lastBrownR float64
	pinkRowsL  [5]float64
	pinkRowsR  [5]float64
}

// New cria um Generator. Valores fora dos limites são ajustados, e um tipo
// inválido vira DefaultType.
func New(cfg Config) *Generator {
	g := &Generator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	g.SetType(cfg.Type)
	g.SetVolume(cfg.Volume)
	g.SetTone(cfg.Tone)
	return g
}

// Type devolve o tipo de ruído atual.
func (g *Generator) Type() Type { return Type(g.noiseType.Load()) }

// SetType altera o tipo de ruído. Tipos inválidos são ignorados.
func (g *Generator) SetType(t Type) {
	if !t.Valid() {
		t = DefaultType
	}
	g.noiseType.Store(int32(t))
}

// Volume devolve o volume atual (0.0 a 1.0).
func (g *Generator) Volume() float64 { return math.Float64frombits(g.volume.Load()) }

// SetVolume altera o volume, limitado a [MinVolume, MaxVolume].
func (g *Generator) SetVolume(v float64) {
	g.volume.Store(math.Float64bits(clamp(v, MinVolume, MaxVolume)))
}

// Tone devolve o coeficiente de tom atual (MinTone a MaxTone).
func (g *Generator) Tone() float64 { return math.Float64frombits(g.tone.Load()) }

// SetTone altera o tom (brilho), limitado a [MinTone, MaxTone].
func (g *Generator) SetTone(t float64) {
	g.tone.Store(math.Float64bits(clamp(t, MinTone, MaxTone)))
}

// Read preenche p com PCM estéreo int16 little-endian. Nunca devolve erro e
// sempre preenche o máximo de amostras completas que cabem em p.
func (g *Generator) Read(p []byte) (n int, err error) {
	const bytesPerFrame = ChannelCount * BitDepth

	vol := g.Volume()
	tCoef := g.Tone()
	nt := g.Type()

	for i := 0; i+bytesPerFrame <= len(p); i += bytesPerFrame {
		whiteL := g.rng.Float64()*2.0 - 1.0
		whiteR := g.rng.Float64()*2.0 - 1.0

		var outL, outR, gain float64

		switch nt {
		case White:
			outL, outR = whiteL, whiteR
			gain = 0.2
		case Pink:
			g.pinkRowsL[g.rng.Intn(5)] = whiteL
			g.pinkRowsR[g.rng.Intn(5)] = whiteR
			var sumL, sumR float64
			for j := 0; j < 5; j++ {
				sumL += g.pinkRowsL[j]
				sumR += g.pinkRowsR[j]
			}
			outL = (sumL / 5.0) * (tCoef * 15.0)
			outR = (sumR / 5.0) * (tCoef * 15.0)
			gain = 0.5
		case Brown:
			g.lastBrownL = (g.lastBrownL + tCoef*whiteL) * 0.99
			g.lastBrownR = (g.lastBrownR + tCoef*whiteR) * 0.99
			outL, outR = g.lastBrownL, g.lastBrownR
			gain = 1.0
		}

		outL = clamp(outL*vol*gain, -1.0, 1.0)
		outR = clamp(outR*vol*gain, -1.0, 1.0)

		sampleL := int16(outL * 32767.0)
		sampleR := int16(outR * 32767.0)

		p[i] = byte(sampleL)
		p[i+1] = byte(sampleL >> 8)
		p[i+2] = byte(sampleR)
		p[i+3] = byte(sampleR >> 8)

		n += bytesPerFrame
	}
	return n, nil
}

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}
