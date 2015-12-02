#!/usr/bin/python

import time
import subprocess

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
    timeToRun = 30 # seconds

    outputName = name + '_' + str(delay) + '_' + str(loss)

    subprocess.call("tc qdisc change dev lo root netem delay " + str(delay) + "ms loss " + str(loss) + "%", shell=True)
    time.sleep(.1)
    subprocess.call('../../' + test + '/' + test + ' -l -d ' + str(timeToRun) + 's -address localhost:8080 > ' + outputName + '.listener.out' + ' 2>'+outputName+'.listener.err &', shell=True)
    subprocess.call('../../' + test + '/' + test + ' -d ' + str(timeToRun) + 's -address localhost:8080 > ' + outputName + '.dialer.out' + ' 2>'+outputName+'.dialer.err &', shell=True)
    time.sleep(timeToRun + .1)

    print outputName

if __name__ == '__main__':
   runTests()
