import numpy as np
import matplotlib.pyplot as plt
from matplotlib import mlab

large_packet_time = float(1010) # sufficiently large time to pretend to be a missed packet
large_bin_size = 100 

def readFile(type, delay, loss):
    packetTimes = []
    recv_filename = type + '_' + str(delay) + '_' + str(loss) + '.dialer.out'
    send_filename = type + '_' + str(delay) + '_' + str(loss) + '.listener.out'
    with open('send_recv_data/' + recv_filename, 'r') as recv_file:
        recv_lines = recv_file.readlines()
    with open('send_recv_data/' + send_filename, 'r') as send_file:
        send_lines = send_file.readlines()

    print "Percent loss for %s: %f" % (type, float(len(recv_lines))/len(send_lines))

    for i in xrange(len(send_lines) - len(recv_lines)):
        packetTimes.append(large_packet_time)

    recv_filename = type + '_' + str(delay) + '_' + str(loss) + '.dialer.out'
    f = open('send_recv_data/' + recv_filename, 'r')
    for line in f:
        start, end = line.split(' ')
        packetTimes.append(float(int(end) - int(start))/10**6)

    return packetTimes

def plot(tcp, udp, rrtcp, delay, loss):
    bins = large_bin_size
    if tcp:
        plt.hist(tcp, bins=bins, normed=1, histtype='step', cumulative=True, label='tcp')
    if udp:
        plt.hist(udp, bins=bins, normed=1, histtype='step', cumulative=True, label='udp')
    if rrtcp:
        plt.hist(rrtcp, bins=bins, normed=1, histtype='step', cumulative=True, label='rrtcp')

    plt.grid(True)
    plt.yticks([i/20.0 for i in range(0, 21)])
    plt.ylim(0, 1.05)
    plt.xlim(0, 1000)
    plt.title('Packet Flight Time with ' + str(delay) + 'ms delay and ' + str(loss) + '% loss')
    plt.xlabel('flight time in milliseconds')
    plt.ylabel('fraction of packets received')
    plt.legend(loc='lower right')

    plt.savefig("send_recv_data/plots/%d_%d.png" % (delay, loss))
    plt.savefig("send_recv_data/plots/%d_%d.eps" % (delay, loss))
    plt.clf()

if __name__ == '__main__':
    for delay in [0, 40, 80, 160]:
        for loss in [0, 5, 10]:
            tcp = readFile('tcp', delay, loss)
            udp = readFile('udp', delay, loss)
            rrtcp = readFile('rrtcp', delay, loss)
            plot(tcp, udp, rrtcp, delay, loss)
