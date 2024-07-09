<template>
  <div class="box">
    <h1>G-Wall Mover</h1>
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
    <div class="row">
      <button @click="flipw" class="flip">W1/W2</button>
      <div v-if="showw" class="wbox">
        <div class="numpad tiny">
          <label for="w1">W1:</label>
          <button @click="tweak('w1', -1)" class="btn left">-</button>
          <span @wheel="spin($event, 'w1')">{{ w1 }}</span>
          <button @click="tweak('w1', 1)" class="btn right">+</button>
        </div>
        <div class="numpad tiny">
          <label for="w2">W2:</label>
          <button @click="tweak('w2', -1)" class="btn left">-</button>
          <span @wheel="spin($event, 'w2')">{{ w2 }}</span>
          <button @click="tweak('w2', 1)" class="btn right">+</button>
        </div>
      </div>
    </div>
    <div class="row">
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
</template>

<script>
import { UpdatePosition, GetPosition, MoveItem } from '../wailsjs/go/main/App'

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
      showw: false
    }
  },
  methods: {
    tweak(field, delta) {
      this[field] += delta
      this.sync()
    },
    spin(e, field) {
      e.preventDefault()
      const delta = e.deltaY > 0 ? -this.step : this.step
      this[field] += delta
      this.sync()
    },
    sync() {
      UpdatePosition(this.l1, this.l2, this.w1, this.w2)
      if (this.timer) clearTimeout(this.timer)
      this.timer = setTimeout(() => {
        MoveItem(this.l1, this.l2, this.w1, this.w2)
      }, 1500)
    },
    async fetch() {
      const pos = await GetPosition()
      this.w1 = pos.W1
      this.w2 = pos.W2
      this.l1 = pos.L1
      this.l2 = pos.L2
    },
    flipw() {
      this.showw = !this.showw
    },
    scrolldown() {
      this.$nextTick(() => {
        const box = this.$refs.logbox
        box.scrollTop = box.scrollHeight
      })
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
  }
}
</script>

<style>
body {
  font-family: sans-serif;
  background-color: #1e1e1e;
  color: #fff;
  margin: 0;
  padding: 0;
}

.box {
  background-color: #2a2a2a;
  padding: 20px;
  width: 300px;
  position: fixed;
  top: 0;
  left: 0;
}

h1 {
  text-align: center;
  color: #4CAF50;
  margin-bottom: 20px;
  font-size: 20px;
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
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

.flip {
  background-color: #4CAF50;
  border: none;
  color: white;
  padding: 5px 10px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 12px;
  cursor: pointer;
  border-radius: 5px;
  transition: background-color 0.3s;
}

.flip:hover {
  background-color: #45a049;
}

.wbox {
  display: flex;
  gap: 10px;
  margin-left: 10px;
}

.wbox .numpad {
  flex-direction: row;
  align-items: center;
}

.wbox label {
  margin-right: 5px;
}

.numpad.tiny span {
  width: 30px;
  font-size: 12px;
}

.numpad.tiny .btn {
  width: 20px;
  height: 20px;
  font-size: 14px;
}

.tip {
  font-size: 12px;
  color: #888;
  text-align: center;
  margin-bottom: 10px;
}

#log {
  height: 150px;
  overflow-y: auto;
  background-color: #3a3a3a;
  border: 1px solid #4CAF50;
  padding: 10px;
  border-radius: 5px;
  font-size: 12px;
}

#log::-webkit-scrollbar {
  width: 0px;
  background: transparent;
}
</style>