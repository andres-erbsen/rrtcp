#!/usr/bin/python

import time
import subprocess
from config import delayIntervals, lossIntervals, tcpTestLocation, udpTestLocation, rrtcpTestLocation

def runTests():
    for delay in delayIntervals:
        for loss in lossIntervals:
            runTest(tcpTestLocation, delay, loss, 'tcp')
            #runTest(udpTestLocation, delay, loss, 'udp')
            runTest(rrtcpTestLocation, delay, loss, 'rrtcp')

def runTest(testLocation, delay, loss, name):
    timeToRun = 30 # seconds

    outputName = name + '_' + str(delay) + '_' + str(loss)

    subprocess.call("tc qdisc change dev lo root netem delay " + str(delay) + "ms loss " + str(loss) + "%", shell=True)
    time.sleep(.1)
    print "Opening tcpdump"
    subprocess.Popen("sudo tcpdump -i lo -tt port 8080 > " + outputName + ".tcpdump.out", shell=True)
    print "Starting application"
    time.sleep(.1)
    subprocess.call( testLocation + ' -l -d ' + str(timeToRun) + 's -address localhost:8080 > ' + outputName + '.listener.out' + ' 2>'+outputName+'.listener.err &', shell=True)
    subprocess.call( testLocation + ' -d ' + str(timeToRun) + 's -address localhost:8080 > ' + outputName + '.dialer.out' + ' 2>'+outputName+'.dialer.err &', shell=True)
    time.sleep(timeToRun + .1)

    print outputName

if __name__ == '__main__':
   runTests()
