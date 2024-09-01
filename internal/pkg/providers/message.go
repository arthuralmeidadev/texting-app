package providers

import "texting-app/internal/pkg/models"

type messageStore interface {
	StoreMessage(sen, cont string, chatId uint, repTo int) error
	UpdateMessageStatus(msgId uint, stat string) error
	GetMessages(chatId, offset uint) ([]*models.Message, error)
	GetUnseenMessagesTotal(usrn, rec string) (uint, error)
	DeleteMessage(msgId uint) error
}

type MessagesProvider struct {
	store messageStore
}

func (p *MessagesProvider) StoreMessage(sen, cont string, chatId uint, repTo int) error {
	return p.store.StoreMessage(sen, cont, chatId, repTo)
}

func (p *MessagesProvider) GetMessages(chatId, offset uint) ([]*models.Message, error) {
	return p.store.GetMessages(chatId, offset)
}

func NewMessageProvider(s messageStore) *MessagesProvider {
	return &MessagesProvider{
		store: s,
	}
}
