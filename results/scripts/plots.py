#!/usr/bin/python

import os
import numpy as np
import matplotlib.pyplot as plt
from matplotlib import mlab
from config import delayIntervals, lossIntervals

large_packet_time = float(1010) # sufficiently large time to pretend to be a missed packet
large_bin_size = 500

def read_file(type, delay, loss):
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

def read_file_unknown(path, type):
    packetTimes = []
    recv_filename = path + '/' + type + '_dialer.out'
    send_filename = path + '/' + type + '_listener.out'
    with open(recv_filename, 'r') as recv_file:
        recv_lines = recv_file.readlines()
    with open(send_filename, 'r') as send_file:
        send_lines = send_file.readlines()

    print "Percent loss for %s: %f" % (type, float(len(recv_lines))/len(send_lines))

    for i in xrange(len(send_lines) - len(recv_lines)):
        packetTimes.append(large_packet_time)

    f = open(recv_filename, 'r')
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

def plot_unknown(tcp, udp, rrtcp, name):
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
    plt.title('Distribution of Application to Application Latency Per Packet\nAnonymous to Non-Anonymous Tor on MIT Guest')
    plt.xlabel('flight time in milliseconds')
    plt.ylabel('fraction of packets received')
    plt.legend(loc='lower right')

    plt.savefig("../plots/%s.png" % name)
    plt.savefig("../plots/%s.eps" % name)
    plt.show()
    plt.clf()

def plot_unknown_data():
    for dirname, dirnames, filenames in os.walk('../test_data/'):
        for subdirname in dirnames:
            path = os.path.join(dirname, subdirname)

            tcp = read_file_unknown(path, 'tcp')
            udp = read_file_unknown(path, 'udp')
            rrtcp = read_file_unknown(path, 'rrtcp')

            plot_unknown(tcp, udp, rrtcp, subdirname)

def plot_delay_loss_data():
    for delay in delayIntervals:
        for loss in lossIntervals:
            tcp = read_file('tcp', delay, loss)
            udp = read_file('udp', delay, loss)
            rrtcp = read_file('rrtcp', delay, loss)
            plot(tcp, udp, rrtcp, delay, loss)

if __name__ == '__main__':
    path = '../test_data/stata-stata'
    tcp = read_file_unknown(path, 'tcp')
    udp = read_file_unknown(path, 'udp')
    rrtcp = read_file_unknown(path, 'rrtcp')

    plot_unknown(tcp, udp, rrtcp, 'stata-stata')
    #path = '../test_data/toranon-tor-mitguest20'
    #tcp = read_file_unknown(path, 'tcp')
    #rrtcp = read_file_unknown(path, 'rrtcp')
    #plot_unknown(tcp, None, rrtcp, 'toranon-tor-mitguest20')
    #:w
    #plot_unknown_data()
