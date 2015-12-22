from config import delayIntervals, lossIntervals

def read_file(type, delay, loss):
    packetTimes = []

    print type, delay, loss

    recv_filename = type + '_' + str(delay) + '_' + str(loss) + '.dialer.out'
    send_filename = type + '_' + str(delay) + '_' + str(loss) + '.listener.out'
    send_network_filename = type + '_' + str(delay) + '_' + str(loss) + '.listener.tcpdump.out'
    recv_network_filename = type + '_' + str(delay) + '_' + str(loss) + '.dialer.tcpdump.out'

    with open('new_tcpdump_data/' + recv_filename, 'r') as recv_file:
        recv_lines = recv_file.readlines()
    with open('new_tcpdump_data/' + send_filename, 'r') as send_file:
        send_lines = send_file.readlines()
    with open('new_tcpdump_data/' + send_network_filename, 'r') as network_file:
        send_network_lines = network_file.readlines()
    with open('new_tcpdump_data/' + recv_network_filename, 'r') as network_file:
        recv_network_lines = network_file.readlines()

    seq_lines = filter(lambda x: "seq" in x and ":" in x[x.index("seq"):], send_network_lines)
    parsed_packets = {}
    network_packets = []
         
    # This code does what we want, let's not think too hard about how it works
    for i in xrange(len(send_lines) - 1):
        line = seq_lines[i]
        seq_index = line.index("seq")
        start_index = line[seq_index:].index(":") + 1 + seq_index
        end_index = line[start_index:].index(",") + start_index
       
        arrow_index = line.index(">")
        port_start = arrow_index + 11
        port_end = line[port_start:].index(":") + port_start

        data_end = line[start_index:end_index]
        port = line[port_start:port_end]
    
        send_timestamp = get_timestamp_micro(line) 
        recv_line = filter(lambda x : ":" + data_end in x and ("." + port) in x, recv_network_lines)[0]
        
        recv_timestamp = get_timestamp_micro(recv_line)

        id = port + ":" + data_end
        if id not in parsed_packets:
            parsed_packets[id] = send_timestamp
            network_packets.append((send_timestamp, recv_timestamp))
    
    for i in xrange(len(network_packets)):
        line = send_lines[i]
        send_timestamp = line[:-4]
        recv_line = filter(lambda x: str(int(line)) in x, recv_lines)
        if len(recv_line) > 0:
            recv_line = recv_line[0]
            recv_timestamp = int(recv_line.split()[1][:-3])

            packetTimes.append((int(send_timestamp), network_packets[i][0], network_packets[i][1], recv_timestamp))

    with open('new_tcpdump_delays/' + type + '_' + str(delay) + '_' + str(loss) + '.delay.log', 'w') as f:
        for time in packetTimes:
            format_str = '%d %d %d %d\n' % (time[0], time[1], time[2], time[3])
            f.write(format_str)

def get_timestamp_micro(line):
    timestamp = int(float(line.split(" ")[0]) * 1000000)
    return timestamp

#for delay in delayIntervals:
#    for loss in lossIntervals:
#        read_file("rrtcp", delay, loss)
read_file("rrtcp", 160, 5)
read_file("rrtcp", 0, 0)
read_file("tcp", 0, 0)
read_file("tcp", 160, 5)
