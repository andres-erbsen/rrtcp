"""
Demo of the histogram (hist) function used to plot a cumulative distribution.

"""
import numpy as np
import matplotlib.pyplot as plt
from matplotlib import mlab

def plot_loss_over_time(file_name):
    with open(file_name, 'r') as f:
        lines = f.readlines()

    events = order_events(lines)
    packets_sent = 0
    packets_received = 0
    time_stamps = []
    packet_percent = []
    
    for event in events:
        time, type = event
        if type == 's':
            packets_sent += 1     
        else:
            packets_received += 1
        time_stamps.append(time)
        packet_percent.append(float(packets_received)/packets_sent)

    plt.plot(time_stamps, packet_percent)
    plt.show()

def order_events(lines):
    events = []
    for line in lines:
        split_line = line.split(" ")

        send_time = int(split_line[0])
        recv_time = int(split_line[1])

        events.append((send_time, 's'))
        events.append((recv_time, 'r'))
    events.sort()    
    return events
        
plot_loss_over_time("tcp_160_10.out")
# 1. time on x-axis, percent packets received on y-axis
# 2. loss percentage on x-axis, 90 percentile latency on y-axis
# time packet sent, time packet received in nanoseconds
