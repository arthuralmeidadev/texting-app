package providers

import "texting-app/internal/pkg/models"

type messageStore interface {
	StoreMessage(sen, cont string, chatId uint, repTo int) error
	UpdateMessageStatus(msgId uint, stat string) error
	GetMessages(chatId, offset uint) ([]*models.Message, error)
	GetUnseenMessagesTotal(usrn, rec string) (uint, error)
	DeleteMessage(msgId uint) error
}

type messagesProvider struct {
	store messageStore
}

func (p *messagesProvider) StoreMessage(sen, cont string, chatId uint, repTo int) error {
	return p.store.StoreMessage(sen, cont, chatId, repTo)
}

func (p *messagesProvider) GetMessages(chatId, offset uint) ([]*models.Message, error) {
	return p.store.GetMessages(chatId, offset)
}

func NewMessageProvider(s messageStore) *messagesProvider {
	return &messagesProvider{
		store: s,
	}
}
