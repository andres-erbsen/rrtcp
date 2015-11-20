#!/usr/bin/python

from mininet.topo import Topo
from mininet.net import Mininet
from mininet.node import CPULimitedHost
from mininet.link import TCLink
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

class SingleSwitchTopo(Topo):
   "Single switch connected to n hosts."
   def build(self, n=2):
       switch = self.addSwitch('s1')
       for h in range(n):
           host = self.addHost('h%s' % (h + 1))
           # 10 Mbps, 5ms delay, 10% loss, 1000 packet queue
           self.addLink(host, switch,
              bw=10, delay='5ms', loss=10, max_queue_size=1000, use_htb=True)

def test():
    topo = SingleSwitchTopo(n=2)
    net = Mininet( topo=topo, link=TCLink )
    net.start()
    
    h1, h2 = net.get('h1', 'h2')
    print h1.cmd( 'ping -c10', h2.IP() )

    net.stop()

if __name__ == '__main__':
   setLogLevel('info')
   test()
