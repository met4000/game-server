<script setup lang="ts">
import ModalWrap from '@/components/ModalWrap.vue'
import type { RoundInfo } from '@/game/liarsdice'
import { useLobbyStore } from '@/stores/lobby'
import DiceHand from './DiceHand.vue'
import { computed, ref, watch } from 'vue'
import DiceHandValues from './DiceHandValues.vue'
import DiceHandTabs from './DiceHandTabs.vue'
import BidRender from './BidRender.vue'

const props = defineProps<{
  show: boolean
  data?: RoundInfo
}>()

const trueBidAmount = computed(() => {
  return props.data?.highestBid.split(',').map((v) => parseInt(v))[1] ?? 0
})

const lobby = useLobbyStore()

const getName = (id: string) => lobby.users[id]?.name ?? ''

const emit = defineEmits<{
  close: []
}>()

const handleClose = () => {
  emit('close')
}

const calculatedAmount = (f: number): number => {
  if (!props.data) {
    return 0
  }

  return Object.values(props.data?.hands ?? {})
    .flat()
    .reduce((a, v) => (v === f || v === 1 ? a + 1 : a), 0)
}

const faceSelected = ref(0)

watch(
  () => props.data?.highestBid,
  (nv) => {
    if (nv) {
      const r = nv.split(',').map((v) => parseInt(v))
      faceSelected.value = r[1]
    }
  },
)

const tabIndex = ref(0)

const handleChangeTab = (n: number) => {
  tabIndex.value = n
}
</script>

<template>
  <ModalWrap title="Call Result" :shown="show" @close="handleClose">
    <template #body>
      <div v-if="data">
        <div>
          <p>
            {{ getName(data.callUser) }} called out {{ getName(data.lastBid) }} for the bid
            <BidRender :bid="data.highestBid" />
          </p>
          <p>There were: <BidRender :bid="[calculatedAmount(trueBidAmount), trueBidAmount]" /></p>
          <p>{{ getName(data.diceLost) }} lost a dice!</p>
        </div>
        <div>
          <div class="is-flex is-justify-content-space-between is-align-items-center">
            <h4 class="title is-4 mb-0">Hands</h4>
            <div class="inline-select-label">
              <label class="label">Highlighted Face:</label>
              <div class="select">
                <select v-model.number="faceSelected">
                  <option>2</option>
                  <option>3</option>
                  <option>4</option>
                  <option>5</option>
                  <option>6</option>
                </select>
              </div>
            </div>
          </div>
          <div>There were: <BidRender :bid="[calculatedAmount(faceSelected), faceSelected]" /></div>
          <DiceHandTabs @change-tab="handleChangeTab" />
          <div :class="{ 'value-icon': tabIndex === 1 }">
            <div v-for="(h, id) in data.hands" :key="id" class="mb-2">
              <h6 class="title is-6 mb-2">{{ getName(id) }}</h6>
              <DiceHandValues :true-value="h" :highlight="faceSelected" :tab-index="tabIndex" />
            </div>
          </div>
        </div>
      </div>
      <div v-else>Loading...</div>
    </template>
    <template #footer>
      <div class="buttons">
        <button class="button" @click="handleClose">Continue</button>
      </div>
    </template>
  </ModalWrap>
</template>

<style>
.value-icon {
  display: flex;
  gap: 2rem;
}

.inline-select-label {
  display: flex;
  gap: 1rem;
  align-items: center;
}
</style>
