<script setup lang="ts">
import api from '@/api'
import DarkModeToggle from '@/components/DarkModeToggle.vue'
import FullLogo from '@/components/FullLogo.vue'
import { onMounted, ref } from 'vue'

const props = defineProps<{
  id?: string
}>()

const createLobby = () => {
  if (name.value !== '') {
    localStorage.setItem('nickname', name.value)
  }

  const lobbyName = name.value || previousName.value

  if (props.id) {
    api.lobby.join(props.id, lobbyName)
  } else {
    api.lobby.create(lobbyName)
  }
}

const name = ref('')

const previousName = ref('')

onMounted(() => {
  const nick = localStorage.getItem('nickname')
  if (nick) {
    previousName.value = nick
  }
})
</script>

<template>
  <div class="box container b-primary wrapper">
    <div class="logo-header mb-6">
      <div class="outer"></div>
      <FullLogo :width="500" />
      <div class="outer">
        <DarkModeToggle />
      </div>
    </div>
    <div class="name">
      <div class="box mb-3 b-primary">
        <h2 class="title is-4 has-text-primary">Choose your Nickname...</h2>
        <form>
          <div class="field">
            <p class="control">
              <input
                class="input"
                type="text"
                :placeholder="previousName || `Enter Nickname...`"
                v-model="name"
              />
            </p>
          </div>
          <hr />
          <div class="field is-grouped is-grouped-centered">
            <div class="control">
              <button class="button is-primary" @click.prevent="createLobby">
                {{ id ? 'Join Lobby' : 'Create Lobby' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
