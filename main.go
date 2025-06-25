package main

import (
    "fmt"
    "github.com/gdamore/tcell/v2"
    "log"
    "time"
    "math/rand"
)

type Monster interface {
	DrawPos(screen tcell.Screen) ()
	DrawHp(screen tcell.Screen, yPos int)
	DrawStats(screen tcell.Screen, yPos int, xPos int)
	IsAlive() bool
	GetPosition() Coordinate
	GetDex() int
	GetXpGain() int
	takeDamage()
}

type Respawner interface {
	Respawn (position Coordinate) ()	
	RespawnArea() (Coordinate, Coordinate)
}

type RespawningMonster struct {
	*basicMonster
	AreaMin Coordinate
	AreaMax Coordinate
	respawnHp int
}

type basicMonster struct {
	name string
    position Coordinate
    hp int
	dex int
	xpGain int
}

type Dungeon struct {
    Tiles []string
}

type Coordinate struct {
    x, y int
}

type Player struct {
    position Coordinate
	hp int
	dex int
	xp int
}

func (monster basicMonster) DrawPos(screen tcell.Screen) () {
	if monster.IsAlive() {
		screen.SetContent(monster.position.x, monster.position.y, []rune(monster.name)[0], nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
}

func (monster basicMonster) DrawHp(screen tcell.Screen, yPos int) {
	mStr := fmt.Sprintf("Monster HP: %d", monster.hp)
	for i, ch := range mStr {
		screen.SetContent(10+i, yPos, ch, nil, tcell.StyleDefault)
	}
}

func (monster basicMonster) DrawStats(screen tcell.Screen, yPos int, xPos int){
	line := fmt.Sprintf("%s (HP:%d DEX:%d)", monster.name, monster.hp, monster.dex)
	for i, ch := range line {
		screen.SetContent(xPos+i, yPos, ch, nil, tcell.StyleDefault)
	}
}

func (monster basicMonster) GetPosition() Coordinate{
	return monster.position
}

func (monster basicMonster) IsAlive() bool {
	return monster.hp > 0
}

func (monster basicMonster) GetDex() int {
	return monster.dex
}

func (monster basicMonster) GetXpGain() int {
	return monster.xpGain
}

func (monster *basicMonster) takeDamage() {
	monster.hp--
}

func (monster *RespawningMonster) Respawn(position Coordinate) {
	monster.hp = monster.respawnHp
	monster.position = position
}

func (monster *RespawningMonster) RespawnArea() (Coordinate, Coordinate) {
	return monster.AreaMin, monster.AreaMax
}

func (dungeon *Dungeon) findRandomTile(minY int, maxY int, minX int, maxX int) (int, int) {
	x := minX+ rand.Intn(maxX - minX)
	y := minY+ rand.Intn(maxY - minY)
	return x,y
}

func (dungeon *Dungeon) FindRandomFreeTile(minY int, maxY int, minX int, maxX int) (Coordinate, bool) {
	for range 10000 {
		x, y := dungeon.findRandomTile(minY, maxY, minX, maxX)
		if dungeon.Tiles[y][x] != '.' {
			continue
		}
		return  Coordinate{x, y}, true
	}
	// found no free space
	return Coordinate{0, 0}, false
}

func (dungeon *Dungeon) Draw(screen tcell.Screen) () {
	for y, line := range dungeon.Tiles {
		for x, ch := range line {
			screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}
}

func (player *Player) DrawPos(screen tcell.Screen) () {
	screen.SetContent(player.position.x, player.position.y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
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
	if fightingMonsterIndex > -1 && monsters[fightingMonsterIndex].IsAlive() {
		monsters[fightingMonsterIndex].DrawHp(screen, len(dungeon.Tiles)+2)
	}

	// draw monster list
	infoX := len(dungeon.Tiles[0]) + 2  // right next to dungeon
	infoY := 0
	for _, monster := range monsters {
		monster.DrawPos(screen)
		if (monster.IsAlive()){
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
	
	playerPosition := Coordinate {x: 1, y: 2}
    player := Player{         position: playerPosition,	hp: 25,  dex: 120,  xp:  0}
	
	//easy initialisation 
	ratPosition := Coordinate {x: 13, y: 6}
	trollPosition := Coordinate {x: 19, y: 8}
	assasinePosition := Coordinate {x: 24, y: 8}
	shoggothPosition := Coordinate {x: 29, y: 8}
	witchKingPosition := Coordinate {x: 35, y: 7}
	
	basicMonsters := []basicMonster{
		{name: "Troll",       position: trollPosition,  hp: 40,  dex: 2,  xpGain: 11},
		{name: "Assassin",    position: assasinePosition,  hp: 12,  dex: 11, xpGain: 20},
		{name: "Shoggoth",    position: shoggothPosition,  hp: 20,  dex: 14, xpGain: 16},
		{name: "Witch King",  position: witchKingPosition,  hp: 150, dex: 25, xpGain: 0},
	}
	
	respawningMonsters := []RespawningMonster{
		{
			basicMonster: &basicMonster{
				name:     "Rat",
				position: ratPosition,
				hp:       3,
				dex:      1,
				xpGain:   1,
			},
			AreaMin:   Coordinate{0, 0},
			AreaMax:   Coordinate{18, len(dungeon.Tiles)},
			respawnHp: 3,
		},
		{
			basicMonster: &basicMonster{
				name:     "Shoggoth",
				position: shoggothPosition,
				hp:       20,
				dex:      14,
				xpGain:   16,
			},
			AreaMin:   Coordinate{26, 0},
			AreaMax:   Coordinate{len(dungeon.Tiles[0]), len(dungeon.Tiles)},
			respawnHp: 30,
		},
	}
	
	monsters := make([]Monster, len(basicMonsters))
	for _, m := range respawningMonsters {
		monsters = append(monsters, &m)
	}
	
	for i := range basicMonsters {
		monsters[i] = &basicMonsters[i] // use & if Monster’s methods have pointer receivers
		//monsters[i] = basicMonsters[i]  // use value if pointer receivers aren’t required
	}
	
	//transformation to []Monster to later mix with other, non-basic Monsters
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
		desiredPlayerPosition := Coordinate {player.position.x + dx, player.position.y + dy}
		
		if dungeon.Tiles[desiredPlayerPosition.y][desiredPlayerPosition.x] == '#' {
			continue
		}

		fought := false
		for i := range monsters {
			if (monsters[i].GetPosition() == desiredPlayerPosition) && monsters[i].IsAlive() {						
				for true {
					if attackSuccesful(player.dex, monsters[i].GetDex()){
						monsters[i].takeDamage()
					}
					if !monsters[i].IsAlive(){
						break
					}
					
					if attackSuccesful(monsters[i].GetDex(), player.dex){
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
				
				if !monsters[3].IsAlive(){
					game_over (screen, "You defeated the evil king", tcell.ColorGreen)
				}
				
				if !monsters[i].IsAlive() {
					// move onto defeated monster's tile
					player.position = desiredPlayerPosition
					if r, ok := monsters[i].(Respawner); ok {
						AreaMin, AreaMax := r.RespawnArea()
						randomPos, positionOkay := dungeon.FindRandomFreeTile(AreaMin.y, AreaMax.y, AreaMin.x, AreaMax.x)
						if (positionOkay){
							r.Respawn(randomPos)
						}
					}
					player.gainExp(monsters[i].GetXpGain(), len(dungeon.Tiles), screen)
				}
				fought = true
				break
			}
		}

		if !fought {
			player.position  = desiredPlayerPosition
		}

		draw(dungeon, monsters, player, -1, screen)
	}
}
