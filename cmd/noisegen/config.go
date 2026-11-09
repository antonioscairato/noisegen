package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/antonioscairato/noisegen/noise"
)

// appConfig é o formato salvo em disco. Os valores são strings para manter
// compatibilidade com o arquivo gerado pela versão anterior do programa.
type appConfig struct {
	Volume string `json:"volume"`
	Tone   string `json:"tone"`
	Noise  string `json:"noise"`
}

// configPath retorna o caminho do arquivo onde o estado é salvo.
func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".app_config.json" // fallback: diretório atual
	}
	return filepath.Join(home, ".meu_app_config.json")
}

// saveConfig grava a configuração do gerador em disco.
func saveConfig(cfg noise.Config) error {
	data, err := json.MarshalIndent(appConfig{
		Volume: strconv.FormatFloat(cfg.Volume, 'f', 2, 64),
		Tone:   strconv.FormatFloat(cfg.Tone, 'f', 3, 64),
		Noise:  strconv.Itoa(int(cfg.Type)),
	}, "", "  ")
	if err != nil {
		return err
	}
	// Permissão de leitura/escrita apenas para o dono.
	return os.WriteFile(configPath(), data, 0600)
}

// loadConfig lê a configuração salva. Qualquer campo ausente ou inválido
// cai no valor padrão do package noise.
func loadConfig() noise.Config {
	cfg := noise.DefaultConfig()

	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	var stored appConfig
	if err := json.Unmarshal(data, &stored); err != nil {
		return cfg
	}

	if v, err := strconv.ParseFloat(stored.Volume, 64); err == nil {
		cfg.Volume = v
	}
	if t, err := strconv.ParseFloat(stored.Tone, 64); err == nil {
		cfg.Tone = t
	}
	if n, err := strconv.ParseInt(stored.Noise, 10, 32); err == nil {
		if nt := noise.Type(n); nt.Valid() {
			cfg.Type = nt
		}
	}
	return cfg
}
