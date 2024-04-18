package telegram

type Update struct {
	Message Message `json:"message"`
}

type Message struct {
	MessageID int  `json:"message_id"`
	From      From   `json:"from"`      
	Chat      Chat   `json:"chat"`      
	Date      int  `json:"date"`      
	Text      string `json:"text"`      
}

type Chat struct {
	ID        int  `json:"id"`        
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"` 
	Username  string `json:"username"`  
	Type      string `json:"type"`      
}

type From struct {
	ID           int  `json:"id"`           
	IsBot        bool   `json:"is_bot"`       
	FirstName    string `json:"first_name"`   
	LastName     string `json:"last_name"`    
	Username     string `json:"username"`     
	LanguageCode string `json:"language_code"`
}
