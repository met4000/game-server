import { handleUserChange, handleUserInit } from './user'
import error from './error'
import { z } from 'zod'
import { IPacketType } from './packetTypes'
import { handleLobbyChange, handleLobbyUserChange } from './lobby'
import { useErrorStore } from '@/stores/error'
import { handleGameAction, handleGameEvent, handleGameState } from './game'
import { notImplemented } from './main'

const routeMap: Record<IPacketType, (data: unknown) => void> = {
  [IPacketType.UserInit]: handleUserInit,
  [IPacketType.UserChange]: handleUserChange,
  [IPacketType.Chat]: notImplemented,
  [IPacketType.LobbyChange]: handleLobbyChange,
  [IPacketType.LobbyUserChange]: handleLobbyUserChange,
  [IPacketType.Error]: error.handle,
  [IPacketType.GameState]: handleGameState,
  [IPacketType.GameAction]: handleGameAction,
  [IPacketType.GameEvent]: handleGameEvent,
  [IPacketType.GameError]: error.handle,
}

const incomingSchema = z.object({
  type: z.nativeEnum(IPacketType),
  data: z.unknown(),
})

export const handleIncomingMessage = (message: string) => {
  const msg = JSON.parse(message)

  const result = incomingSchema.safeParse(msg)

  const es = useErrorStore()

  if (!result.success) {
    // TODO: better error
    es.add({
      type: 'warning',
      message: 'Unknown data received from server',
    })

    console.log(result.error.format())
    return
  }

  routeMap[result.data.type](result.data.data)
}
