package auth

type Data struct {
	phone     string
	sessionId string
	code      string
}

type Store struct {
	data []Data
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Upsert(phone, sessionId, code string) Data {
	for i  := range s.data {
		if s.data[i].phone == phone {
			s.data[i].sessionId = sessionId
			s.data[i].code = code
			return s.data[i]
		}
	}
	
	newData := Data{
		phone:     phone,
		sessionId: sessionId,
		code:      code,
	}
	s.data = append(s.data, newData)
	return newData

}

func (s *Store) GetByPhone(phone string) (*Data) {
	for _, d := range s.data {
		if d.phone == phone {
			return &d
		}
	}
	return nil
}

func (s *Store) GetBySessionId(sessionId string) (*Data) {
	for _, d := range s.data {
		if d.sessionId == sessionId {
			return &d
		}
	}
	return nil
}