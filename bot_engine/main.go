package botengine

import (
	"github.com/judeosbert/bus_tracker_bot/admin"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
)

type botEngine struct {
	admin admin.AdminUtils
	handlers   [](func(telego.Message,state.Saver,admin.AdminUtils) (*telego.SendMessageParams, error))
	outChan    chan telego.SendMessageParams
	msgChan    chan telego.Message
	stateSaver state.Saver
}

// RegisterHandler implements BotEngine.
func (b *botEngine) RegisterHandler(h func(telego.Message,state.Saver,admin.AdminUtils) (*telego.SendMessageParams, error)) {
	b.handlers = append(b.handlers, h)
}

// PostUpdate implements BotEngine.
func (b *botEngine) PostUpdate(update telego.Message) error {
	go func() {
		b.msgChan <- update
	}()
	return nil
}

// Start implements BotEngine.
func (b *botEngine) Start() {
	b.RegisterHandler(HandleAfterPnrSent)
	b.RegisterHandler(HandleAfterProviderSent)
	go func(){
		for update := range b.msgChan {
			i := 0
			for i<len(b.handlers){
				h := b.handlers[i]
				msg,err := h(update,b.stateSaver,b.admin)
				if err != nil{
					i++
					continue
				}
				b.outChan <- *msg
				break
			}
		}
	}()
}

func (b *botEngine) OutChan() <-chan telego.SendMessageParams {
	return b.outChan
}

type BotEngine interface {
	Start()
	PostUpdate(update telego.Message) error
	OutChan() <-chan telego.SendMessageParams
	RegisterHandler(func(telego.Message,state.Saver,admin.AdminUtils) (*telego.SendMessageParams, error))
}

func NewBotEnginer(saver state.Saver,admin admin.AdminUtils) BotEngine {
	return &botEngine{
		admin: admin,
		msgChan:    make(chan telego.Message),
		outChan:    make(chan telego.SendMessageParams),
		stateSaver: saver,
	}
}
