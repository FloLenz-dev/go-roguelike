package main

import (
    "fmt"
    "github.com/gdamore/tcell/v2"
    "log"
    "time"
    "math/rand"
)

var dungeon = []string{
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
	xp_gain int
}

func findRandomFreeTile(monsters []Monster, player Player, min_x int, max_x int) (int, int, bool) {
	for attempt := 0; attempt < 1000; attempt++ {
		x := min_x + rand.Intn(max_x)
		y := rand.Intn(len(dungeon))
		if dungeon[y][x] != '.' {
			continue
		}
		return x, y, true
	}

	// found no free space
	return 0, 0, false
}

func main() {
	
    screen, err := tcell.NewScreen()
    if err != nil {
        log.Fatal(err)
    }
    if err := screen.Init(); err != nil {
        log.Fatal(err)
    }
    defer screen.Fini()

    player := Player{2, 2, 20, 2, 0}
	
	    monsters := []Monster{
        {"Rat", 13, 6, 3, 1, 1},
        {"Troll", 19, 8, 40, 3, 16},
        {"Black Knight", 24, 8, 15, 15, 15},
        {"Shoggoth", 29, 8, 20, 3, 10},
        {"Witch King", 35, 7, 150, 40, 0},
    }
	
	draw := func(monster *Monster) {
        
		screen.Clear()
        for y, line := range dungeon {
            for x, ch := range line {
                screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
            }
        }
		
		// Draw monsters
        for _, m := range monsters {
            if m.hp > 0 {
                screen.SetContent(m.x, m.y, []rune(m.name)[0], nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
            }
        }
		
		// Draw player
        screen.SetContent(player.x, player.y, '@', nil, tcell.StyleDefault.Bold(true))
        screen.Show()
		
		// Draw player HP + XP
		statusY := len(dungeon)
		statusStr := fmt.Sprintf("HP: %d  DEX: %d  XP: %d", player.hp, player.dex, player.xp)

		for i, ch := range statusStr {
			screen.SetContent(i, statusY, ch, nil, tcell.StyleDefault)
		}
		
		//draw monster list
		infoX := len(dungeon[0]) + 2  // right next to dungeon
		infoY := 0
		
		// If fighting, draw monster HP
		if monster != nil && monster.hp > 0 {
			mStr := fmt.Sprintf("Monster HP: %d", monster.hp)
			for i, ch := range mStr {
				screen.SetContent(10+i, len(dungeon)+2, ch, nil, tcell.StyleDefault)
			}
		}

		for _, m := range monsters {
			if m.hp <= 0 {
				continue
			}
			line := fmt.Sprintf("%s (HP:%d DEX:%d)", m.name, m.hp, m.dex)
			for i, ch := range line {
				screen.SetContent(infoX+i, infoY, ch, nil, tcell.StyleDefault)
			}
			infoY++
		}
				
		
		//Show everything
        screen.Show()
    }

	fight := func(m *Monster) {
		var player_dice, monster_dice int	
		for m.hp > 0 && player.hp > 0 {
			player_dice = rand.Intn(6)
			monster_dice = rand.Intn(6)
			if (monster_dice + m.dex) < (player_dice + player.dex){
			    m.hp--
			}
			draw(m)
			time.Sleep(200 * time.Millisecond)
			if m.hp <= 0 {
				break
			}
			
			player_dice = rand.Intn(6)
			monster_dice = rand.Intn(6)
			if (player_dice + player.dex) < (monster_dice + m.dex){
			    player.hp--
			}
			draw(m)
			time.Sleep(200 * time.Millisecond)
			draw(nil)
		}

		if player.hp <= 0 {
			screen.Clear()
			msg := "You died. Game Over."
			for i, ch := range msg {
				screen.SetContent(5+i, 5, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
			}
			screen.Show()
			time.Sleep(3 * time.Second)
			screen.Fini()
			log.Fatal("Game Over")
		}
		if monsters[4].hp <= 0{
			screen.Clear()
			msg := "You defeated the evil King. Game Over."
			for i, ch := range msg {
				screen.SetContent(5+i, 5, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
			}
			screen.Show()
			time.Sleep(3 * time.Second)
			screen.Fini()
			log.Fatal("Game Over")
		}
	}

    draw(nil)

    for {
        ev := screen.PollEvent()
        switch ev := ev.(type) {
        case *tcell.EventKey:
            switch ev.Key() {
            case tcell.KeyEscape, tcell.KeyCtrlC:
                return
            default:
                dx, dy := 0, 0
                switch ev.Rune() {
                case 'w':
                    dy = -1
                case 's':
                    dy = 1
                case 'a':
                    dx = -1
                case 'd':
                    dx = 1
                }
                newX := player.x + dx
                newY := player.y + dy
                if dungeon[newY][newX] == '#' {
                    break
                }

                fought := false
                for i := range monsters {
                    if monsters[i].x == newX && monsters[i].y == newY && monsters[i].hp > 0 {
                        fight(&monsters[i])
                        if monsters[i].hp <= 0 {
                            // move onto defeated monster's tile
                            player.x = newX
                            player.y = newY
							if monsters[i].name == "Rat" {
								//create new Rat in first dungeon
                                x, y, ok := findRandomFreeTile(monsters, player, 0, 18)
                                if ok {
									monsters[i].x = x
									monsters[i].y = y
									monsters[i].hp = 3
								}
							}
							if monsters[i].name == "Shoggoth" {
								//create new Rat in second dungeon
                                x, y, ok := findRandomFreeTile(monsters, player, 24, len(dungeon[0])-24)
                                if ok {
									monsters[i].x = x
									monsters[i].y = y
									monsters[i].hp = 3
								}
							}
								levels_up := ((player.xp % 3) + monsters[i].xp_gain) / 3
							    player.xp += monsters[i].xp_gain
							if levels_up > 0 {
								// Level-Up!
								mStr := "You gained levels!"
								for i, ch := range mStr {
									screen.SetContent(9+i, len(dungeon)+1, ch, nil, tcell.StyleDefault)
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
											screen.SetContent(10+i, len(dungeon)+ 2 + level, ch, nil, tcell.StyleDefault)
									}
								}
								screen.Show()
								time.Sleep(5 * time.Second)
							}
                        }
                        fought = true
                        break
                    }
                }

                if !fought {
                    player.x = newX
                    player.y = newY
                }

                draw(nil)
			}
        }
    }
}
