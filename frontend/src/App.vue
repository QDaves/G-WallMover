<template>
  <div class="box">
    <h1>G-Wall Mover</h1>
    <div class="mode-toggle">
      <button @click="toggleMode" :class="{ active: !isMultiMode }">Single</button>
      <button @click="toggleMode" :class="{ active: isMultiMode }">Multi</button>
    </div>
    <div class="content-area">
      <div class="controls-section">
        <div v-if="!isMultiMode">
          <div class="row">
            <label for="l1">L1:</label>
            <div class="numpad">
              <button @click="tweak('l1', -step)" class="btn left">-</button>
              <span @wheel="spin($event, 'l1')">{{ l1 }}</span>
              <button @click="tweak('l1', step)" class="btn right">+</button>
            </div>
          </div>
          <div class="row">
            <label for="l2">L2:</label>
            <div class="numpad">
              <button @click="tweak('l2', -step)" class="btn left">-</button>
              <span @wheel="spin($event, 'l2')">{{ l2 }}</span>
              <button @click="tweak('l2', step)" class="btn right">+</button>
            </div>
          </div>
          <div class="row" style="justify-content: flex-start;">
  <button @click="toggleAdvancedOptions" 
          :class="{ active: advancedOptionsEnabled }" 
          class="advanced-options-btn">
    {{ advancedOptionsEnabled ? 'Disable' : 'Enable' }} Advanced Options
  </button>
</div>
          <div class="row">
            <label for="w1">W1:</label>
            <div class="numpad">
              <button @click="tweak('w1', -1)" class="btn left" :disabled="!advancedOptionsEnabled">-</button>
              <span @wheel="advancedOptionsEnabled && spin($event, 'w1')">{{ w1 }}</span>
              <button @click="tweak('w1', 1)" class="btn right" :disabled="!advancedOptionsEnabled">+</button>
            </div>
          </div>
          <div class="row">
            <label for="w2">W2:</label>
            <div class="numpad">
              <button @click="tweak('w2', -1)" class="btn left" :disabled="!advancedOptionsEnabled">-</button>
              <span @wheel="advancedOptionsEnabled && spin($event, 'w2')">{{ w2 }}</span>
              <button @click="tweak('w2', 1)" class="btn right" :disabled="!advancedOptionsEnabled">+</button>
            </div>
          </div>
        </div>
        <div v-else>
          <div class="row">
            <label for="mainL1">L1:</label>
            <div class="numpad">
              <button @click="tweakMain('L1', -step)" class="btn left" :disabled="capturedItemsCount === 0">-</button>
              <span @wheel="capturedItemsCount > 0 && spinMain($event, 'L1')">{{ mainLocation.L1 }}</span>
              <button @click="tweakMain('L1', step)" class="btn right" :disabled="capturedItemsCount === 0">+</button>
            </div>
          </div>
          <div class="row">
            <label for="mainL2">L2:</label>
            <div class="numpad">
              <button @click="tweakMain('L2', -step)" class="btn left" :disabled="capturedItemsCount === 0">-</button>
              <span @wheel="capturedItemsCount > 0 && spinMain($event, 'L2')">{{ mainLocation.L2 }}</span>
              <button @click="tweakMain('L2', step)" class="btn right" :disabled="capturedItemsCount === 0">+</button>
            </div>
          </div>
          <div class="row">
            <div class="numpad">
              <label for="l1Diff">V1:</label>
              <button @click="tweakDiff('L1', -step)" class="btn left" :disabled="capturedItemsCount === 0">-</button>
              <span @wheel="capturedItemsCount > 0 && spinDiff($event, 'L1')">{{ diffLocation.L1 }}</span>
              <button @click="tweakDiff('L1', step)" class="btn right" :disabled="capturedItemsCount === 0">+</button>
            </div>
            <div class="numpad">
              <label for="l2Diff">V2:</label>
              <button @click="tweakDiff('L2', -step)" class="btn left" :disabled="capturedItemsCount === 0">-</button>
              <span @wheel="capturedItemsCount > 0 && spinDiff($event, 'L2')">{{ diffLocation.L2 }}</span>
              <button @click="tweakDiff('L2', step)" class="btn right" :disabled="capturedItemsCount === 0">+</button>
            </div>
          </div>
          <div class="captured-items">Captured Items: {{ capturedItemsCount }}</div>
        </div>
      </div>
      <div class="action-buttons">
        <button v-if="isMultiMode" @click="toggleListening" :class="{ active: isListening }">
          {{ isListening ? 'Stop' : 'Capture' }}
        </button>
        <button v-if="isMultiMode" @click="clearCapturedItems" class="clear-btn">Clear</button>
        <button v-if="isMultiMode" @click="moveAllItems" :disabled="capturedItemsCount === 0">Move All</button>
      </div>
      <div class="step-size">
        <label for="step">Step size:</label>
        <select id="step" v-model="step">
          <option v-for="s in steps" :key="s" :value="s">{{ s }}</option>
        </select>
      </div>
      <div class="tip">Scroll above values or use buttons</div>
      <div id="log" ref="logbox">
        <div v-for="(msg, index) in log" :key="index">{{ msg }}</div>
      </div>
    </div>
  </div>
</template>

<script>
import { UpdatePosition, GetPosition, MoveItem, ToggleListening, ClearCapturedItems, UpdateMainLocation, UpdateDiffLocation, MoveAllItems, ToggleMode } from '../wailsjs/go/main/App'

export default {
  data() {
    return {
      l1: 0,
      l2: 0,
      w1: 0,
      w2: 0,
      step: 5,
      steps: [1, 5, 10, 15, 20, 30, 40, 50],
      log: [],
      timer: null,
      isMultiMode: false,
      isListening: false,
      capturedItemsCount: 0,
      mainLocation: {
        L1: 0,
        L2: 0,
        W1: 0,
        W2: 0
      },
      diffLocation: {
        L1: 0,
        L2: 0
      },
      advancedOptionsEnabled: false
    }
  },
  methods: {
    tweak(field, delta) {
      if ((field === 'w1' || field === 'w2') && !this.advancedOptionsEnabled) return;
      this[field] += delta
      this.sync()
    },
    tweakMain(field, delta) {
      if (this.capturedItemsCount === 0) return;
      this.mainLocation[field] += delta
      this.syncMain()
    },
    tweakDiff(field, delta) {
      if (this.capturedItemsCount === 0) return;
      this.diffLocation[field] += delta
      this.syncDiff()
    },
    spin(e, field) {
      if ((field === 'w1' || field === 'w2') && !this.advancedOptionsEnabled) return;
      e.preventDefault()
      const delta = e.deltaY > 0 ? -this.step : this.step
      this[field] += delta
      this.sync()
    },
    spinMain(e, field) {
      if (this.capturedItemsCount === 0) return;
      e.preventDefault()
      const delta = e.deltaY > 0 ? -this.step : this.step
      this.mainLocation[field] += delta
      this.syncMain()
    },
    spinDiff(e, field) {
      if (this.capturedItemsCount === 0) return;
      e.preventDefault()
      const delta = e.deltaY > 0 ? -this.step : this.step
      this.diffLocation[field] += delta
      this.syncDiff()
    },
    sync() {
      UpdatePosition(this.l1, this.l2, this.w1, this.w2)
      if (this.timer) clearTimeout(this.timer)
      this.timer = setTimeout(() => {
        MoveItem(this.l1, this.l2, this.w1, this.w2)
      }, 1500)
    },
    syncMain() {
      UpdateMainLocation(this.mainLocation.L1, this.mainLocation.L2, this.mainLocation.W1, this.mainLocation.W2)
    },
    syncDiff() {
      UpdateDiffLocation(this.diffLocation.L1, this.diffLocation.L2)
    },
    async fetch() {
      const pos = await GetPosition()
      this.w1 = pos.W1
      this.w2 = pos.W2
      this.l1 = pos.L1
      this.l2 = pos.L2
      this.mainLocation = { ...pos }
    },
    scrolldown() {
      this.$nextTick(() => {
        const box = this.$refs.logbox
        box.scrollTop = box.scrollHeight
      })
    },
    toggleMode() {
      this.isMultiMode = !this.isMultiMode
      ToggleMode(this.isMultiMode)
    },
    toggleListening() {
      this.isListening = !this.isListening
      ToggleListening(this.isListening)
    },
    clearCapturedItems() {
      ClearCapturedItems()
      this.capturedItemsCount = 0
    },
    moveAllItems() {
      MoveAllItems()
    },
    toggleAdvancedOptions() {
      this.advancedOptionsEnabled = !this.advancedOptionsEnabled
    }
  },
  mounted() {
    this.fetch()
    window.runtime.EventsOn("logUpdate", (message) => {
      this.log = message.split('\n')
      this.scrolldown()
    })
    window.runtime.EventsOn("positionUpdate", (pos) => {
      this.w1 = pos.W1
      this.w2 = pos.W2
      this.l1 = pos.L1
      this.l2 = pos.L2
    })
    window.runtime.EventsOn("modeChanged", (isMultiMode) => {
      this.isMultiMode = isMultiMode
    })
    window.runtime.EventsOn("itemsCaptured", (count) => {
      this.capturedItemsCount = count
    })
    window.runtime.EventsOn("mainLocationUpdate", (location) => {
      this.mainLocation = location
    })
  }
}
</script>

<style>
body {
  font-family: Arial, sans-serif;
  background-color: #1e1e1e;
  color: #fff;
  margin: 0;
  padding: 0;
}

.box {
  background-color: #2a2a2a;
  padding: 15px;
  width: 300px;
  position: fixed;
  top: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 30px);
}

h1 {
  text-align: center;
  color: #4CAF50;
  margin: 0 0 15px 0;
  font-size: 20px;
}

.mode-toggle {
  display: flex;
  justify-content: center;
  margin-bottom: 15px;
}

.mode-toggle button {
  background-color: #3a3a3a;
  border: 1px solid #4CAF50;
  color: white;
  padding: 5px 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.mode-toggle button:first-child {
  border-radius: 5px 0 0 5px;
}

.mode-toggle button:last-child {
  border-radius: 0 5px 5px 0;
}

.mode-toggle button.active {
  background-color: #4CAF50;
}

.content-area {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow-y: auto;
}

.controls-section {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

label {
  font-weight: bold;
  margin-right: 10px;
}

.numpad {
  display: flex;
  align-items: center;
}

.numpad span {
  width: 40px;
  text-align: center;
  user-select: none;
  cursor: ns-resize;
}

.btn {
  background-color: #4CAF50;
  border: none;
  color: white;
  width: 24px;
  height: 24px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.1s;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 4px;
}

.btn:hover {
  background-color: #45a049;
}

.btn:active {
  transform: scale(0.95);
}

.btn:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.btn.left {
  margin-right: 2px;
}

.btn.right {
  margin-left: 2px;
}

select {
  background-color: #3a3a3a;
  color: white;
  border: 1px solid #4CAF50;
  padding: 5px;
  border-radius: 5px;
}

.tip {
  font-size: 12px;
  color: #888;
  text-align: center;
  margin: 10px 0;
}

#log {
  height: 120px;
  overflow-y: auto;
  background-color: #3a3a3a;
  border: 1px solid #4CAF50;
  padding: 10px;
  border-radius: 5px;
  font-size: 12px;
  margin-top: auto;
}

#log::-webkit-scrollbar {
  width: 0px;
  background: transparent;
}

button {
  background-color: #4CAF50;
  border: none;
  color: white;
  padding: 10px 15px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 14px;
  margin: 4px 2px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #45a049;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

button.active {
  background-color: #45a049;
}

.clear-btn {
  background-color: #f44336;
}

.clear-btn:hover {
  background-color: #d32f2f;
}

.captured-items {
  text-align: center;
  margin-top: 10px;
  font-weight: bold;
}

.row .numpad {
  flex: 1;
  margin-right: 10px;
}

.row .numpad:last-child {
  margin-right: 0;
}

.advanced-options-btn {
  font-size: 12px;
  padding: 5px 10px;
  width: auto;
  margin: 0;
  display: inline-block;
}

.advanced-options-btn.active {
  background-color: #45a049;
}

.advanced-options {
  margin-top: 15px;
  border-top: 1px solid #4CAF50;
  padding-top: 15px;
}

.action-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: auto;
  margin-bottom: 15px;
}

.step-size {
  margin-bottom: 15px;
}
</style>