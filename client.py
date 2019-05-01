import pygame
import sys
import game
import socket
import json
import threading

pygame.init()
player = game.Player(500, 300)
state = {}
screen = pygame.display.set_mode((1000, 600), pygame.DOUBLEBUF)
image = pygame.image.load("assets/player1.png").convert_alpha()
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(("127.0.0.1", 8080))
# s.sendall(b"HAJ\n")
# data = s.recv(1024)
# print(data.decode())


def main():
    while True:
        events()
        keypress()
        render()


def listener():
    while True:
        data = s.recv(1024).decode("utf-8")
        print(data)
        state = json.loads(data.split("}", 1)[0])
        print(state)


def render():
    screen.fill((0, 0, 0))
    screen.blit(image, (player.x_pos, player.y_pos))
    # pygame.display.flip()
    pygame.display.update()


def events():
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            pygame.quit()
            sys.exit()


def keypress():
    keys = pygame.key.get_pressed()
    if keys[pygame.K_UP]:
        s.sendall(b'{ "action":"up"}')
    if keys[pygame.K_LEFT]:
        s.sendall(b'{ "action":"left"}')
    if keys[pygame.K_DOWN]:
        s.sendall(b'{ "action":"down"}')
    if keys[pygame.K_RIGHT]:
        s.sendall(b'{ "action":"right"}')


if __name__ == "__main__":
    l = threading.Thread(target=listener, args=(), daemon=True)
    l.start()
    main()
