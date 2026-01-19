package main

import "math"

type Share struct {
	UserID int     `json:"userId"`
	Amount float64 `json:"share"`
}

type SplitStrategy interface {
	Calculate(amount float64, participants []Participant) []Share
}

type EqualSplit struct{}
func (s EqualSplit) Calculate(amount float64, participants []Participant) []Share {
	share := math.Floor((amount/float64(len(participants)))*100) / 100
	shares := make([]Share, len(participants))
	for i, p := range participants {
		shares[i] = Share{UserID: p.UserID, Amount: share}
	}
	return shares
}

type ExactSplit struct{}
func (s ExactSplit) Calculate(amount float64, participants []Participant) []Share {
	shares := make([]Share, len(participants))
	for i, p := range participants {
		shares[i] = Share{UserID: p.UserID, Amount: p.Amount}
	}
	return shares
}