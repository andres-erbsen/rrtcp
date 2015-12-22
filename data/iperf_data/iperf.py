#!/usr/bin/python

import os
import numpy


def readFile():
    send = []
    rcv = []
    sendId = 5
    rcvId = 4
    filename = 'rrtcp.txt'
    with open(filename) as datafile:
        lines = datafile.readlines()

    for line in lines:
        split = line.split(",")
        id = split[5]
        rate = split[8]
        if int(id) == sendId:
            send.append(int(rate))
        else:
            rcv.append(int(rate))

    return (send, rcv)



if __name__ == '__main__':
    send, rcv = readFile()
    sendArr = numpy.array(send)
    rcvArr = numpy.array(rcv)

    print "send mean:", numpy.mean(sendArr)/1e6, "Mbits/sec"
    print "send dev:", numpy.std(sendArr)/1e6, "Mbits/sec"
    print "rcv mean:", numpy.mean(rcvArr)/1e6, "Mbits/sec"
    print "rcv dev:", numpy.std(rcvArr)/1e6, "Mbits/sec"

    # Test 1, Olga/Asya's computer:

    # alone:
    # send mean: 2.88666804706 Mbits/sec
    # send dev: 1.84028475924 Mbits/sec
    # rcv mean: 0.749631264368 Mbits/sec
    # rcv dev: 0.370029607402 Mbits/sec

    # with tcp
    # send mean: 2.77526693333 Mbits/sec
    # send dev: 1.69548369341 Mbits/sec
    # rcv mean: 2.31358222222 Mbits/sec
    # rcv dev: 1.24814878256 Mbits/sec

    # with rrtcp
    # send mean: 1.61606208 Mbits/sec
    # send dev: 1.22209210692 Mbits/sec
    # rcv mean: 2.66356216822 Mbits/sec
    # rcv dev: 0.903527392379 Mbits/sec

    # Test 2, Andres/Asya's computer:
