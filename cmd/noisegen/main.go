// Comando noisegen: gerador de ruídos sintéticos (branco, rosa e marrom)
// com controle pelo teclado no terminal.
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ebitengine/oto/v3"
	"github.com/eiannone/keyboard"

	"github.com/antonioscairato/noisegen/noise"
)

const (
	volumeStep = 0.05
	toneStep   = 0.005
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   noise.SampleRate,
		ChannelCount: noise.ChannelCount,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		fmt.Printf("Erro no áudio: %v\n", err)
		return 1
	}
	<-ready

	gen := noise.New(loadConfig())
	player := ctx.NewPlayer(gen)
	player.Play()
	defer player.Close()

	if err := keyboard.Open(); err != nil {
		fmt.Printf("Erro ao iniciar teclado: %v\n", err)
		return 1
	}
	defer keyboard.Close()

	// Com o terminal em modo raw o Ctrl+C chega como tecla (tratado abaixo),
	// mas mantemos o handler de sinal para o caso de o SO ainda entregá-lo.
	quit := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		close(quit)
		keyboard.Close() // desbloqueia keyboard.GetKey
	}()

	renderUI(gen)

loop:
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}
		select {
		case <-quit:
			break loop
		default:
		}

		changed := false

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			break loop
		case keyboard.KeyArrowUp:
			gen.SetVolume(gen.Volume() + volumeStep)
			changed = true
		case keyboard.KeyArrowDown:
			gen.SetVolume(gen.Volume() - volumeStep)
			changed = true
		case keyboard.KeyArrowRight: // tom mais agudo/brilhante
			gen.SetTone(gen.Tone() + toneStep)
			changed = true
		case keyboard.KeyArrowLeft: // tom mais abafado/grave
			gen.SetTone(gen.Tone() - toneStep)
			changed = true
		}

		switch char {
		case '1':
			gen.SetType(noise.White)
			changed = true
		case '2':
			gen.SetType(noise.Pink)
			changed = true
		case '3':
			gen.SetType(noise.Brown)
			changed = true
		}

		if !changed {
			continue
		}

		renderUI(gen)
		err = saveConfig(noise.Config{
			Type:   gen.Type(),
			Volume: gen.Volume(),
			Tone:   gen.Tone(),
		})
		if err != nil {
			fmt.Println("Erro ao salvar preferências:", err)
		}
	}

	fmt.Print(clearScreen)
	fmt.Println("👋🎧 Até logo!")
	return 0
}
