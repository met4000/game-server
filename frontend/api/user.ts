import { useUserStore, userSchema } from '@/stores/user'
import { parseData, sendMessage } from './main'
import { OPacketType } from './packetTypes'
import { useSocketStore } from '@/stores/socket'
import { useLobbyStore } from '@/stores/lobby'
import router from '@/router'

const idLSKey = 'id_token'

export const handleUserInit = (data: unknown) => {
  const parsedData = parseData(data, userSchema)

  localStorage.setItem(idLSKey, parsedData.token ?? '')

  const user = useUserStore()
  user.data = parsedData

  const lobbyStore = useLobbyStore()
  const socket = useSocketStore()
  socket.userReady = true

  if (user.data.lobbies.length > 0) {
    lobbyStore.getInfo()
  } else {
    router.replace({ name: 'home', query: router.currentRoute.value.query })
  }
}

export const handleUserChange = (data: unknown) => {
  const parsedData = parseData(data, userSchema)

  const user = useUserStore()
  user.data = parsedData
}

const getLocalToken = () => {
  return localStorage.getItem(idLSKey)
}

const sendNameChange = (name: string) => {
  sendMessage({
    type: OPacketType.UserNameChange,
    data: {
      name,
    },
  })
}

export default { getLocalToken, sendNameChange }
