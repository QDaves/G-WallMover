package main

import (
	"context"
	"embed"
	"fmt"
	"os"
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
	furniID     string
	log         []string
	logMu       sync.Mutex
	lastAction  time.Time
	packetStr   string
	captureMode bool
	captureMu   sync.Mutex
	ctx         context.Context
	moveTimer   *time.Timer

	multiMode        bool
	capturedItems    map[string]struct{}
	mainLocation     struct{ L1, L2, W1, W2 int }
	diffLocation     struct{ L1, L2 int }
	isListening      bool
	listeningMu      sync.Mutex
	lastCapturedItem struct {
		ID             string
		L1, L2, W1, W2 int
		Direction      string
	}
}

func NewApp(ext *g.Ext, assets embed.FS) *App {
	return &App{
		ext:           ext,
		assets:        assets,
		capturedItems: make(map[string]struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.setupExt()
	go func() {
		a.runExt()
	}()
}

func (a *App) setupExt() {
	a.ext.Intercept(out.CHAT).With(a.handleTalk)
	a.ext.Intercept(out.SHOUT).With(a.handleTalk)
	a.ext.Intercept(in.ITEMS_2).With(a.handleItems2)
	a.ext.Intercept(in.UPDATEITEM).With(a.handleUpdateItem)
	a.ext.Intercept(out.ADDSTRIPITEM).With(a.handleAddStripItem)
}

func (a *App) runExt() {
	defer os.Exit(0)
	a.ext.Run()
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

func (a *App) UpdatePosition(l1, l2, w1, w2 int) {
	a.pos.L1 = l1
	a.pos.L2 = l2
	a.pos.W1 = w1
	a.pos.W2 = w2

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
	msg := e.Packet.ReadString()
	if msg == "#wallmover" {
		runtime.WindowShow(a.ctx)
		e.Block()
	}
}

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
	var newPos struct {
		W1, W2, L1, L2 int
		Direction      string
	}
	for _, part := range posParts {
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

	a.listeningMu.Lock()
	isListening := a.isListening
	a.listeningMu.Unlock()

	if a.multiMode && isListening {
		if _, exists := a.capturedItems[newID]; !exists {
			a.capturedItems[newID] = struct{}{}
			a.AddLogMsg(fmt.Sprintf("Item added: %s", newID))
			runtime.EventsEmit(a.ctx, "itemsCaptured", len(a.capturedItems))
			a.lastCapturedItem.ID = newID
			a.lastCapturedItem.L1 = newPos.L1
			a.lastCapturedItem.L2 = newPos.L2
			a.lastCapturedItem.W1 = newPos.W1
			a.lastCapturedItem.W2 = newPos.W2
			a.lastCapturedItem.Direction = newPos.Direction
		}
	}

	if !a.multiMode && time.Since(a.lastAction) >= 2*time.Second && (newID != a.furniID || newPos != a.pos) {
		a.furniID = newID
		a.pos = newPos
		runtime.EventsEmit(a.ctx, "positionUpdate", a.pos)
		a.AddLogMsg(fmt.Sprintf("%s: ID %s, Position: w=%d,%d l=%d,%d %s",
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
			a.AddLogMsg(fmt.Sprintf("Captured furni ID: %s", a.furniID))
		}
		return
	}
	a.captureMu.Unlock()
}

func (a *App) updatePacketStr() string {
	parts := strings.Split(a.packetStr, "\t")
	if len(parts) < 5 {
		return a.packetStr
	}
	parts[3] = fmt.Sprintf(":w=%d,%d l=%d,%d %s", a.pos.W1, a.pos.W2, a.pos.L1, a.pos.L2, a.pos.Direction)
	return strings.Join(parts, "\t")
}

func (a *App) ToggleListening(isListening bool) {
	a.listeningMu.Lock()
	a.isListening = isListening
	a.listeningMu.Unlock()
	if isListening {
		a.AddLogMsg("Started listening for new items")
	} else {
		a.AddLogMsg("Stopped listening for new items")
		a.mainLocation.L1 = a.lastCapturedItem.L1
		a.mainLocation.L2 = a.lastCapturedItem.L2
		a.mainLocation.W1 = a.lastCapturedItem.W1
		a.mainLocation.W2 = a.lastCapturedItem.W2
		a.AddLogMsg(fmt.Sprintf("Set Main location to: W1=%d, W2=%d, L1=%d, L2=%d",
			a.mainLocation.L1, a.mainLocation.L2, a.mainLocation.W1, a.mainLocation.W2))
		runtime.EventsEmit(a.ctx, "mainLocationUpdate", a.mainLocation)
	}
}

func (a *App) ClearCapturedItems() {
	a.capturedItems = make(map[string]struct{})
	a.AddLogMsg("Cleared all IDs")
	runtime.EventsEmit(a.ctx, "itemsCaptured", 0)
}

func (a *App) UpdateMainLocation(l1, l2, w1, w2 int) {
	a.mainLocation.L1 = l1
	a.mainLocation.L2 = l2
	a.mainLocation.W1 = w1
	a.mainLocation.W2 = w2
	a.AddLogMsg(fmt.Sprintf("Main location: L1=%d, L2=%d, W1=%d, W2=%d", l1, l2, w1, w2))
	a.updateMultiItemPositions()
}

func (a *App) UpdateDiffLocation(l1Diff, l2Diff int) {
	a.diffLocation.L1 = l1Diff
	a.diffLocation.L2 = l2Diff
	a.AddLogMsg(fmt.Sprintf("Location: L1=%d, L2=%d", l1Diff, l2Diff))
	a.updateMultiItemPositions()
}

func (a *App) updateMultiItemPositions() {
	count := 0
	for id := range a.capturedItems {
		l1 := a.mainLocation.L1 + (count * a.diffLocation.L1)
		l2 := a.mainLocation.L2 + (count * a.diffLocation.L2)
		updatedPacket := a.createUpdatePacket(id, l1, l2, a.mainLocation.W1, a.mainLocation.W2)
		a.ext.Send(in.UPDATEITEM, []byte(updatedPacket))
		count++
	}
}

func (a *App) createUpdatePacket(id string, l1, l2, w1, w2 int) string {
	return fmt.Sprintf("%s\t0\t0\t:w=%d,%d l=%d,%d %s", id, w1, w2, l1, l2, a.lastCapturedItem.Direction)
}

func (a *App) MoveAllItems() {
	if !a.multiMode {
		return
	}

	count := 0
	for id := range a.capturedItems {
		numID, _ := strconv.Atoi(id)
		l1 := a.mainLocation.L1 + (count * a.diffLocation.L1)
		l2 := a.mainLocation.L2 + (count * a.diffLocation.L2)
		placestuffData := fmt.Sprintf(":w=%d,%d l=%d,%d %s", a.mainLocation.W1, a.mainLocation.W2, l1, l2, a.lastCapturedItem.Direction)

		go func(id string, numID int, placestuffData string, count int) {
			time.Sleep(time.Duration(count) * 500 * time.Millisecond)
			a.ext.Send(g.Out.Id("MOVEITEM"), numID, placestuffData)
			a.AddLogMsg(fmt.Sprintf("Moved %s L1=%d, L2=%d, W1=%d, W2=%d", id, l1, l2, a.mainLocation.W1, a.mainLocation.W2))
		}(id, numID, placestuffData, count)

		count++
	}
	a.AddLogMsg(fmt.Sprintf("Moving %d items", count))
}

func (a *App) ToggleMode(isMultiMode bool) {
	a.multiMode = isMultiMode
	if isMultiMode {
		a.AddLogMsg("Switched to Multi mode")
	} else {
		a.AddLogMsg("Switched to Single mode")
	}
	runtime.EventsEmit(a.ctx, "modeChanged", isMultiMode)
}
