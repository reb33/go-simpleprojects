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

func (s *Store) Upsert(phone, sessionId, code string) *Data {
	data := s.GetByPhone(phone)
	if data == nil {
		data := Data{
			phone:     phone,
			sessionId: sessionId,
			code:      code,
		}
		s.data = append(s.data, data)
		return &data
	}
	data.sessionId = sessionId
	data.code = code
	return data
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