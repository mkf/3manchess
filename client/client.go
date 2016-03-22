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

func rerr(httpError error, ale multi.APIListErr) error {
	if httpError != nil {
		return httpError
	}
	return ale.ToErr()
}

func (s *Service) SignUp(sp multi.SignUpPost) (*multi.SignUpGive, *http.Response, error) {
	give := new(multi.SignUpGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/signup").BodyJSON(sp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) LogIn(lp multi.LoggingIn) (*multi.Authorization, *http.Response, error) {
	give := new(multi.Authorization)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/login").BodyJSON(lp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) BotKey(bkg multi.BotKeyGetting) (*multi.Authorization, *http.Response, error) {
	give := new(multi.Authorization)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/botkey").BodyJSON(bkg).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) NewBot(nbp multi.NewBotPost) (*multi.NewBotGive, *http.Response, error) {
	give := new(multi.NewBotGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/newbot").BodyJSON(nbp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) AddGame(gpp multi.GameplayPost) (*multi.GameplayGive, *http.Response, error) {
	give := new(multi.GameplayGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/addgame").BodyJSON(gpp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) Turn(gameid int64, turnp multi.TurnPost) (*multi.MoveAndAfterKeys, *http.Response, error) {
	give := new(multi.MoveAndAfterKeys)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post(fmt.Sprintf("api/play/%d", gameid)).BodyJSON(turnp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) Play(gameid int64) (*server.GameplayData, *http.Response, error) {
	give := new(server.GameplayData)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/play/%d", gameid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) State(stateid int64) (*game.StateData, *http.Response, error) {
	give := new(game.StateData)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/state/%d", stateid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) VFTPGen(stateid int64) (*multi.VFTPGenGive, *http.Response, error) {
	give := new(multi.VFTPGenGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/state/%d/vftpgen", stateid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) Move(moveid int64) (*server.MoveData, *http.Response, error) {
	give := new(server.MoveData)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/move/%d", moveid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) WhoIsIt(playerid int64) (*multi.InfoWhoIsIt, *http.Response, error) {
	give := new(multi.InfoWhoIsIt)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/player/%d", playerid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) UserInfo(userid int64) (*multi.InfoUser, *http.Response, error) {
	give := new(multi.InfoUser)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/user/%d", userid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

func (s *Service) BotInfo(botid int64) (*multi.InfoBot, *http.Response, error) {
	give := new(multi.InfoBot)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/bot/%d", botid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}
