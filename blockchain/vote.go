package blockchain

import "encoding/json"

type Voting struct {
	Title string `json:"title"`
}

func NewVoting(title string) Voting {
	return Voting{
		Title: title,
	}
}

func (v Voting) Data() []byte {
	data, _ := json.Marshal(v)
	return data
}

type Vote struct {
	BlockHash string `json:"blockHash"`
	Value     bool   `json:"value"`
}

func NewVote(blockHash string, value bool) Vote {
	return Vote{
		BlockHash: blockHash,
		Value:     value,
	}
}

func (v Vote) Data() []byte {
	data, _ := json.Marshal(v)
	return data
}
