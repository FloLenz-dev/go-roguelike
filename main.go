package main

import (
    "fmt"
    "github.com/gdamore/tcell/v2"
    "log"
    "time"
    "math/rand"
)

type Dungeon struct {
    Tiles []string
}

type Player struct {
    x, y int
	hp int
	dex int
	xp int
}

type Monster struct {
	name string
    x, y int
    hp int
	dex int
	xpGain int
}

func (dungeon *Dungeon) findRandomTile(minY int, maxY int, minX int, maxX int) (int, int) {
	x := minX+ rand.Intn(maxX - minX)
	y := minY+ rand.Intn(maxY - minY)
	return x,y
}


func (dungeon *Dungeon) FindRandomFreeTile(minY int, maxY int, minX int, maxX int) (int, int, bool) {
	for range 10000 {
		x, y := dungeon.findRandomTile(minY, maxY, minX, maxX)
		if dungeon.Tiles[y][x] != '.' {
			continue
		}
		return x, y, true
	}
	// found no free space
	return 0, 0, false
}

func (dungeon *Dungeon) Draw(screen tcell.Screen) () {
	for y, line := range dungeon.Tiles {
		for x, ch := range line {
			screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}
}

func (monster *Monster) DrawPos(screen tcell.Screen) () {
	if monster.hp > 0 {
		screen.SetContent(monster.x, monster.y, []rune(monster.name)[0], nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
}

func (monster *Monster) DrawHp(screen tcell.Screen, yPos int) {
	mStr := fmt.Sprintf("Monster HP: %d", monster.hp)
	for i, ch := range mStr {
		screen.SetContent(10+i, yPos, ch, nil, tcell.StyleDefault)
	}
}

func (player *Player) DrawPos(screen tcell.Screen) () {
	screen.SetContent(player.x, player.y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func (player *Player) DrawStats(screen tcell.Screen, yPos int ) () {
		statusStr := fmt.Sprintf("HP: %d  DEX: %d  XP: %d", player.hp, player.dex, player.xp)
		for i, ch := range statusStr {
			screen.SetContent(i, yPos, ch, nil, tcell.StyleDefault)
		}
}

func (player *Player) gainExp(xpToGain int, messageYOffSet int, screen tcell.Screen)(){
	levels_up := ((player.xp % 2) + xpToGain) / 2
	player.xp += xpToGain
	if levels_up > 0 {
		// Level-Up!
		mStr := "You gained levels!"
		for i, ch := range mStr {
			screen.SetContent(9+i, messageYOffSet + 1, ch, nil, tcell.StyleDefault)
		}
		for level := range(levels_up) {
			if rand.Intn(2) == 0 {
				mStr = "You gained 5 HP!"
				player.hp += 5
			} else {
				mStr = "You gained 1 DEX!"
				player.dex += 1
			}
			for i, ch := range mStr {
					screen.SetContent(10+i, messageYOffSet + 2 + level, ch, nil, tcell.StyleDefault)
			}
		}
		screen.Show()
		time.Sleep(5 * time.Second)
	}
}

func (monster *Monster) DrawStats(screen tcell.Screen, yPos int, xPos int){
	line := fmt.Sprintf("%s (HP:%d DEX:%d)", monster.name, monster.hp, monster.dex)
	for i, ch := range line {
		screen.SetContent(xPos+i, yPos, ch, nil, tcell.StyleDefault)
	}
}

func (monster *Monster) revive(dungeon Dungeon, minY int, maxY int, minX int, maxX int, hp int){
	// revive monster at randomized position
	x, y, ok := dungeon.FindRandomFreeTile(minY, maxY, minX, maxX)
	if ok {
		monster.x = x
		monster.y = y
		monster.hp = hp
	}
}

func calculate_desired_movement (pressedKey rune) (int,int) {
	dx, dy := 0, 0
	switch pressedKey {
		case 'w':
			dy = -1
		case 's':
			dy = 1
		case 'a':
			dx = -1
		case 'd':
			dx = 1
	}
	return dx,dy
}

func initScreen() tcell.Screen{
    screen, err := tcell.NewScreen()
    if err != nil {
        log.Fatal("Screen creation failed:", err)
    }
    if err := screen.Init(); err != nil {
        log.Fatal("Screen init failed:", err)
    }
    return screen
}

func draw(dungeon Dungeon, monsters []Monster, player Player, fightingMonsterIndex int, screen tcell.Screen) {
	
	screen.Clear()
	
	// draw
	dungeon.Draw(screen)
	
	player.DrawPos(screen)
	player.DrawStats(screen, len(dungeon.Tiles))
	
	// If fighting, draw monster HP
	if fightingMonsterIndex > -1 && monsters[fightingMonsterIndex].hp > 0 {
		monsters[fightingMonsterIndex].DrawHp(screen, len(dungeon.Tiles)+2)
	}

	// draw monster list
	infoX := len(dungeon.Tiles[0]) + 2  // right next to dungeon
	infoY := 0
	for _, monster := range monsters {
		monster.DrawPos(screen)
		if (monster.hp > 0){
			monster.DrawStats(screen, infoY, infoX)
			infoY++
		}
	}
			
	// Show everything
	screen.Show()
}

func game_over (screen tcell.Screen, messageText string, messageColor tcell.Color){
	screen.Clear()
	for i, ch := range messageText {
		screen.SetContent(5+i, 5, ch, nil, tcell.StyleDefault.Foreground(messageColor))
	}
	screen.Show()
	time.Sleep(5 * time.Second)
	screen.Fini()
	log.Fatal("Game Over")
}

func attackSuccesful (dexAttacker int, dexDefender int) bool{
	attackerDice := rand.Intn(6)
	defenderDice := rand.Intn(6)
	if (attackerDice + dexAttacker) > (defenderDice + dexDefender){
		return true
	}
	return false
}

func monitorKeyboard(screen tcell.Screen) (bool, rune) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				// Exit condition: user pressed Escape or Ctrl+C
				return true, '#'
			case tcell.KeyRune:
				// Normal character input (letters, numbers, symbols)
				return false, ev.Rune()
			default:
				// Ignore special keys like arrows, function keys, etc.
			}
		default:
			// Ignore non-keyboard events (e.g. mouse or resize events)
		}
	}
}

func initGame() (tcell.Screen, Dungeon, Player, []Monster){
	screen := initScreen()
	dungeon := Dungeon{}
	dungeon.Tiles = []string{
		"##########                              ",
		"#........#                              ",
		"#........#                              ",
		"#........#################              ",
		"#........................#              ",
		"#....#..........#........###############",
		"#........................#.............#",
		"#.........#........#######.............#",
		"#......................................#",
		"#..................#######.............#",
		"####################     ###############",
	}

    player := Player{         x: 2,  y: 2,	hp: 25,  dex: 1,  xp:  0}
	
	monsters := []Monster{
		{name: "Rat",         x: 13, y: 6,  hp: 3,   dex: 1,  xpGain: 1},
		{name: "Troll",       x: 19, y: 8,  hp: 40,  dex: 2,  xpGain: 11},
		{name: "Assassin",    x: 24, y: 8,  hp: 12,  dex: 11, xpGain: 20},
		{name: "Shoggoth",    x: 29, y: 8,  hp: 20,  dex: 14, xpGain: 16},
		{name: "Witch King",  x: 35, y: 7,  hp: 150, dex: 25, xpGain: 0},
	}
	return screen, dungeon, player, monsters
}

func main() {
	
	screen, dungeon, player, monsters := initGame()
	
    defer screen.Fini()
	time.Sleep(200 * time.Millisecond)
    draw(dungeon, monsters, player, -1, screen)

    for {
		exitRequested, pressedRune := monitorKeyboard(screen)
		if (exitRequested){
			return
		}
		dx, dy := calculate_desired_movement(pressedRune)
		newX := player.x + dx
		newY := player.y + dy
		
		if dungeon.Tiles[newY][newX] == '#' {
			continue
		}

		fought := false
		for i := range monsters {
			if monsters[i].x == newX && monsters[i].y == newY && monsters[i].hp > 0 {						
				for true {
					if attackSuccesful(player.dex, monsters[i].dex){
						monsters[i].hp--
					}
					if monsters[i].hp <= 0 {
						break
					}
					
					if attackSuccesful(monsters[i].dex, player.dex){
						player.hp--
					}
					if player.hp <= 0 {
						game_over (screen, "You've been defeated", tcell.ColorRed)
						break
					}
					
					draw(dungeon, monsters, player, i, screen)
					time.Sleep(200 * time.Millisecond)
				}
				draw(dungeon, monsters, player, i, screen)
				time.Sleep(200 * time.Millisecond)
				
				if monsters[4].hp <= 0{
					game_over (screen, "You defeated the evil king", tcell.ColorGreen)
				}
				
				if monsters[i].hp <= 0 {
					// move onto defeated monster's tile
					player.x = newX
					player.y = newY
					if monsters[i].name == "Rat" {
						monsters[i].revive(dungeon, 0, len(dungeon.Tiles), 0, 18, 3)
					}
					if monsters[i].name == "Shoggoth" {
						monsters[i].revive(dungeon, 0, len(dungeon.Tiles), 24, len(dungeon.Tiles[0]), 30)
					}
					player.gainExp(monsters[i].xpGain, len(dungeon.Tiles), screen)
				}
				fought = true
				break
			}
		}

		if !fought {
			player.x, player.y  = newX, newY
		}

		draw(dungeon, monsters, player, -1, screen)
	}
}
