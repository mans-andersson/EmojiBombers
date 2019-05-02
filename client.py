import pygame
import sys
import socket
import json
import threading

pygame.init()
state = {"players": []}
screen = pygame.display.set_mode((800, 800), pygame.DOUBLEBUF)
player_image_1 = pygame.image.load("assets/player1.png").convert_alpha()
player_image_2 = pygame.image.load("assets/player2.png").convert_alpha()
player_image_3 = pygame.image.load("assets/player3.png").convert_alpha()
player_image_4 = pygame.image.load("assets/player4.png").convert_alpha()
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
        buffer += s.recv(512).decode("utf-8")
        result = buffer.split(";")
        if len(result) > 2:
            return result[1]


def render():
    screen.fill((0, 0, 0))
    for player in state["players"]:
        render_player(player["id"])
    pygame.display.update()


def render_player(id):
    player = state["players"][id]
    screen.blit(player_images[id], (player["x_pos"] - 50, player["y_pos"] - 50))


def events():
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()


def keypress():
    global last_update
    current_time = pygame.time.get_ticks()
    if current_time - last_update < 5:
        return
    last_update = current_time
    keys = pygame.key.get_pressed()
    if keys[pygame.K_UP]:
        s.sendall(b'{ "action":"up"};')
    if keys[pygame.K_LEFT]:
        s.sendall(b'{ "action":"left"};')
    if keys[pygame.K_DOWN]:
        s.sendall(b'{ "action":"down"};')
    if keys[pygame.K_RIGHT]:
        s.sendall(b'{ "action":"right"};')


if __name__ == "__main__":
    l = threading.Thread(target=listener, args=(), daemon=True)
    l.start()
    main()
