package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/game"
	"github.com/google/uuid"
)

// All outgoing lobby data
type m_Lobby struct {
	Id         uuid.UUID   `json:"id"`
	MaxPlayers *int        `json:"maxPlayers"`
	MinPlayers *int        `json:"minPlayers"`
	Name       *string     `json:"name"`
	GameType   *string     `json:"gameType"`
	InGame     *bool       `json:"inGame"`
	Options    interface{} `json:"gameOptions"`
}

// incoming client lobby change
type m_LobbyChange struct {
	Id         uuid.UUID
	MaxPlayers int
	MinPlayers int
	Name       string
	GameType   string
}

type m_LobbyJoin struct {
	Id   uuid.UUID
	Name string
}

const pt_OLobbyUserChange OPacketType = "server_lobby_change_user"

const pt_OLobbyChange OPacketType = "server_lobby_change"

// Create Lobby Event
const pt_ICreateLobby IPacketType = "client_lobby_create"

func createLobby(i HandlerInput) error {
	data, err := hp[m_LobbyJoin](i.Packet)
	if err != nil {
		return err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = db.CreateLobby(id, i.UserId)
	if err != nil {
		return err
	}

	if data.Name != "" {
		err = db.SetUserName(i.UserId, data.Name)
		if err != nil {
			return err
		}
	}

	sendLobbyDataSingle(i.C, id, i.UserId)

	return nil
}

// Join lobby event
const pt_IJoinLobby IPacketType = "client_lobby_join"

func joinLobby(i HandlerInput) error {
	data, err := hp[m_LobbyJoin](i.Packet)
	if err != nil {
		return err
	}

	exists, err := db.LobbyExists(data.Id)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("lobby does not exist")
	}

	if data.Name != "" {
		err = db.SetUserName(i.UserId, data.Name)
		if err != nil {
			return err
		}
	}

	err = db.SaveLobbyUser(data.Id, i.UserId)
	if err != nil {
		return err
	}

	sendLobbyUserChange(i.C, data.Id)
	sendLobbyDataSingle(i.C, data.Id, i.UserId)

	return nil
}

// User request all lobby users
const pt_ILobbyUsers IPacketType = "client_lobby_users"

func lobbyUsers(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	_, packet := makeLobbyUsersMessage(*lobbyId)
	Send(i.C, i.UserId, packet)

	return nil
}

func makeLobbyUsersMessage(lobbyId uuid.UUID) ([]uuid.UUID, *Packet[db.LobbyUserList]) {
	users, err := db.GetLobbyUsers(lobbyId)
	if err != nil {
		return []uuid.UUID{}, nil
	}
	keys := make([]uuid.UUID, 0, len(users))

	for u := range users {
		keys = append(keys, u)
	}

	newMessage := mp(pt_OLobbyUserChange, users)
	return keys, &newMessage
}

func sendLobbyUserChange(c Client, lobbyId uuid.UUID) {
	ids, packet := makeLobbyUsersMessage(lobbyId)
	SendMany(c, ids, packet)
}

const pt_ILobbyInfo IPacketType = "client_lobby_info"

func lobbyInfo(i HandlerInput) error {
	lobbyIds, err := db.GetUserLobbies(i.UserId)
	if err != nil {
		return err
	}

	for _, lobbyId := range lobbyIds {
		sendLobbyDataSingle(i.C, lobbyId, i.UserId)
	}

	return nil
}

func lobbyInfoPacket(lobbyId uuid.UUID) (*Packet[m_Lobby], error) {
	var info m_Lobby

	info.Id = lobbyId

	properties := []string{"name", "maxPlayers", "minPlayers", "gameType", "inGame", "options"}
	m, err := db.GetLobbyProperties(lobbyId, properties)
	if err != nil {
		return nil, err
	}

	err = decode(m, &info)
	if err != nil {
		return nil, err
	}

	packet := mp(pt_OLobbyChange, info)

	return &packet, nil
}

func sendLobbyDataSingle(c Client, lobbyId uuid.UUID, userId uuid.UUID) {
	lobbyPacket, err := lobbyInfoPacket(lobbyId)
	if err != nil {
		log.Println("failed to make lobby info packet")
	}

	Send(c, userId, lobbyPacket)
}

func sendLobbyData(c Client, lobbyId uuid.UUID) {
	lobbyPacket, err := lobbyInfoPacket(lobbyId)
	if err != nil {
		log.Println("failed to make lobby info packet")
	}

	ids, err := db.GetLobbyUserIds(lobbyId)
	if err != nil {
		log.Println("failed to get lobby users")
	}

	SendMany(c, ids, lobbyPacket)
}

func makeChangesAllowed(lobbyId uuid.UUID, userId uuid.UUID) error {
	isHost, err := db.IsUserLobbyHost(lobbyId, userId)
	if err != nil {
		return err
	}
	if !isHost {
		return fmt.Errorf("not lobby host")
	}

	return nil
}

const pt_ILobbyChange = "client_lobby_change"

func lobbyChange(i HandlerInput) error {
	data, err := hp[m_LobbyChange](i.Packet)
	if err != nil {
		return err
	}

	err = makeChangesAllowed(data.Id, i.UserId)
	if err != nil {
		return err
	}

	oldGameType, err := db.GetLobbyProperty[string](data.Id, "gameType")
	if err != nil {
		return err
	}
	if data.GameType != oldGameType {
		err = db.SetLobbyProperty(data.Id, "gameType", data.GameType)
		if err != nil {
			return err
		}
		options, err := game.GetDefaultOptions(data.GameType)
		if err != nil {
			return err
		}
		err = db.SetLobbyProperty(data.Id, "options", options)
		if err != nil {
			return err
		}
	}

	sendLobbyData(i.C, data.Id)

	return nil
}

const pt_ILobbyOptions = "client_lobby_options"

type m_LobbyOptions struct {
	Id   uuid.UUID
	Data json.RawMessage
}

func lobbyOptions(i HandlerInput) error {
	data, err := hp[m_LobbyOptions](i.Packet)
	if err != nil {
		return err
	}

	err = makeChangesAllowed(data.Id, i.UserId)
	if err != nil {
		return err
	}

	err = db.SetLobbyProperty(data.Id, "options", string(data.Data))
	if err != nil {
		return err
	}

	sendLobbyData(i.C, data.Id)
	return nil
}

// Leave lobby event
const pt_ILeaveLobby IPacketType = "client_lobby_leave"

func leaveLobby(i HandlerInput) error {
	lobbyId, err := hp[uuid.UUID](i.Packet)
	if err != nil {
		return err
	}

	ig, err := db.LobbyInGame(*lobbyId)
	if err != nil {
		return err
	}
	if ig {
		g := &gc{
			GameId: *lobbyId,
			C:      i.C,
		}
		err = game.HandleLeave(g, *lobbyId, i.UserId)
		if err != nil {
			return err
		}
	}

	err = db.RemoveLobbyUser(*lobbyId, i.UserId)
	if err != nil {
		return err
	}

	sendLobbyUserChange(i.C, *lobbyId)
	sendLobbyDataSingle(i.C, uuid.Nil, i.UserId)

	return nil
}
