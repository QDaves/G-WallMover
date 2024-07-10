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
}

func NewApp(ext *g.Ext, assets embed.FS) *App {
	return &App{
		ext:    ext,
		assets: assets,
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
	if time.Since(a.lastAction) >= 2*time.Second && (newID != a.furniID || newPos != a.pos) {
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
