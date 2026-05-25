package missions

import "crypto-desert/internal/enemies"

// ── Difficulty ────────────────────────────────────────────────────────────────

type Difficulty int

const (
	DifficultyNormal Difficulty = 0
	DifficultyNGPlus Difficulty = 1
	DifficultyNGPP   Difficulty = 2
	DifficultyNGPPP  Difficulty = 3
)

func (d Difficulty) Label() string {
	switch d {
	case DifficultyNormal: return "Normal"
	case DifficultyNGPlus: return "NG+"
	case DifficultyNGPP:   return "NG++"
	case DifficultyNGPPP:  return "NG+++"
	default:               return "NG+++"
	}
}

func (d Difficulty) StatMultiplier() float64 { return 1.0 + float64(d)*0.25 }
func (d Difficulty) XPMultiplier() float64   { return 1.0 + float64(d)*0.5 }

// ── Wave ──────────────────────────────────────────────────────────────────────

type Wave struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Intro      string   `json:"intro"`
	EnemyNames []string `json:"enemy_names"`
	IsBossWave bool     `json:"is_boss_wave"`
}

// ── Mission ───────────────────────────────────────────────────────────────────

type Mission struct {
	ID     string `json:"id"`
	CityID string `json:"city_id"`
	Title  string `json:"title"`
	Waves  []Wave `json:"waves"`
}

// ── City ──────────────────────────────────────────────────────────────────────

type City struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Subtitle   string  `json:"subtitle"`
	Faction    string  `json:"faction"`
	Lore       string  `json:"lore"`
	Problem    string  `json:"problem"`
	Reward     string  `json:"reward"`
	Icon       string  `json:"icon"`
	Mission    Mission `json:"mission"`
	UnlockedBy string  `json:"unlocked_by"`
}

// ── Campaign ──────────────────────────────────────────────────────────────────

var Campaign = []City{
	{
		ID: "genesis_block", Name: "Genesis Block", Subtitle: "Capital da Ordem dos Blocos",
		Faction: "BTC", Icon: "⛓",
		Lore: `Em 2087, o colapso dos governos deixou um vácuo de poder preenchido por facções que controlam blockchains privadas. Genesis Block foi a primeira cidade construída sobre a cadeia original — o Bitcoin. Seus habitantes são mineradores, validadores e guardiões do protocolo primordial. Por décadas, Genesis Block prosperou como a cidade mais estável do deserto digital. Mas algo mudou.`,
		Problem: `Especuladores e bots de pump estão drenando a liquidez da cidade, manipulando os preços do BTC local e corrompendo os mineradores mais fracos. O conselho da Ordem dos Blocos enviou um pedido de socorro urgente.`,
		Reward:  `A Ordem dos Blocos está em dívida com você. As rotas para Ether Citadel foram abertas. O primeiro bloco está seguro.`,
		UnlockedBy: "",
		Mission: Mission{
			ID: "mission_genesis", CityID: "genesis_block", Title: "Operação: Cadeia Limpa",
			Waves: []Wave{
				{ID: "genesis_w1", Title: "Encontro 1 — Portões Externos", EnemyNames: []string{"Especulador Novato", "Especulador Novato"}, Intro: "Dois especuladores novatos bloqueiam a entrada da cidade. Eles parecem desesperados — e perigosos."},
				{ID: "genesis_w2", Title: "Encontro 2 — Sala dos Servidores", EnemyNames: []string{"Bot de Pump"}, Intro: "Um Bot de Pump corrompido assumiu o controle dos servidores de mineração. Seus ataques são aleatórios e imprevisíveis."},
				{ID: "genesis_w3", Title: "Encontro 3 — Núcleo da Rede", EnemyNames: []string{"Minerador Fantasma"}, Intro: "O Minerador Fantasma guarda o núcleo central. Ele serviu à Ordem por décadas — mas foi corrompido pelo caos do mercado."},
				{ID: "genesis_boss", Title: "BOSS — Câmara do Bloco Genesis", EnemyNames: []string{"Satoshi das Trevas"}, Intro: "No coração de Genesis Block, uma entidade sombria aguarda. Satoshi das Trevas corrompeu o bloco original.", IsBossWave: true},
			},
		},
	},
	{
		ID: "ether_citadel", Name: "Ether Citadel", Subtitle: "Sede do Conclave dos Contratos",
		Faction: "ETH", Icon: "🔮",
		Lore: `Ether Citadel é uma cidade vertical — torres de contratos inteligentes que se estendem até a névoa digital. O Conclave dos Contratos governa com precisão matemática: cada lei é um contrato, cada punição é uma execução automática. Era uma cidade de ordem perfeita — até os oráculos mentirem.`,
		Problem: `Os oráculos de dados foram comprometidos por um agente desconhecido. Contratos maliciosos estão sendo executados automaticamente, drenando fundos e paralisando os cidadãos. O Conclave perdeu o controle dos andares inferiores.`,
		Reward:  `Os oráculos estão limpos. Os contratos voltam a executar com dados verdadeiros. A rota para Sol Dunes foi desbloqueada pelo Conclave como agradecimento.`,
		UnlockedBy: "genesis_block",
		Mission: Mission{
			ID: "mission_ether", CityID: "ether_citadel", Title: "Operação: Dados Verdadeiros",
			Waves: []Wave{
				{ID: "ether_w1", Title: "Encontro 1 — Andar Térreo", EnemyNames: []string{"Fomo Cultist", "Fomo Cultist"}, Intro: "Cultistas do FOMO tomaram o lobby da Citadel. Eles gritam previsões absurdas de preço enquanto atacam."},
				{ID: "ether_w2", Title: "Encontro 2 — Câmara dos Oráculos", EnemyNames: []string{"Oráculo Corrompido"}, Intro: "O Oráculo Corrompido manipula dados em tempo real. Ele já escreveu um contrato para contra-atacar."},
				{ID: "ether_w3", Title: "Encontro 3 — Mempool Proibido", EnemyNames: []string{"Sombra do Mempool"}, Intro: "A Sombra do Mempool existe entre as transações pendentes. Você mal consegue vê-la antes do veneno entrar."},
				{ID: "ether_boss", Title: "BOSS — Coração do Contrato", EnemyNames: []string{"Vitalik Void"}, Intro: "Vitalik Void, o arquiteto que perdeu o controle de sua própria criação. O contrato que ele escreveu agora o controla.", IsBossWave: true},
			},
		},
	},
	{
		ID: "sol_dunes", Name: "Sol Dunes", Subtitle: "Território dos Rastreadores Solares",
		Faction: "SOL", Icon: "🏜",
		Lore: `Sol Dunes é uma cidade construída sobre velocidade. Os Rastreadores Solares processam milhares de transações por segundo. Mas Dust Raiders descobriram como injetar transações falsas na rede de alta velocidade — e tudo travou.`,
		Problem: `Dust Raiders estão injetando transações falsas na rede SOL, causando congestionamento e corrompendo validadores. A velocidade da rede caiu 80% — um colapso sem precedentes.`,
		Reward:  `A rede SOL voltou à velocidade máxima. Os Rastreadores abriram as rotas para BNB Quarter.`,
		UnlockedBy: "ether_citadel",
		Mission: Mission{
			ID: "mission_sol", CityID: "sol_dunes", Title: "Operação: Alta Frequência",
			Waves: []Wave{
				{ID: "sol_w1", Title: "Encontro 1 — Dunas Externas", EnemyNames: []string{"Dust Raider", "Dust Raider"}, Intro: "Dust Raiders emboscam você nas dunas antes mesmo de chegar à cidade."},
				{ID: "sol_w2", Title: "Encontro 2 — Torre de Validação", EnemyNames: []string{"Validador Traidor"}, Intro: "Um Validador Traidor vendeu seu nó para os Raiders. Ele conhece todos os seus padrões."},
				{ID: "sol_w3", Title: "Encontro 3 — Núcleo de Throughput", EnemyNames: []string{"Whale Corrupta"}, Intro: "A Whale Corrupta financia toda a operação. Imensa, furiosa e com bolsos fundos."},
				{ID: "sol_boss", Title: "BOSS — Câmara Central do Protocolo", EnemyNames: []string{"O Liquidador"}, Intro: "O Liquidador chegou antes de você. Um protocolo de liquidação forçada que ganhou consciência.", IsBossWave: true},
			},
		},
	},
	{
		ID: "bnb_quarter", Name: "BNB Quarter", Subtitle: "Domínio da Guilda das Taxas",
		Faction: "BNB", Icon: "💰",
		Lore: `BNB Quarter é o coração comercial do deserto digital. Cada transação paga uma taxa à Guilda das Taxas. Era um sistema que funcionava para todos — até alguém decidir que preferia um sistema que funciona só para ele.`,
		Problem: `Um esquema sofisticado de front-running está desviando as taxas da Guilda para carteiras fantasma. As Sombras do Mempool são os executores — mas alguém mais poderoso as coordena.`,
		Reward:  `A Guilda das Taxas recuperou suas receitas. As rotas para DOGE Wasteland foram abertas.`,
		UnlockedBy: "sol_dunes",
		Mission: Mission{
			ID: "mission_bnb", CityID: "bnb_quarter", Title: "Operação: Taxa Justa",
			Waves: []Wave{
				{ID: "bnb_w1", Title: "Encontro 1 — Beco das Trocas", EnemyNames: []string{"Bot de Pump", "Especulador Novato"}, Intro: "Bots de Pump tomaram o beco onde ficam as exchanges menores."},
				{ID: "bnb_w2", Title: "Encontro 2 — Câmara de Compensação", EnemyNames: []string{"Sombra do Mempool"}, Intro: "A Sombra do Mempool intercepta transações em tempo real."},
				{ID: "bnb_w3", Title: "Encontro 3 — Cofre da Guilda", EnemyNames: []string{"Oráculo Corrompido"}, Intro: "O Oráculo Corrompido alimenta dados falsos para os sistemas de auditoria."},
				{ID: "bnb_boss", Title: "BOSS — Sala dos Registros", EnemyNames: []string{"O Liquidador"}, Intro: "O Liquidador foi contratado pela Guilda rival para eliminar você e apagar as evidências.", IsBossWave: true},
			},
		},
	},
	{
		ID: "doge_wasteland", Name: "DOGE Wasteland", Subtitle: "Terra da Horda Lunar — Fim do Mundo",
		Faction: "DOGE", Icon: "🌕",
		Lore: `DOGE Wasteland não foi planejada. Simplesmente surgiu — um acidente de mercado que ganhou vida própria. A Horda Lunar não obedece lógica, não respeita contratos, não segue padrões. E agora o DOGE Primordial acordou.`,
		Problem: `O DOGE Primordial está corroendo a estabilidade de todo o deserto digital. Cada hora que passa, os preços de todas as cryptos ficam mais voláteis. Você precisa chegar até ele e encerrar isso.`,
		Reward:  `O DOGE Primordial foi contido. O deserto digital respira. Você limpou todas as cidades — o NG+ começa.`,
		UnlockedBy: "bnb_quarter",
		Mission: Mission{
			ID: "mission_doge", CityID: "doge_wasteland", Title: "Operação: Lua Final",
			Waves: []Wave{
				{ID: "doge_w1", Title: "Encontro 1 — Ruínas da Exchange", EnemyNames: []string{"Fomo Cultist", "Especulador Novato", "Bot de Pump"}, Intro: "Fomo Cultists, Especuladores e Bots do Wasteland. Caóticos e assustados com o que acordou aqui."},
				{ID: "doge_w2", Title: "Encontro 2 — Campo de Memes", EnemyNames: []string{"Whale Corrupta", "Sombra do Mempool"}, Intro: "Uma Whale Corrupta e uma Sombra do Mempool formaram uma aliança improvável."},
				{ID: "doge_w3", Title: "Encontro 3 — Portal do Caos", EnemyNames: []string{"Validador Traidor", "Oráculo Corrompido"}, Intro: "O Validador Traidor e o Oráculo Corrompido alimentam o ritual que mantém o DOGE Primordial desperto."},
				{ID: "doge_boss", Title: "BOSS FINAL — Núcleo do Primordial", EnemyNames: []string{"DOGE Primordial"}, Intro: "O DOGE Primordial em sua forma completa. O meme que virou realidade, o caos que virou carne.", IsBossWave: true},
			},
		},
	},
}

func CityByID(id string) (City, bool) {
	for _, c := range Campaign {
		if c.ID == id { return c, true }
	}
	return City{}, false
}

func NextCity(cityID string) (City, bool) {
	for i, c := range Campaign {
		if c.ID == cityID && i+1 < len(Campaign) { return Campaign[i+1], true }
	}
	return City{}, false
}

func ScaleTemplate(tmpl enemies.EnemyTemplate, d Difficulty) enemies.EnemyTemplate {
	if d == DifficultyNormal { return tmpl }
	mult := d.StatMultiplier()
	scaled := tmpl
	scaled.Level = int(float64(tmpl.Level) * mult)
	if tmpl.HPOverride > 0 { scaled.HPOverride = int(float64(tmpl.HPOverride) * mult) }
	if tmpl.AttackModOverride > 0 { scaled.AttackModOverride = int(float64(tmpl.AttackModOverride) * mult) }
	scaled.XPReward = int(float64(tmpl.XPReward) * d.XPMultiplier())
	scaled.GoldReward = int(float64(tmpl.GoldReward) * mult)
	return scaled
}
