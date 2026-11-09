# noisegen

Gerador de ruídos sintéticos (branco, rosa e marrom) em Go, com controle pelo teclado no terminal.

O projeto tem duas partes:

- **`noise/`**: package reutilizável que gera o áudio (PCM estéreo, 16 bits, 44100 Hz) e implementa `io.Reader`.
- **`cmd/noisegen/`**: aplicativo de terminal que toca o ruído via [oto](https://github.com/ebitengine/oto) e salva as preferências do usuário.

## Aplicativo

```sh
go install github.com/antonioscairato/noisegen/cmd/noisegen@latest
noisegen
```

| Tecla | Ação |
|-------|------|
| `1` / `2` / `3` | Ruído branco / rosa / marrom |
| `↑` / `↓` | Volume |
| `←` / `→` | Tom (brilho) |
| `Esc` ou `Ctrl+C` | Sair |

Volume, tom e tipo de ruído são salvos em `~/.meu_app_config.json` e restaurados na próxima execução.

No Windows, baixe o `noisegen-windows-amd64.exe` na página de [Releases](https://github.com/antonioscairato/noisegen/releases) e execute-o em um terminal (Windows Terminal ou PowerShell).

No Linux, o oto precisa das bibliotecas de desenvolvimento do ALSA (por exemplo, `libasound2-dev` no Debian/Ubuntu).

## Usando o package

```go
import "github.com/antonioscairato/noisegen/noise"

gen := noise.New(noise.Config{
	Type:   noise.Brown,
	Volume: 0.5,
	Tone:   0.02,
})

// gen é um io.Reader com PCM estéreo int16 little-endian a 44100 Hz.
player := otoCtx.NewPlayer(gen)
player.Play()

// Ajustes em tempo real, de qualquer goroutine:
gen.SetVolume(0.3)
gen.SetType(noise.Pink)
gen.SetTone(0.05)
```

`Read` deve ser chamado por uma única goroutine; os setters são seguros para uso concorrente.

## Testes

```sh
go test ./noise/
```
