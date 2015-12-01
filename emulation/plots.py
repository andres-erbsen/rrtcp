#!/usr/bin/python

import numpy as np
import matplotlib.pyplot as plt
from matplotlib import mlab
from config import delayIntervals, lossIntervals

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

    tcpHist = plt.hist(tcp, bins=bins, normed=1, histtype='step', cumulative=True, label='TCP')
    udpHist = plt.hist(udp, bins=bins, normed=1, histtype='step', cumulative=True, label='UDP')
    rrtcpHist = plt.hist(rrtcp, bins=bins, normed=1, histtype='step', cumulative=True, label='RRTCP')

    plt.grid(True)
    plt.ylim(0, 1.05)
    plt.xlim(0, max(tcp + udp + rrtcp))

    plt.title('Packet Flight Time with ' + str(delay) + 'ms delay and ' + str(loss) + '% loss')
    plt.xlabel('milliseconds')
    plt.ylabel('packets')
    plt.legend(loc='upper right')

    plt.show()

if __name__ == '__main__':
    thres = 3 # must have more than 3 lines

    # iterate over all delays and losses
    for delay in delayIntervals:
        for loss in lossIntervals:
            tcp = readFile('tcp', delay, loss)
            udp = readFile('udp', delay, loss)
            rrtcp = readFile('rrtcp', delay, loss)

            # only plot if we have data on all of them
            if len(tcp) > thres and len(udp) > thres and len(rrtcp) > thres:
                plot(tcp, udp, rrtcp, delay, loss)
