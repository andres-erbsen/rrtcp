import numpy as np
import matplotlib.pyplot as plt
from matplotlib import mlab

def readFile(type, delay, loss):
    packetTimes = []
    filename = type + '_' + str(delay) + '_' + str(loss) + '.out'
    f = open('data/' + filename, 'r')
    for line in f:
        start, end = line.split(' ')
        packetTimes.append(float(int(end) - int(start))/10**6)

    return packetTimes

def plot(tcp, udp, rrtcp, delay, loss):
    bins = 50
    if tcp:
        plt.hist(tcp, bins=bins, normed=1, histtype='step', cumulative=True)
    if udp:
        plt.hist(udp, bins=bins, normed=1, histtype='step', cumulative=True)
    if rrtcp:
        plt.hist(rrtcp, bins=bins, normed=1, histtype='step', cumulative=True)

    plt.grid(True)
    plt.ylim(0, 1.05)
    plt.xlim(0, max(tcp + udp + rrtcp))
    plt.title('Packet Flight Time with ' + str(delay) + 'ms delay and ' + str(loss) + '% loss')
    plt.xlabel('milliseconds')
    plt.ylabel('packets')

    plt.show()

if __name__ == '__main__':
    delay = 80
    loss = 10
    tcp = readFile('tcp', delay, loss)
    udp = readFile('udp', delay, loss)
    rrtcp = readFile('rrtcp', delay, loss)
    plot(tcp, udp, rrtcp, delay, loss)
