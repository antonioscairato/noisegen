package main

import (
	"fmt"

	"github.com/antonioscairato/noisegen/noise"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	italic  = "\033[3m"
	faint   = "\033[2m"
	magenta = "\033[95m"
	brown   = "\033[38;5;94m"
	white   = "\033[97m"

	clearScreen = "\033[H\033[2J"
)

type noiseInfo struct {
	title       string
	color       string
	description string
}

var noiseInfos = map[noise.Type]noiseInfo{
	noise.White: {
		title: "BRANCO (Estático/TV Antiga)",
		color: white,
		description: `
Contém todas as frequências audíveis na mesma intensidade.
O som se assemelha a uma TV fora do ar ou a um ventilador forte:

• Bloqueio de barulhos externos: Excelente para mascarar sons
repentinos (como latidos de cachorro, trânsito ou portas batendo)
porque preenche o ambiente.
• Sono de bebês: Muito utilizado para acalmar recém-nascidos,
pois simula o som uterino e bloqueia os ruídos da casa.
• Alívio de zumbido (Tinnitus): Ajuda a neutralizar o chiado
constante no ouvido de quem sofre com a condição.`,
	},
	noise.Pink: {
		title: "ROSA (Chuva/Cachoeira Equilibrada)",
		color: magenta,
		description: `
Tem frequências graves mais altas e agudos mais
suaves que o ruído branco, resultando em um som
mais natural, como o de uma chuva constante ou
vento nas árvores.

• Melhoria do sono profundo: Estudos mostram que ele
sincroniza as ondas cerebrais, ajudando a prolongar o
sono de ondas lentas (estágio mais reparador).
• Consolidação da memória: Por melhorar o sono profundo,
está diretamente ligado à retenção de memória e aprendizado.
• Pessoas sensíveis a agudos: É mais confortável e menos
"irritante" aos ouvidos do que o ruído branco.`,
	},
	noise.Brown: {
		title: "MARROM (Mar Profundo/Trovão Distante)",
		color: brown,
		description: `
Dá ainda mais destaque aos sons graves e reduz drasticamente
os agudos. O som é profundo e robusto, lembrando uma cachoeira
distante ou o interior de um avião.

• Foco profundo e TDAH: Muito procurado por pessoas que precisam
se concentrar em tarefas complexas ou que têm TDAH, pois ajuda a
"acalmar" o cérebro hiperativo sem causar distração.
• Relaxamento e redução de ansiedade: O tom grave e aveludado
gera uma sensação imediata de aconchego e diminuição do estresse.
• Meditação: Ajuda a silenciar os pensamentos internos
(o "falatório" mental).`,
	},
}

// renderUI limpa a tela e desenha o painel de controle atualizado.
func renderUI(gen *noise.Generator) {
	info := noiseInfos[gen.Type()]
	vol := gen.Volume()
	tone := gen.Tone()

	fmt.Print(clearScreen)
	fmt.Println("=========================================================")
	fmt.Println(bold + "          🎧  GERADOR DE RUÍDOS SINTÉTICOS  🎧" + reset)
	fmt.Println("=========================================================")
	fmt.Printf(" [1] Tipo de Ruído :  %s\n", bold+info.color+info.title+reset)
	fmt.Printf(" [2] Volume        :  %-3.0f%%  (Use ↑ / ↓  para ajustar)\n", vol*100)
	fmt.Printf(" [3] Tom (Brilho)  :  %-3.0f%%  (Use ← / →  para ajustar)\n", (tone/noise.MaxTone)*100)
	fmt.Println("=========================================================")
	fmt.Println(" Teclas de Atalho:")
	fmt.Println("   [1] Ruído Branco  |  [2] Ruído Rosa  |  [3] Ruído Marrom")
	fmt.Println("   [ESC] ou [Ctrl+C] para Sair")
	fmt.Println("=========================================================")
	fmt.Println(italic + info.description + reset)
	fmt.Println("---------------------------------------------------------")
	fmt.Println(faint + italic + "Antonio Scairato - Software Developer (antonio.scairato@gmail.com)" + reset)
}
