package client

import "github.com/ArchieT/3manchess/multi"
import "github.com/ArchieT/3manchess/server"
import "github.com/ArchieT/3manchess/game"
import "github.com/dghubble/sling"
import "net/http"
import "fmt"

//type WhatServer

type Client struct {
	BaseURL string
	*Service
}

type Service struct {
	sling *sling.Sling
}

func NewService(httpClient *http.Client, baseURL string) *Service {
	return &Service{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		Service: NewService(httpClient, baseURL),
		BaseURL: baseURL,
	}
}

func (s *Service) SignUp(sp multi.SignUpPost) (*multi.SignUpGive, *http.Response, error) {
	give := new(multi.SignUpGive)
	ser := new(multi.OurError)
	resp, err := s.sling.New().Post("api/signup").BodyJSON(sp).Receive(give, ser)
	if err == nil {
		err = ser
	}
	return give, resp, err
}

func (s *Service) LogIn(lp multi.LoggingIn) (*multi.Authorization, *http.Response, error) {}

func (s *Service) BotKey(bkg multi.BotKeyGetting) (*multi.Authorization, *http.Response, error) {}

func (s *Service) NewBot(nbp multi.NewBotPost) (*multi.NewBotGive, *http.Response, error) {}

func (s *Service) AddGame(gpp multi.GameplayPost) (*multi.GameplayGive, *http.Response, error) {}

func (s *Service) Turn(gameid int64, turnp multi.TurnPost) (*multi.MoveAndAfterKeys, *http.Response, error) {
}

func (s *Service) Play(gameid int64) (*server.GameplayData, *http.Response, error) {}

func (s *Service) State(stateid int64) (*game.StateData, *http.Response, error) {}

func (s *Service) Move(moveid int64) (*server.MoveData, *http.Response, error) {}

func (s *Service) WhoIsIt(playerid int64) (*multi.InfoWhoIsIt, *http.Response, error) {}

func (s *Service) UserInfo(userid int64) (*multi.InfoUser, *http.Response, error) {}

func (s *Service) BotInfo(botid int64) (*multi.InfoBot, *http.Response, error) {}
