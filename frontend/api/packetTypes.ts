export enum IPacketType {
  Chat = 'server_chat',
  Error = 'server_error',
  LobbyChange = 'server_lobby_change',
  LobbyUserChange = 'server_lobby_change_user',
  UserInit = 'server_user_init',
  UserChange = 'server_user_change',
  GameEvent = 'server_game_event',
  GameState = 'server_game_state',
  GameAction = 'server_game_action',
  GameError = 'server_error_game',
}

export enum OPacketType {
  Chat = 'client_chat',
  CreateLobby = 'client_lobby_create',
  JoinLobby = 'client_lobby_join',
  LeaveLobby = 'client_lobby_leave',
  LobbyUsers = 'client_lobby_users',
  LobbyInfo = 'client_lobby_info',
  LobbyChange = 'client_lobby_change',
  LobbyOptions = 'client_lobby_options',
  UserNameChange = 'client_user_name_change',
  GameAction = 'client_game_action',
  GameNew = 'client_game_new',
  GameReady = 'client_game_ready',
}
