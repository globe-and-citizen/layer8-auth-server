<template>
  <div v-if="active" class="cnc-pay-popup" @click.self="closePopup">
    <div class="cnc-pay-popup__content">
      <button
        class="cnc-pay-popup__close"
        type="button"
        @click="closePopup"
      >
        ×
      </button>
      <div ref="widgetContainer"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {nextTick, onMounted, onUnmounted, ref, watch} from 'vue'
import {config} from '@/config'

const props = defineProps<{
  active: boolean
  factureId: string
  amount: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

type CncPayApi = {
  setFactureId(factureId: string): void
  setAmount(amount: string): void
  setOnStatus(callback: (status: string) => void): void
  show(target: string | HTMLElement): void
}

const cncPay = () =>
  (window as unknown as { CncPay?: CncPayApi }).CncPay

const widgetReady = ref(false)
const widgetContainer = ref<HTMLElement | null>(null)

let widgetScript: HTMLScriptElement | undefined

function loadWidgetScript(): HTMLScriptElement {
  const existing = document.getElementById('cnc-pay-widget')

  if (existing instanceof HTMLScriptElement) return existing

  const script = document.createElement('script')
  script.id = 'cnc-pay-widget'
  script.src = 'https://payment-gate-widget-staging.up.railway.app/widget.js'
  script.setAttribute('data-bank', config.CONTRACT_ADDRESS)
  script.setAttribute('data-token', 'USDC')
  script.async = true

  document.body.appendChild(script)
  return script
}

const onWidgetLoad = () => {
  widgetReady.value = !!cncPay()
}

onMounted(() => {
  if (cncPay()) {
    onWidgetLoad()
    return
  }

  widgetScript = loadWidgetScript()
  widgetScript.addEventListener('load', onWidgetLoad)

  // Handle the case where the script loaded before
  // this component registered its listener.
  if (cncPay()) onWidgetLoad()
})

onUnmounted(() => {
  widgetScript?.removeEventListener('load', onWidgetLoad)
})

function closePopup() {
  emit('close')
}

async function openCheckout() {
  if (!props.active || !widgetReady.value) return

  await nextTick()

  const widget = cncPay()
  if (!widget || !widgetContainer.value || !props.active) return

  widget.setFactureId(props.factureId)
  widget.setAmount(props.amount)

  widget.setOnStatus((status) => {
    console.log('payment status:', status)

    if (status === 'success' || status === 'completed') {
      emit('success')
      closePopup()
    }
  })

  widget.show(widgetContainer.value)
}

watch(
  [() => props.active, widgetReady],
  ([active, ready]) => {
    if (active && ready) {
      void openCheckout()
    }
  },
  {immediate: true}
)
</script>

<style scoped>
.cnc-pay-popup {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  background: rgb(0 0 0 / 60%);
}

.cnc-pay-popup__content {
  position: relative;
  width: min(100%, 480px);
  max-height: 90vh;
  overflow: auto;
  padding: 2.5rem 1rem 1rem;
  border-radius: 12px;
  background: #fff;
}

.cnc-pay-popup__close {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  border: 0;
  background: transparent;
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
}
</style>
