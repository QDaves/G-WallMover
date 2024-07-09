package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	g "xabbo.b7c.io/goearth"
	"xabbo.b7c.io/goearth/shockwave/in"
	"xabbo.b7c.io/goearth/shockwave/out"
)

type App struct {
	ext    *g.Ext
	assets embed.FS
	pos    struct {
		W1, W2, L1, L2 int
		Direction      string
	}
<<<<<<< HEAD
	furniID     string
	log         []string
	logMu       sync.Mutex
	lastAction  time.Time
	packetStr   string
	captureMode bool
	captureMu   sync.Mutex
	ctx         context.Context
	moveTimer   *time.Timer
=======
	furniID         string
	logMessages     []string
	logMu           sync.Mutex
	lastActionTime  time.Time
	packetStructure string
	captureMode     bool
	captureMu       sync.Mutex
	ctx             context.Context
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
}

func NewApp(ext *g.Ext, assets embed.FS) *App {
	return &App{
		ext:    ext,
		assets: assets,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
<<<<<<< HEAD
	a.setupExt()
	go a.runExt()
}

func (a *App) setupExt() {
	a.ext.Intercept(out.CHAT).With(a.handleTalk)
	a.ext.Intercept(out.SHOUT).With(a.handleTalk)
	a.ext.Intercept(in.ITEMS_2).With(a.handleItems2)
	a.ext.Intercept(in.UPDATEITEM).With(a.handleUpdateItem)
	a.ext.Intercept(out.ADDSTRIPITEM).With(a.handleAddStripItem)
}

func (a *App) runExt() {
=======
	a.setupExtension()
	go a.runExtension()
}

func (a *App) setupExtension() {
	a.ext.Intercept(out.CHAT).With(a.handleTalk)
	a.ext.Intercept(out.SHOUT).With(a.handleTalk)
	a.ext.Intercept(in.ITEMS_2).With(a.handleItems2Packet)
	a.ext.Intercept(in.UPDATEITEM).With(a.handleUpdateItemPacket)
	a.ext.Intercept(out.ADDSTRIPITEM).With(a.handleAddStripItem)
}

func (a *App) runExtension() {
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
	a.ext.Run()
	log.Println("running")
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

<<<<<<< HEAD
=======
func (a *App) MoveItem(l1, l2, w1, w2 int) {
	if a.furniID == "" {
		a.AddLogMessage("Place some wall furni first")
		return
	}

	a.pos.L1 = l1
	a.pos.L2 = l2
	a.pos.W1 = w1
	a.pos.W2 = w2

	a.removeWallItem()
	a.placeWallItem()
	a.AddLogMessage("Moved wall item")
	a.lastActionTime = time.Now()
}

>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
func (a *App) UpdatePosition(l1, l2, w1, w2 int) {
	a.pos.L1 = l1
	a.pos.L2 = l2
	a.pos.W1 = w1
	a.pos.W2 = w2

<<<<<<< HEAD
	updatedPacket := a.updatePacketStr()
	a.ext.Send(in.UPDATEITEM, []byte(updatedPacket))
	a.AddLogMsg(fmt.Sprintf("Updated Location: w=%d,%d l=%d,%d %s", a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction))

	if a.moveTimer != nil {
		a.moveTimer.Stop()
	}
	a.moveTimer = time.AfterFunc(1500*time.Millisecond, a.moveWallItem)
}

func (a *App) moveWallItem() {
	placestuffData := fmt.Sprintf(":w=%d,%d l=%d,%d %s", a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction)
	numFurniID, _ := strconv.Atoi(a.furniID)
	a.ext.Send(g.Out.Id("MOVEITEM"), numFurniID, placestuffData)
}

func (a *App) MoveItem(l1, l2, w1, w2 int) {
	a.pos.L1 = l1
	a.pos.L2 = l2
	a.pos.W1 = w1
	a.pos.W2 = w2
	a.moveWallItem()
}

func (a *App) AddLogMsg(msg string) {
	a.logMu.Lock()
	defer a.logMu.Unlock()
	a.log = append(a.log, msg)
	if len(a.log) > 100 {
		a.log = a.log[1:]
	}
	runtime.EventsEmit(a.ctx, "logUpdate", strings.Join(a.log, "\n"))
=======
	updatedPacket := a.updatePacketStructure()
	a.ext.Send(in.ITEMS_2, []byte(updatedPacket))
	a.AddLogMessage(fmt.Sprintf("position: w=%d,%d l=%d,%d %s", a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction))
}

func (a *App) AddLogMessage(message string) {
	a.logMu.Lock()
	defer a.logMu.Unlock()
	a.logMessages = append(a.logMessages, message)
	if len(a.logMessages) > 100 {
		a.logMessages = a.logMessages[1:]
	}
	runtime.EventsEmit(a.ctx, "logUpdate", strings.Join(a.logMessages, "\n"))
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
}

func (a *App) GetPosition() map[string]interface{} {
	return map[string]interface{}{
		"W1":        a.pos.W1,
		"W2":        a.pos.W2,
		"L1":        a.pos.L1,
		"L2":        a.pos.L2,
		"Direction": a.pos.Direction,
	}
}

func (a *App) handleTalk(e *g.InterceptArgs) {
<<<<<<< HEAD
	msg := e.Packet.ReadString()
	if msg == "#wallmover" {
=======
	message := e.Packet.ReadString()
	if message == "#wallmover" {
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
		runtime.WindowShow(a.ctx)
		e.Block()
	}
}

<<<<<<< HEAD
func (a *App) handleItems2(e *g.InterceptArgs) {
	a.packetStr = e.Packet.ReadString()
	a.handleItemPacket(a.packetStr, "Item placed")
}

func (a *App) handleUpdateItem(e *g.InterceptArgs) {
	a.packetStr = e.Packet.ReadString()
	a.handleItemPacket(a.packetStr, "Item updated")
}

func (a *App) handleItemPacket(packetStr, actionType string) {
	parts := strings.Split(packetStr, "\t")
	if len(parts) < 5 {
		return
	}
	newID := parts[0]
	pos := parts[3]
	posParts := strings.Fields(pos)
=======
func (a *App) handleItems2Packet(e *g.InterceptArgs) {
	a.packetStructure = e.Packet.ReadString()
	a.handleItemPacket(a.packetStructure, "Item placed")
}

func (a *App) handleUpdateItemPacket(e *g.InterceptArgs) {
	a.packetStructure = e.Packet.ReadString()
	a.handleItemPacket(a.packetStructure, "New Item ID Applied")
}

func (a *App) handleItemPacket(packetStructure string, actionType string) {
	parts := strings.Split(packetStructure, "\t")
	if len(parts) < 5 {
		return
	}
	newFurniID := parts[0]
	position := parts[3]
	positionParts := strings.Fields(position)
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
	var newPos struct {
		W1, W2, L1, L2 int
		Direction      string
	}
<<<<<<< HEAD
	for _, part := range posParts {
=======
	for _, part := range positionParts {
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
		if strings.HasPrefix(part, ":w=") {
			wParts := strings.Split(strings.TrimPrefix(part, ":w="), ",")
			if len(wParts) == 2 {
				newPos.W1, _ = strconv.Atoi(wParts[0])
				newPos.W2, _ = strconv.Atoi(wParts[1])
			}
		} else if strings.HasPrefix(part, "l=") {
			lParts := strings.Split(strings.TrimPrefix(part, "l="), ",")
			if len(lParts) == 2 {
				newPos.L1, _ = strconv.Atoi(lParts[0])
				newPos.L2, _ = strconv.Atoi(lParts[1])
			}
		} else if part == "r" || part == "l" {
			newPos.Direction = part
		}
	}
<<<<<<< HEAD
	if time.Since(a.lastAction) >= 2*time.Second && (newID != a.furniID || newPos != a.pos) {
		a.furniID = newID
		a.pos = newPos
		runtime.EventsEmit(a.ctx, "positionUpdate", a.pos)
		a.AddLogMsg(fmt.Sprintf("%s: ID %s, Position: w=%d,%d l=%d,%d %s",
=======
	if time.Since(a.lastActionTime) >= 2*time.Second && (newFurniID != a.furniID || newPos != a.pos) {
		a.furniID = newFurniID
		a.pos = newPos
		runtime.EventsEmit(a.ctx, "positionUpdate", a.pos)
		a.AddLogMessage(fmt.Sprintf("%s: ID %s, Position: w=%d,%d l=%d,%d %s",
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
			actionType, a.furniID, a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction))
	}
}

func (a *App) handleAddStripItem(e *g.InterceptArgs) {
	a.captureMu.Lock()
	if a.captureMode {
		a.captureMode = false
		a.captureMu.Unlock()
		packetContent := e.Packet.ReadString()
		parts := strings.Split(packetContent, " ")
		if len(parts) >= 3 {
			a.furniID = parts[2]
<<<<<<< HEAD
			a.AddLogMsg(fmt.Sprintf("Captured furni ID: %s", a.furniID))
=======
			a.AddLogMessage(fmt.Sprintf("Captured furni ID: %s", a.furniID))
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
		}
		return
	}
	a.captureMu.Unlock()
}

<<<<<<< HEAD
func (a *App) updatePacketStr() string {
	parts := strings.Split(a.packetStr, "\t")
	if len(parts) < 5 {
		return a.packetStr
=======
func (a *App) removeWallItem() {
	a.ext.Send(out.ADDSTRIPITEM, []byte("new item "+a.furniID))
}

func (a *App) placeWallItem() {
	placestuffData := fmt.Sprintf("%s :w=%d,%d l=%d,%d %s", a.furniID, a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction)
	a.ext.Send(out.PLACESTUFF, []byte(placestuffData))
}

func (a *App) updatePacketStructure() string {
	parts := strings.Split(a.packetStructure, "\t")
	if len(parts) < 5 {
		return a.packetStructure
>>>>>>> 6d9554d632396d8767a32690167928b9484ee64f
	}
	parts[3] = fmt.Sprintf(":w=%d,%d l=%d,%d %s", a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction)
	return strings.Join(parts, "\t")
}
