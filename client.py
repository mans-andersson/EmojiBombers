import pygame
import sys
import socket
import json
import threading

pygame.init()
state = {
    "players": [],
    "blockades": [],
    "bombs": [],
    "placed_bombs": [],
    "explosions": [],
    "winner": -1,
}
screen = pygame.display.set_mode((800, 800), pygame.DOUBLEBUF)
player_image_1 = pygame.image.load("assets/player1.png").convert_alpha()
player_image_2 = pygame.image.load("assets/player2.png").convert_alpha()
player_image_3 = pygame.image.load("assets/player3.png").convert_alpha()
player_image_4 = pygame.image.load("assets/player4.png").convert_alpha()
block_image = pygame.image.load("assets/block.png").convert_alpha()
gift_image = pygame.image.load("assets/gift.png").convert_alpha()
bomb_image = pygame.image.load("assets/bomb.png").convert_alpha()
damaged_image = pygame.image.load("assets/damaged.png").convert_alpha()
dead_image = pygame.image.load("assets/dead.png").convert_alpha()
explosion_image = pygame.image.load("assets/explosion.png").convert_alpha()
victory_text = pygame.image.load("assets/victory_text.png").convert_alpha()
player_images = [player_image_1, player_image_2, player_image_3, player_image_4]
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(("127.0.0.1", 8080))
last_update = pygame.time.get_ticks()


def main():
    while True:
        events()
        keypress()
        render()


def listener():
    while True:
        new_state = receive_state()
        try:
            state.update(json.loads(new_state))
        except:
            print("Couldn't parse json")


def receive_state():
    buffer = ""
    while True:
        buffer += s.recv(1024).decode("utf-8")
        result = buffer.split(";")
        if len(result) > 2:
            for res in result[1 : len(result) - 1]:
                if len(res) > 0:
                    return res


def render():
    screen.fill((0, 0, 0))
    if state["winner"] != -1:
        render_victory_message(state["winner"])
        return
    for player in state["players"]:
        render_player(player["id"])
    for blockade in state["blockades"]:
        screen.blit(block_image, (blockade["x_pos"] - 25, blockade["y_pos"] - 25))
    for bomb in state["bombs"]:
        if bomb["spawned"]:
            screen.blit(gift_image, (bomb["x_pos"] - 25, bomb["y_pos"] - 25))
    for bomb in state["placed_bombs"]:
        if bomb["spawned"]:
            screen.blit(bomb_image, (bomb["x_pos"] - 25, bomb["y_pos"] - 25))
    for explosion in state["explosions"]:
        if explosion["spawned"]:
            screen.blit(
                explosion_image, (explosion["x_pos"] - 100, explosion["y_pos"] - 100)
            )
    pygame.display.update()


def render_player(id):
    player = state["players"][id]
    if player["spawned"]:
        if player["dead"]:
            screen.blit(dead_image, (player["x_pos"] - 50, player["y_pos"] - 50))
        elif player["damage_taken"]:
            screen.blit(damaged_image, (player["x_pos"] - 50, player["y_pos"] - 50))
        else:
            screen.blit(player_images[id], (player["x_pos"] - 50, player["y_pos"] - 50))


def render_victory_message(winnerID):
    screen.blit(player_images[winnerID], (350, 350))
    screen.blit(victory_text, (100, 460))
    pygame.display.update()


def events():
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()


def keypress():
    global last_update
    current_time = pygame.time.get_ticks()
    if current_time - last_update < 16:
        return
    last_update = current_time
    keys = pygame.key.get_pressed()
    if keys[pygame.K_UP]:
        s.sendall(b'{ "action":"up"};')
    elif keys[pygame.K_LEFT]:
        s.sendall(b'{ "action":"left"};')
    elif keys[pygame.K_DOWN]:
        s.sendall(b'{ "action":"down"};')
    elif keys[pygame.K_RIGHT]:
        s.sendall(b'{ "action":"right"};')
    elif keys[pygame.K_SPACE]:
        s.sendall(b'{ "action":"space"};')


if __name__ == "__main__":
    l = threading.Thread(target=listener, args=(), daemon=True)
    l.start()
    main()
