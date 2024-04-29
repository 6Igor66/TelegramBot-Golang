package telegram

import (
	"bot/internal/models"
	"context"
	"log"
	"time"

	"github.com/cameo-engineering/tonconnect"
	"golang.org/x/exp/maps"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	wallet string
	amount string
}

func NewHandler(wallet, amount string) *Handler {
	return &Handler{
		wallet: wallet,
		amount: amount,
	}
}

var UsersState = make(map[int64]*models.UserState)

var Links = make(map[string]string)

func (h *Handler) Wallet(c tele.Context) error {
	link := Links["Wallet"]
	inlineBtn := tele.InlineButton{
		Unique: "Wallet",
		URL:    link,
		Text:   "Connect",
	}

	inlineKeys := [][]tele.InlineButton{
		[]tele.InlineButton{inlineBtn},
	}

	path := "../qr/Wallet_qrcode.png"
	photo := tele.FromDisk(path)

	return c.Send(&tele.Photo{File: photo},
		"Press the button to connect:",
		&tele.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}

func (h *Handler) TonKeeper(c tele.Context) error {
	link := Links["Tonkeeper"]
	inlineBtn := tele.InlineButton{
		Unique: "Tonkeeper",
		URL:    link,
		Text:   "Connect",
	}

	inlineKeys := [][]tele.InlineButton{
		[]tele.InlineButton{inlineBtn},
	}

	path := "../qr/Tonkeeper_qrcode.png"
	photo := tele.FromDisk(path)

	return c.Send(&tele.Photo{File: photo},
		"Press the button to connect:",
		&tele.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}

func (h *Handler) TonHub(c tele.Context) error {
	link := Links["Tonhub"]
	inlineBtn := tele.InlineButton{
		Unique: "Tonhub",
		URL:    link,
		Text:   "Connect",
	}

	inlineKeys := [][]tele.InlineButton{
		[]tele.InlineButton{inlineBtn},
	}

	path := "../qr/Tonhub_qrcode.png"
	photo := tele.FromDisk(path)

	return c.Send(&tele.Photo{File: photo},
		"Press the button to connect:",
		&tele.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}

func (h *Handler) MyTonWallet(c tele.Context) error {
	link := Links["MyTonWallet"]
	inlineBtn := tele.InlineButton{
		Unique: "MyTonWallet",
		URL:    link,
		Text:   "Connect",
	}

	inlineKeys := [][]tele.InlineButton{
		[]tele.InlineButton{inlineBtn},
	}

	path := "../qr/MyTonWallet_qrcode.png"
	photo := tele.FromDisk(path)

	return c.Send(&tele.Photo{File: photo},
		"Press the button to connect:",
		&tele.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}

func (h *Handler) Transaction(c tele.Context) error {
	if !(UsersState[c.Chat().ID].Connected) {
		ctx := context.Background()
		UsersState[c.Chat().ID].Ctx = ctx
		s := UsersState[c.Chat().ID].S

		_, err := s.Connect(UsersState[c.Chat().ID].Ctx, (maps.Values(tonconnect.Wallets))...)
		if err != nil {
			log.Fatal(err)
		}
		UsersState[c.Chat().ID].Connected = true
	}

	ctx := UsersState[c.Chat().ID].Ctx
	s := UsersState[c.Chat().ID].S
	msg, err := tonconnect.NewMessage(h.wallet, h.amount)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := tonconnect.NewTransaction(
		tonconnect.WithTimeout(10*time.Minute),
		tonconnect.WithMainnet(),
		tonconnect.WithMessage(*msg),
	)
	if err != nil {
		log.Fatal(err)
	}
	boc, err := s.SendTransaction(ctx, *tx)
	if err != nil {
		log.Println(err)
		return c.Send("you rejected the transaction")
	} else {
		log.Printf("Bag of Cells: %x", boc)
		return c.Send("success")
	}
}

func (h *Handler) BtnDisconnect(c tele.Context) error {
	_, ok := UsersState[c.Chat().ID]
	if !ok || !UsersState[c.Chat().ID].Connected {
		c.Send("you are not connected")
		return nil
	}
	ctx := UsersState[c.Chat().ID].Ctx
	s := UsersState[c.Chat().ID].S
	delete(UsersState, c.Chat().ID)
	c.Send("you've successfully disconnected")
	return s.Disconnect(ctx)
}
