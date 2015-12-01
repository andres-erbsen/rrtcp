#!/usr/bin/python

from mininet.topo import Topo
from mininet.net import Mininet
from mininet.node import CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel
import time

class SingleSwitchTopo(Topo):
   "Single switch connected to n hosts."
   def build(self, delay, loss):
       switch = self.addSwitch('s1')
       for h in range(2):
           host = self.addHost('h%s' % (h + 1))
           # 10 Mbps, 5ms delay, 10% loss, 1000 packet queue
           self.addLink(host, switch,
              bw=10, delay=str(delay/2) + 'ms', loss=loss, max_queue_size=1000, use_htb=True)

def runTests():
    tcpTest = 'tcp-clock-station'
    udpTest = 'udp-clock-station'
    rrtcpTest = 'rrtcp-clock-station'

    for delay in [0, 40, 80, 160]:
        for loss in [0, 5, 10]:
            runTest(tcpTest, delay, loss, 'tcp')
            runTest(udpTest, delay, loss, 'udp')
            runTest(rrtcpTest, delay, loss, 'rrtcp')


def runTest(test, delay, loss, name):
    timeToRun = 5 # seconds
    topo = SingleSwitchTopo(delay, loss)
    net = Mininet( topo=topo, link=TCLink )
    net.start()

    h1, h2 = net.get('h1', 'h2')
    outputName = name + '_' + str(delay) + '_' + str(loss)+'.out'

    h1.cmd( '../' + test + '/' + test + ' -l -d ' + str(timeToRun) + 's -address ' + h1.IP() + ':8080 2>'+outputName+'.listener.err &' )
    h2.cmd( '../' + test + '/' + test + ' -d ' + str(timeToRun) + 's -address ' + h1.IP() + ':8080 > ' + outputName + ' 2>'+outputName+'.dialer.err &' )
    time.sleep(timeToRun + .1)

    net.stop()

if __name__ == '__main__':
   setLogLevel('info')
   runTests()
