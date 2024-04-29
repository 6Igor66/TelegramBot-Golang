package telegram

import (
	"bot/internal/config"
	"bot/internal/models"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log"
	"time"

	"github.com/cameo-engineering/tonconnect"
	"github.com/skip2/go-qrcode"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot    *tele.Bot
	wallet string
}

func NewBot(cfg *config.Config) *Bot {
	var b Bot
	b.initBot(cfg)
	b.initHandlers()
	return &b
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) initBot(cfg *config.Config) {
	settings := tele.Settings{
		Token:  cfg.Bot.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	var err error
	b.bot, err = tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}
	b.wallet = cfg.Bot.Wallet
}

func (b *Bot) initHandlers() {
	var (
		menu = &tele.ReplyMarkup{ResizeKeyboard: true}

		btnWallet      = menu.Text("Wallet")
		btnTonKeeper   = menu.Text("Tonkeeper")
		btnTonHub      = menu.Text("Tonhub")
		btnMyTonWallet = menu.Text("MyTonWallet")
	)

	menu.Reply(
		menu.Row(btnWallet),
		menu.Row(btnTonKeeper),
		menu.Row(btnTonHub),
		menu.Row(btnMyTonWallet),
	)

	b.bot.Handle("/start", func(c tele.Context) error {
		s, err := tonconnect.NewSession()

		UsersState[c.Chat().ID] = &models.UserState{
			S:         s,
			Connected: false,
		}

		if err != nil {
			log.Fatal(err)
		}

		data := make([]byte, 32)
		_, err = rand.Read(data)
		if err != nil {
			log.Fatal(err)
		}

		connreq, err := tonconnect.NewConnectRequest(
			"https://raw.githubusercontent.com/XaBbl4/pytonconnect/main/pytonconnect-manifest.json", //TODO create own manifest
			tonconnect.WithProofRequest(base32.StdEncoding.EncodeToString(data)),
		)
		if err != nil {
			log.Fatal(err)
		}

		for _, wallet := range tonconnect.Wallets {
			link, err := s.GenerateUniversalLink(wallet, *connreq)
			Links[wallet.Name] = link
			if err != nil {
				log.Fatal(err)
			}

			qrCode, _ := qrcode.New(link, qrcode.Medium)
			filename := fmt.Sprintf("../qr/%s_qrcode.png", wallet.Name)
			_ = qrCode.WriteFile(256, filename)
		}

		return c.Send("Choose wallet to connect", menu)
	})
	h := NewHandler(b.wallet, "1000000")
	b.bot.Handle(&btnWallet, h.Wallet)
	b.bot.Handle(&btnTonKeeper, h.TonKeeper)
	b.bot.Handle(&btnTonHub, h.TonHub)
	b.bot.Handle(&btnMyTonWallet, h.MyTonWallet)
	b.bot.Handle("/transaction", h.Transaction)
	b.bot.Handle("/disconnect", h.BtnDisconnect)
}
