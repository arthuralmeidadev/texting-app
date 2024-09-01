package providers

import "texting-app/internal/pkg/models"

type chatStore interface {
	CreateChat(name string, usrs []string) (uint, error)
	GetChats(usrn string, offset uint) ([]*models.Chat, error)
    GetChatMembers(chatId uint) ([]string, error)
	DeleteChat(chatId uint) error
}

type ChatProvider struct {
	store chatStore
}


func (p *ChatProvider) GetChatMembers(charId uint) ([]string, error) {
    return p.store.GetChatMembers(charId)
}

func NewChatProvider(s chatStore) *ChatProvider {
	return &ChatProvider{
		store: s,
	}
}
