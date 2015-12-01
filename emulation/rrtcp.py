#!/usr/bin/python

from mininet.topo import Topo
from mininet.net import Mininet
from mininet.node import CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections, waitListening
from mininet.log import setLogLevel
import time

timeToRun = 60
def runTests():
    tcpTest = 'tcp-clock-station'
    udpTest = 'udp-clock-station'
    rrtcpTest = 'rrtcp-clock-station'

    for delay in [120]:# [0, 40, 80, 160]:
        for loss in [3]:# [0, 5, 10]:
            runTest(tcpTest, delay, loss, 'tcp')
            runTest(udpTest, delay, loss, 'udp')
            runTest(rrtcpTest, delay, loss, 'rrtcp')


def runTest(test, delay, loss, name):
    net = Mininet()
    h1 = net.addHost('h1')
    h2 = net.addHost('h2')
    net.addLink(h1, h2, cls=TCLink, bw=10, delay=str(delay) + 'ms', loss=loss, max_queue_size=1000, use_htb=True)

    net.start()
    takeMeasurements(h1, h2, test, delay, loss, name)
    net.stop()

def takeMeasurements(h1, h2, test, delay, loss, name):
    outputName = name + '_' + str(delay) + '_' + str(loss)+'.out'
    print '*** starting listener...'
    h1.cmd('../' + test + '/' + test + ' -l -d ' + str(timeToRun) + 's -address ' + h1.IP() + ':8080 2>'+outputName+'.listener.err &' )
    time.sleep(.1)
    print '*** running dialer...'
    h2.cmd('../' + test + '/' + test + ' -d ' + str(timeToRun) + 's -address ' + h1.IP() + ':8080 > ' + outputName + ' 2>'+outputName+'.dialer.err' )
    h1.cmd('killall '+name)

if __name__ == '__main__':
   setLogLevel('info')
   runTests()
