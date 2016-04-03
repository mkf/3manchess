package client

import "github.com/ArchieT/3manchess/multi"
import "github.com/ArchieT/3manchess/server"
import "github.com/ArchieT/3manchess/game"
import "github.com/dghubble/sling"
import "net/http"
import "fmt"

//type WhatServer

//Client is the base url with a pointer to a Sling pointer
type Client struct {
	BaseURL string
	*Service
}

//Service is a Sling pointer
type Service struct {
	sling *sling.Sling
}

//NewService creates a Sling instance
func NewService(httpClient *http.Client, baseURL string) *Service {
	return &Service{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

//NewClient creates a Sling instance, with baseurl in Client struct
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

//SignUp : /api/signup
func (s *Service) SignUp(sp multi.SignUpPost) (*multi.SignUpGive, *http.Response, error) {
	give := new(multi.SignUpGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/signup").BodyJSON(sp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//LogIn : /api/login
func (s *Service) LogIn(lp multi.LoggingIn) (*multi.Authorization, *http.Response, error) {
	give := new(multi.Authorization)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/login").BodyJSON(lp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//BotKey : /api/botkey
func (s *Service) BotKey(bkg multi.BotKeyGetting) (*multi.Authorization, *http.Response, error) {
	give := new(multi.Authorization)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/botkey").BodyJSON(bkg).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//NewBot : /api/newbot
func (s *Service) NewBot(nbp multi.NewBotPost) (*multi.NewBotGive, *http.Response, error) {
	give := new(multi.NewBotGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/newbot").BodyJSON(nbp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//AddGame : /api/addgame
func (s *Service) AddGame(gpp multi.GameplayPost) (*multi.GameplayGive, *http.Response, error) {
	give := new(multi.GameplayGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post("api/addgame").BodyJSON(gpp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//Turn : /api/play/{gameId}
func (s *Service) Turn(gameid int64, turnp multi.TurnPost) (*multi.MoveAndAfterKeys, *http.Response, error) {
	give := new(multi.MoveAndAfterKeys)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Post(fmt.Sprintf("api/play/%d", gameid)).BodyJSON(turnp).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

var colquernms = [3]string{"white", "gray", "black"}

func queraft(p [3]*int64) string {
	for is := 0; is < 3; is++ {
		if p[is] != nil {
			o := fmt.Sprintf("?%s=%d", colquernms[is], *p[is])
			for i := is + 1; i < 3; i++ {
				if p[i] != nil {
					o += fmt.Sprintf("&%s=%d", colquernms[i], *p[i])
				}
			}
			return o
		}
	}
	return ""
}

//After : /api/play/{gameId}/after?white=123&gray=456
func (s *Service) After(gameid int64, filterplayers [3]*int64) (*[]server.MoveFollow, *http.Response, error) {
	give := new([]server.MoveFollow)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/play/%d/after", gameid)+queraft(filterplayers)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//Before : /api/play/{gameId}/before
func (s *Service) Before(gameid int64) (*[]server.MoveFollow, *http.Response, error) {
	give := new([]server.MoveFollow)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/play/%d/before", gameid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//OwnersBots : /api/user/{userId}/bots
func (s *Service) OwnersBots(owner int64) (*[]server.BotFollow, *http.Response, error) {
	give := new([]server.BotFollow)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/user/%d/bots", owner)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//Play : /api/play/{gameId}
func (s *Service) Play(gameid int64) (*server.GameplayData, *http.Response, error) {
	give := new(server.GameplayData)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/play/%d", gameid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//State : /api/state/{stateId}
func (s *Service) State(stateid int64) (*game.State, *http.Response, error) {
	give := new(game.State)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/state/%d", stateid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//VFTPGen : /api/state/{stateId}/vftpgen
func (s *Service) VFTPGen(stateid int64) (*multi.VFTPGenGive, *http.Response, error) {
	give := new(multi.VFTPGenGive)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/state/%d/vftpgen", stateid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//Move : /api/move/{moveId}
func (s *Service) Move(moveid int64) (*server.MoveData, *http.Response, error) {
	give := new(server.MoveData)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/move/%d", moveid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//WhoIsIt : /api/player/{playerId}
func (s *Service) WhoIsIt(playerid int64) (*multi.InfoWhoIsIt, *http.Response, error) {
	give := new(multi.InfoWhoIsIt)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/player/%d", playerid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//UserInfo : /api/user/{userId}
func (s *Service) UserInfo(userid int64) (*server.InfoUser, *http.Response, error) {
	give := new(server.InfoUser)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/user/%d", userid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}

//BotInfo : /api/bot/{botId}
func (s *Service) BotInfo(botid int64) (*server.InfoBot, *http.Response, error) {
	give := new(server.InfoBot)
	ser := new(multi.APIListErr)
	resp, err := s.sling.New().Get(fmt.Sprintf("api/bot/%d", botid)).Receive(give, ser)
	return give, resp, rerr(err, *ser)
}
